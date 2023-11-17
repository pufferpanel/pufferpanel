package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/servers"
	"github.com/spf13/cast"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var templatesToSkip []string
var mustTest []string

// The purpose of this is to simply test all templates in our repo to the best of our ability
// This will download all the templates, spin up a fake server using it, and attempt to run everything
// For now, it will just test that we can create, install, start, and stop the server
// Note though, we will wait 1 minute before we stop a server, to "give it time" to start up
// After 5 minutes, if the server did not stop, we will consider it a failed template
// Arguments are templates to ignore, for ones which require data that we cannot actually safely test
// such as ones which need Steam credentials or to actually own the game
func main() {
	var gitRef string
	var skipStr string
	var requiredStr string
	var templateFolder string
	var deleteTemp bool
	var workingDir string
	var reuse bool
	flag.StringVar(&gitRef, "gitref", "refs/heads/v3", "")
	flag.StringVar(&skipStr, "skip", "", "")
	flag.StringVar(&requiredStr, "require", "", "")
	flag.StringVar(&templateFolder, "path", "", "")
	flag.BoolVar(&deleteTemp, "delete", true, "")
	flag.StringVar(&workingDir, "workDir", "", "")
	flag.BoolVar(&reuse, "reuse", false, "")
	flag.Parse()

	if skipStr != "" {
		templatesToSkip = strings.Split(skipStr, ",")
		fmt.Printf("Skip rules: %s\n", strings.Join(templatesToSkip, " "))
	}
	if requiredStr != "" {
		mustTest = strings.Split(requiredStr, ",")
		fmt.Printf("Require rules: %s\n", strings.Join(mustTest, " "))
	}

	var err error
	if workingDir == "" {
		tmpDir := os.TempDir()
		pattern := "puffertemplatetest"

		toDelete, _ := filepath.Glob(filepath.Join(tmpDir, pattern+"*"))
		for _, z := range toDelete {
			err := os.RemoveAll(z)
			panicIf(err)
		}

		workingDir, err = os.MkdirTemp("", pattern)
		panicIf(err)
	} else {
		err = os.RemoveAll(workingDir)
		panicIf(err)
		err = os.Mkdir(workingDir, 0755)
		panicIf(err)
	}

	if deleteTemp {
		defer os.RemoveAll(workingDir)
	}

	config.DatabaseDialect.Set("sqlite3", false)
	config.DatabaseUrl.Set("file:test.db?cache=shared&mode=memory", false)
	config.ConsoleForward.Set(true, false)
	config.ServersFolder.Set(filepath.Join(workingDir, "servers"), false)
	config.BinariesFolder.Set(filepath.Join(workingDir, "binaries"), false)
	config.CacheFolder.Set(filepath.Join(workingDir, "cache"), false)
	config.LogsFolder.Set(filepath.Join(workingDir, "logs"), false)

	_ = os.MkdirAll(config.ServersFolder.Value(), 0755)
	_ = os.MkdirAll(config.BinariesFolder.Value(), 0755)
	_ = os.MkdirAll(config.CacheFolder.Value(), 0755)
	_ = os.MkdirAll(config.LogsFolder.Value(), 0755)

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
	if templateFolder == "" {
		templateFolder = filepath.Join(workingDir, "templates")
		err = os.MkdirAll(templateFolder, 0755)
		panicIf(err)

		log.Printf("Cloning template repo")
		_, err = git.PlainClone(templateFolder, false, &git.CloneOptions{
			URL:           "https://github.com/PufferPanel/templates",
			ReferenceName: plumbing.ReferenceName(gitRef),
		})
		panicIf(err)
	}

	var templateFolders []os.DirEntry
	templateFolders, err = os.ReadDir(templateFolder)
	panicIf(err)

	testScenarios := make([]*TestScenario, 0)

	for _, folder := range templateFolders {
		if !folder.IsDir() || strings.HasPrefix(folder.Name(), ".") {
			continue
		}

		if _, err = os.Stat(filepath.Join(templateFolder, folder.Name(), ".skip")); err == nil {
			continue
		}

		var files []os.DirEntry
		files, err = os.ReadDir(filepath.Join(templateFolder, folder.Name()))
		panicIf(err)

		for _, file := range files {
			if file.Name() == "data.json" {
				continue
			}

			filePath := filepath.Join(templateFolder, folder.Name(), file.Name())
			if strings.HasSuffix(file.Name(), ".json") {
				tmp := &TestTemplate{}
				tmp.Name = strings.TrimSuffix(file.Name(), ".json")

				tmp.Template, err = os.ReadFile(filePath)
				panicIf(err)

				_, err = os.Stat(filepath.Join(templateFolder, folder.Name(), "data.txt"))
				if err == nil {
					tmp.Variables, err = readDataTxtFile(filepath.Join(templateFolder, folder.Name(), "data.txt"))
					panicIf(err)
				} else if !os.IsNotExist(err) {
					panicIf(err)
				}

				_, err = os.Stat(filepath.Join(templateFolder, folder.Name(), "data.json"))
				if err == nil {
					tests, err := readDataJsonFile(filepath.Join(templateFolder, folder.Name(), "data.json"))
					for _, v := range tests {
						testScenarios = append(testScenarios, &TestScenario{
							Name: v.Name,
							Test: &TestTemplate{
								Template:       tmp.Template,
								Name:           tmp.Name,
								Variables:      v.Variables,
								Environment:    v.Environment,
								IgnoreExitCode: v.IgnoreExitCode,
							},
						})
					}
					panicIf(err)
				} else if !os.IsNotExist(err) {
					panicIf(err)
				} else {
					//no data json, which means it's a single test
					//but, each template could support envs, so auto-process each

					template := pufferpanel.Server{}
					err = json.NewDecoder(bytes.NewReader(tmp.Template)).Decode(&template)
					panicIf(err)

					if len(template.SupportedEnvironments) > 0 {
						for _, v := range template.SupportedEnvironments {
							z := &TestTemplate{
								Template:    tmp.Template,
								Name:        tmp.Name,
								Environment: make(map[string]interface{}),
								Variables:   make(map[string]interface{}),
							}

							scenario := &TestScenario{
								Name: z.Name,
								Test: z,
							}
							if v.Type != "host" {
								scenario.Name = scenario.Name + "-" + v.Type
							}

							for r, p := range v.Metadata {
								scenario.Test.Environment[r] = p
							}

							scenario.Test.Environment["type"] = v.Type

							testScenarios = append(testScenarios, scenario)
						}
					} else {
						testScenarios = append(testScenarios, &TestScenario{
							Name: tmp.Name,
							Test: tmp,
						})
					}
				}
			}
		}
	}

	var docker *client.Client
	ctx := context.Background()

	finalScenarioList := make([]*TestScenario, 0)
	for _, scenario := range testScenarios {
		skip := false
		for _, v := range templatesToSkip {
			if match(v, scenario.Name) {
				skip = true
				break
			}
		}
		if skip {
			for _, v := range mustTest {
				if match(v, scenario.Name) {
					skip = false
					break
				}
			}

			if skip {
				log.Printf("Skipping %s", scenario.Name)
				continue
			}
		}
		log.Printf("Will run test for: %s", scenario.Name)
		finalScenarioList = append(finalScenarioList, scenario)
	}

	//now... we can create servers from each one of them
	for _, scenario := range finalScenarioList {
		log.Printf("Starting test for %s", scenario.Name)

		template := scenario.Test

		if strings.HasSuffix(scenario.Name, "-docker") {
			//kill off any existing docker containers
			if docker == nil {
				docker, err = client.NewClientWithOpts(client.FromEnv)
				panicIf(err)
				docker.NegotiateAPIVersion(ctx)
			}

			opts := types.ContainerListOptions{
				Filters: filters.NewArgs(),
			}

			opts.All = true
			opts.Filters.Add("name", scenario.Name)

			existingContainers, err := docker.ContainerList(ctx, opts)
			panicIf(err)
			if len(existingContainers) > 0 {
				err = docker.ContainerRemove(ctx, template.Name, types.ContainerRemoveOptions{
					Force: true,
				})
				panicIf(err)
			}
		}

		buf := bytes.NewReader(template.Template)

		log.Printf("Creating server")
		prg := servers.CreateProgram()
		err = json.NewDecoder(buf).Decode(prg)
		panicIf(err)
		prg.Identifier = strings.ReplaceAll(scenario.Name, "+", "")

		if template.Variables != nil {
			for k, v := range template.Variables {
				variable := prg.Variables[k]
				variable.Value = v
				prg.Variables[k] = variable
			}
		}

		if template.Environment != nil && len(template.Environment) > 0 {
			prg.Environment = pufferpanel.MetadataType{
				Type:     cast.ToString(template.Environment["type"]),
				Metadata: template.Environment,
			}
		}

		if reuse {
			_ = os.Remove(filepath.Join(config.ServersFolder.Value(), prg.Identifier+".json"))
		} else {
			_ = os.RemoveAll(filepath.Join(config.ServersFolder.Value(), prg.Identifier))
		}

		prg, err = servers.Create(prg)
		panicIf(err)

		err = prg.Install()
		panicIf(err)

		err = runServer(prg)
		panicIf(err)

		running, err := prg.IsRunning()
		panicIf(err)

		if running {
			panic(errors.New("server is still running"))
		}

		if !template.IgnoreExitCode && prg.RunningEnvironment.GetLastExitCode() != 0 {
			panicIf(fmt.Errorf("exit code status %d", prg.RunningEnvironment.GetLastExitCode()))
		}

		err = prg.Destroy()
		panicIf(err)
	}
}

func readDataTxtFile(fileName string) (map[string]interface{}, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer pufferpanel.Close(file)

	result := make(map[string]interface{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		key := parts[0]
		value := parts[1]
		result[key] = value
	}
	return result, nil
}

func readDataJsonFile(fileName string) ([]*TestData, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer pufferpanel.Close(file)

	result := make([]*TestData, 0)
	err = json.NewDecoder(file).Decode(&result)
	return result, err
}

func runServer(prg *servers.Server) (err error) {
	err = prg.Start()
	panicIf(err)

	c := make(chan error, 1)
	go func() {
		c <- prg.RunningEnvironment.WaitForMainProcess()
	}()
	t := time.After(time.Minute * 1)

	//we need to make sure we were running for a minute
	//if we did not, something went wrong
	running, err := prg.IsRunning()
	panicIf(err)

	if !running {
		panic(errors.New("server did not run for a minute"))
	}

	select {
	case <-t:
		break
	case err = <-c:
		panicIf(err)
		break
	}

	err = prg.Stop()
	panicIf(err)

	return prg.GetEnvironment().WaitForMainProcessFor(time.Minute * 5)
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
	result, _ := regexp.MatchString("^"+wildCardToRegexp(pattern)+"$", value)
	return result
}

type TestScenario struct {
	Name string
	Test *TestTemplate
}

type TestTemplate struct {
	Template       []byte
	Name           string
	Variables      map[string]interface{}
	Environment    map[string]interface{}
	IgnoreExitCode bool
}

type TestData struct {
	Name           string                 `json:"name"`
	Variables      map[string]interface{} `json:"variables"`
	Environment    map[string]interface{} `json:"environment"`
	IgnoreExitCode bool                   `json:"ignoreExitCode"`
}
