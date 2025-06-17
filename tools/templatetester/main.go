package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"io/fs"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// The purpose of this is to simply test all templates in our repo to the best of our ability
// This will download all the templates, spin up a fake server using it, and attempt to run everything
// For now, it will just test that we can create, install, start, and stop the server
// Note though, we will wait 1 minute before we stop a server, to "give it time" to start up
// After 5 minutes, if the server did not stop, we will consider it a failed template
// Arguments are templates to ignore, for ones which require data that we cannot actually safely test
// such as ones which need Steam credentials or to actually own the game
func main() {
	if len(CmdFlags.Skip) != 0 {
		log.Printf("Skip rules: %s", strings.Join(CmdFlags.Skip, " "))
	}
	if len(CmdFlags.Required) != 0 {
		log.Printf("Require rules: %s", strings.Join(CmdFlags.Required, " "))
	}
	if len(CmdFlags.Files) != 0 {
		log.Printf("Files to test rules: %s", strings.Join(CmdFlags.Files, " "))
	}

	var err error
	if CmdFlags.WorkingDir == "" {
		tmpDir := os.TempDir()
		pattern := "puffertemplatetest"

		toDelete, _ := filepath.Glob(filepath.Join(tmpDir, pattern+"*"))
		for _, z := range toDelete {
			err := os.RemoveAll(z)
			panicIf(err)
		}

		CmdFlags.WorkingDir, err = os.MkdirTemp("", pattern)
		panicIf(err)
	} else {
		err = filepath.WalkDir(CmdFlags.WorkingDir, func(path string, info fs.DirEntry, err error) error {
			if path == CmdFlags.WorkingDir {
				return err
			}
			return os.RemoveAll(path)
		})
		if err != nil && !os.IsNotExist(err) {
			panicIf(err)
		}
	}

	if CmdFlags.DeleteTemp {
		defer func() {
			_ = os.RemoveAll(CmdFlags.WorkingDir)
		}()
	}

	var tests = buildTests()

	if len(tests) == 0 {
		log.Printf("No tests were found")
		return
	}

	if CmdFlags.PrintTests {
		data := make([]string, 0)
		for _, v := range tests {
			data = append(data, "\""+v.Name+"\"")
		}
		msg := strings.Join(data, ",")
		envFile := os.ExpandEnv("GITHUB_ENV")
		if envFile == "" {
			log.Printf("No env file found, printing to console")
			log.Println(msg)
		} else {
			file, err := os.OpenFile(envFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			panicIf(err)
			defer func() {
				_ = file.Close()
			}()
			_, err = file.WriteString(fmt.Sprintf("TEMPLATES=%s\n", msg))
			panicIf(err)
		}
	}

	//we have our test set, let's kick off a panel instance
	//for this, we're going to run the binary, and wait for the "service" to start up (using the unix socket)
	//once that's done, we'll then create our servers
	//this is the best way to truly "model" what's going on

	waiter := make(chan bool, 1)

	var unixSocketPath = fmt.Sprintf("/tmp/pufferpanel-%d.sock", time.Now().Unix())
	_ = os.Remove(unixSocketPath)
	listener, err := net.ListenPacket("unixgram", unixSocketPath)
	panicIf(err)

	var dbConn = "file:" + filepath.Join(CmdFlags.WorkingDir, "test.db") + "?cache=shared"

	go socketListener(listener, waiter)

	c := exec.Command(CmdFlags.PufferpanelBinary, "runService")

	defer func() {
		if c.Process.Pid != 0 {
			_ = c.Process.Kill()
		}
		_ = listener.Close()
		_ = os.Remove(unixSocketPath)
	}()

	go func() {
		c.Dir = CmdFlags.WorkingDir
		c.Env = append(
			os.Environ(),
			"PUFFER_PANEL_DATABASE_URL="+dbConn,
			"PUFFER_LOGS=logs",
			"PUFFER_DAEMON_CONSOLE_FORWARD=true",
			"PUFFER_DAEMON_DATA_ROOT="+CmdFlags.WorkingDir,
			"PUFFER_WEB_HOST=0.0.0.0:8080",
			"NOTIFY_SOCKET="+unixSocketPath,
			"GIN_MODE=release",
		)
		stdout, err := c.StdoutPipe()
		panicIf(err)
		stderr, err := c.StderrPipe()
		panicIf(err)

		panicIf(c.Start())

		go ioCopy(os.Stdout, stdout)
		go ioCopy(os.Stderr, stderr)

		e := c.Wait()
		var d *exec.ExitError
		if errors.As(e, &d) {
			if d.Error() == "signal: killed" {
				e = nil
			}
		}
		panicIf(e)
	}()

	//wait for panel to be up, so the db is fully created and we're good to go
	<-waiter

	//now we can inject our admin user in, so we can proceed to spin up the servers
	log.Println("Starting database edits")
	db, err := gorm.Open(sqlite.Open(dbConn))
	panicIf(err)
	panicIf(initLoginAdminUser(db))

	//now, start the web calls
	//the concern is the session length, we'll "force" it to last for 24 hours. if it expires, then the tests should
	//be failed anyways
	client := &http.Client{}
	session, err := createSession(db)
	panicIf(err)

	log.Println("Now starting tests")

	for i := range tests {
		runTest(client, session, tests[i])
	}
}

func runTest(client *http.Client, session string, test *TestScenario) {
	log.Println("\nStarting: " + test.Name)
	urlPrefix := CmdFlags.Host + "/api/servers/" + test.Name

	var data []byte
	var err error

	//create server
	_, err = call(client, &http.Request{
		Method: "PUT",
		URL:    createUrl(urlPrefix),
		Header: createHeaders(session),
		Body:   createCreateBody(test),
	})
	panicIf(err)

	//install server
	_, err = call(client, &http.Request{
		Method: "POST",
		URL:    createUrl(urlPrefix + "/install"),
		Header: createHeaders(session),
	})
	panicIf(err)

	//wait for install to complete
	for {
		time.Sleep(10 * time.Second)
		data, err = call(client, &http.Request{
			Method: "GET",
			URL:    createUrl(urlPrefix + "/status"),
			Header: createHeaders(session),
		})
		panicIf(err)

		var status pufferpanel.ServerRunning
		err = json.NewDecoder(bytes.NewReader(data)).Decode(&status)
		panicIf(err)

		if !status.Installing {
			break
		}
	}

	//start server
	_, err = call(client, &http.Request{
		Method: "POST",
		URL:    createUrl(urlPrefix + "/start"),
		Header: createHeaders(session),
	})
	panicIf(err)

	//wait for 5 minutes
	started := time.Now()
	for {
		time.Sleep(1 * time.Minute)
		data, err = call(client, &http.Request{
			Method: "GET",
			URL:    createUrl(urlPrefix + "/status"),
			Header: createHeaders(session),
		})
		panicIf(err)

		var status pufferpanel.ServerRunning
		err = json.NewDecoder(bytes.NewReader(data)).Decode(&status)
		panicIf(err)

		if !status.Running {
			panicIf(errors.New("server did not run for 5 minutes"))
		}
		if math.Abs(time.Since(started).Seconds()) >= 360 {
			break
		}
	}

	//stop server
	_, err = call(client, &http.Request{
		Method: "POST",
		URL:    createUrl(urlPrefix + "/stop"),
		Header: createHeaders(session),
	})
	panicIf(err)

	//wait for the stop
	started = time.Now()
	for {
		time.Sleep(1 * time.Minute)
		data, err = call(client, &http.Request{
			Method: "GET",
			URL:    createUrl(urlPrefix + "/status"),
			Header: createHeaders(session),
		})
		panicIf(err)

		var status pufferpanel.ServerRunning
		err = json.NewDecoder(bytes.NewReader(data)).Decode(&status)
		panicIf(err)

		if !status.Running {
			break
		}
		if math.Abs(time.Since(started).Seconds()) >= 360 {
			panicIf(errors.New("server did not stop after 360 seconds"))
		}
	}

	//delete server
	_, err = call(client, &http.Request{
		Method: "DELETE",
		URL:    createUrl(urlPrefix),
		Header: createHeaders(session),
	})

	panicIf(err)
}

func socketListener(listener net.PacketConn, ch chan bool) {
	defer func() {
		_ = listener.Close()
	}()

	const NotifyReady = "READY=1"
	const NotifyShutdown = "STOPPING=1"

	for {
		buf := make([]byte, 16)
		n, _, err := listener.ReadFrom(buf)
		if err != nil {
			return
		}
		panicIf(err)

		str := string(buf[:n])

		if str == NotifyReady {
			ch <- true
		} else if str == NotifyShutdown {
			break
		}
	}
}

var loginAdminUser = &models.User{
	Username:       "loginAdminUser",
	Email:          "admin@example.com",
	OtpActive:      false,
	HashedPassword: "asdf",
}

func initLoginAdminUser(db *gorm.DB) error {
	err := db.Create(loginAdminUser).Error
	if err != nil {
		return err
	}

	perms := &models.Permissions{
		UserId: &loginAdminUser.ID,
		Scopes: []*scopes.Scope{scopes.ScopeAdmin},
	}
	err = db.Create(perms).Error
	return err
}

func createSession(db *gorm.DB) (string, error) {
	ss := &services.Session{DB: db}
	token, err := ss.CreateForUser(loginAdminUser)
	if err != nil {
		return "", err
	}
	err = db.Model(&models.Session{}).Where("token = ?", token).Update("expiration_time", time.Now().Add(time.Hour*24)).Error
	return token, err
}

func call(client *http.Client, request *http.Request) (data []byte, err error) {
	response, err := client.Do(request)
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	if err != nil {
		return
	}

	data, err = io.ReadAll(response.Body)
	if err == io.EOF {
		err = nil
	}

	if err != nil || response.StatusCode < 200 || response.StatusCode >= 400 {
		err = fmt.Errorf("failed to create server (%d)\n%s", response.StatusCode, data)
	}
	return
}

func createUrl(str string) *url.URL {
	u, err := url.Parse(str)
	panicIf(err)
	return u
}

func createHeaders(session string) http.Header {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+session)
	return headers
}
