package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/go-git/go-git/v5"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/environments"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/operations"
	"github.com/pufferpanel/pufferpanel/v2/programs"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	tmpDir, err := os.MkdirTemp("", "puffertemplatetest")
	panicIf(err)
	defer os.RemoveAll(tmpDir)

	config.DatabaseDialect.Set("sqlite3", false)
	config.DatabaseUrl.Set("file:test.db?cache=shared&mode=memory", false)
	config.ConsoleForward.Set(true, false)
	config.ServersFolder.Set(filepath.Join(tmpDir, "servers"), false)
	config.BinariesFolder.Set(filepath.Join(tmpDir, "binaries"), false)
	config.CacheFolder.Set(filepath.Join(tmpDir, "cache"), false)
	config.LogsFolder.Set(filepath.Join(tmpDir, "logs"), false)
	templateFolder := filepath.Join(tmpDir, "templates")

	_ = os.MkdirAll(config.ServersFolder.Value(), 0755)
	_ = os.MkdirAll(config.BinariesFolder.Value(), 0755)
	_ = os.MkdirAll(config.CacheFolder.Value(), 0755)
	_ = os.MkdirAll(config.LogsFolder.Value(), 0755)
	_ = os.MkdirAll(templateFolder, 0755)

	newPath := os.Getenv("PATH")
	fullPath, _ := filepath.Abs(config.BinariesFolder.Value())
	if !strings.Contains(newPath, fullPath) {
		_ = os.Setenv("PATH", newPath+":"+fullPath)
	}

	logging.Initialize(false)

	//this may require a DB, so we are going to pretend we have one
	//because of how code works, we're going to abuse our own system
	//db, err := database.GetConnection()
	//panicIf(err)

	//get all templates
	_, err = git.PlainClone(templateFolder, false, &git.CloneOptions{
		URL:           "https://github.com/PufferPanel/templates",
		ReferenceName: "refs/heads/v2.6",
	})
	panicIf(err)

	environments.LoadModules()
	operations.LoadOperations()

	templateFolders, err := os.ReadDir(templateFolder)
	panicIf(err)

	templates := make([]*TestTemplate, 0)

	for _, folder := range templateFolders {
		if !folder.IsDir() || strings.HasPrefix(folder.Name(), ".") {
			continue
		}

		if _, err := os.Stat(filepath.Join(templateFolder, folder.Name(), ".skip")); err == nil {
			continue
		}

		files, err := os.ReadDir(filepath.Join(templateFolder, folder.Name()))
		panicIf(err)

		for _, file := range files {
			filePath := filepath.Join(templateFolder, folder.Name(), file.Name())
			if strings.HasSuffix(file.Name(), ".json") {
				tmp := &TestTemplate{}
				tmp.Name = strings.TrimSuffix(file.Name(), ".json")

				skip := false
				for _, v := range os.Args[1:] {
					if match(v, tmp.Name) {
						skip = true
						break
					}
				}
				if skip {
					continue
				}

				tmp.Template, err = os.ReadFile(filePath)
				panicIf(err)

				var dataFile os.FileInfo
				if dataFile, err = os.Stat(filepath.Join(templateFolder, folder.Name(), tmp.Name+".txt")); os.IsNotExist(err) {
					dataFile, _ = os.Stat(filepath.Join(templateFolder, folder.Name(), "data.txt"))
				}
				if dataFile != nil {
					tmp.Data, err = os.ReadFile(filepath.Join(templateFolder, folder.Name(), dataFile.Name()))
					panicIf(err)
				}

				templates = append(templates, tmp)
			}
		}
	}

	//now... we can create servers from each one of them
	for _, template := range templates {
		log.Printf("Starting test for %s", template.Name)

		buf := bytes.NewReader(template.Template)

		log.Printf("Creating server")
		prg := programs.CreateProgram()
		err = json.NewDecoder(buf).Decode(prg)
		panicIf(err)
		prg.Identifier = strings.ReplaceAll(template.Name, "+", "")

		//use data file to fill in extra template data
		if template.Data != nil {
			reader := bytes.NewReader(template.Data)
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.Split(line, "=")
				key := parts[0]
				value := parts[1]
				variable := prg.Variables[key]
				variable.Value = value
				prg.Variables[key] = variable
			}
		}

		err = programs.Create(prg)
		panicIf(err)

		err = prg.Install()
		panicIf(err)

		err = runServer(prg)
		panicIf(err)

		err = prg.Destroy()
		panicIf(err)
	}
}

func runServer(prg *programs.Program) (err error) {
	err = prg.Start()
	panicIf(err)

	c := make(chan error, 1)
	go func() {
		c <- prg.RunningEnvironment.WaitForMainProcess()
	}()
	t := time.After(time.Minute * 1)

	select {
	case <-t:
		break
	case <-c:
		break
	}

	err = prg.Stop()
	panicIf(err)

	return prg.GetEnvironment().WaitForMainProcessFor(time.Minute * 3)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// wildCardToRegexp converts a wildcard pattern to a regular expression pattern.
func wildCardToRegexp(pattern string) string {
	var result strings.Builder
	for i, literal := range strings.Split(pattern, "*") {

		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return result.String()
}

func match(pattern string, value string) bool {
	result, _ := regexp.MatchString(wildCardToRegexp(pattern), value)
	return result
}

type TestTemplate struct {
	Template []byte
	Name     string
	Data     []byte
}
