package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func readDataTxtFile(fileName string) (map[string]interface{}, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer utils.Close(file)

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
	defer utils.Close(file)

	result := make([]*TestData, 0)
	err = json.NewDecoder(file).Decode(&result)
	return result, err
}

func buildTests() []*TestScenario {
	var templateFolder = CmdFlags.TemplateFolder

	if templateFolder == "" {
		if templateFolder = os.Getenv("TEMPLATE_PATH"); templateFolder == "" {
			templateFolder = filepath.Join(CmdFlags.WorkingDir, "templates")
			err := os.MkdirAll(templateFolder, 0755)
			panicIf(err)

			log.Printf("Cloning template repo")
			_, err = git.PlainClone(templateFolder, false, &git.CloneOptions{
				URL:           "https://github.com/PufferPanel/templates",
				ReferenceName: plumbing.ReferenceName(CmdFlags.GitRef),
				SingleBranch:  true,
				Depth:         1,
			})
			panicIf(err)
		}
	}

	templateFolders, err := os.ReadDir(templateFolder)
	panicIf(err)

	testScenarios := make([]*TestScenario, 0)

	for _, folder := range templateFolders {
		if !folder.IsDir() || strings.HasPrefix(folder.Name(), ".") {
			continue
		}

		log.Printf("Scanning %s", folder.Name())

		if _, err = os.Stat(filepath.Join(templateFolder, folder.Name(), ".skip")); err == nil {
			log.Printf("  Skipping as skip file present")
			continue
		}

		var files []os.DirEntry
		files, err = os.ReadDir(filepath.Join(templateFolder, folder.Name()))
		panicIf(err)

		for _, file := range files {
			if file.Name() == "data.json" {
				continue
			}

			log.Printf("  Scanning %s", file.Name())

			if len(CmdFlags.Files) > 0 {
				test := false
				for _, t := range CmdFlags.Files {
					if t == file.Name() {
						test = true
						break
					}
				}
				if !test {
					log.Printf("  Skipping %s as not in files", file.Name())
					continue
				}
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

				template := pufferpanel.Server{}
				err = json.NewDecoder(bytes.NewReader(tmp.Template)).Decode(&template)
				panicIf(err)

				_, err = os.Stat(filepath.Join(templateFolder, folder.Name(), "data.json"))
				if err == nil {
					tests, err := readDataJsonFile(filepath.Join(templateFolder, folder.Name(), "data.json"))
					for _, v := range tests {
						log.Printf("  Considering %s", v.Name)
						testScenarios = append(testScenarios, &TestScenario{
							Name: v.Name,
							Test: &TestTemplate{
								Template:           tmp.Template,
								Name:               tmp.Name,
								Variables:          v.Variables,
								Environment:        v.Environment,
								IgnoreExitCode:     v.IgnoreExitCode,
								RuntimeRequirement: v.RuntimeRequirement,
							},
						})
					}
					panicIf(err)
				} else if !os.IsNotExist(err) {
					panicIf(err)
				} else {
					//no data json, which means it's a single test
					//but, each template could support envs, so auto-process each
					if len(template.SupportedEnvironments) > 0 {
						for _, v := range template.SupportedEnvironments {
							z := &TestTemplate{
								Template:           tmp.Template,
								Name:               tmp.Name,
								Environment:        make(map[string]interface{}),
								Variables:          make(map[string]interface{}),
								RuntimeRequirement: tmp.RuntimeRequirement,
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
							log.Printf("  Considering %s", scenario.Name)
							testScenarios = append(testScenarios, scenario)
						}
					} else {
						log.Printf("  Considering %s", tmp.Name)
						testScenarios = append(testScenarios, &TestScenario{
							Name: tmp.Name,
							Test: tmp,
						})
					}
				}
			}
		}
	}

	finalScenarioList := make([]*TestScenario, 0)
	for _, scenario := range testScenarios {
		skip := false
		for _, v := range CmdFlags.Skip {
			if utils.CompareWildcard(scenario.Name, v) {
				skip = true
				break
			}
		}
		if skip {
			for _, v := range CmdFlags.Required {
				if utils.CompareWildcard(scenario.Name, v) {
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

	return finalScenarioList
}

func ioCopy(dest io.Writer, src io.Reader) {
	var buf = make([]byte, 1024)
	var readErr error
	var n int

	for readErr != io.EOF {
		n, readErr = src.Read(buf)
		_, _ = dest.Write(buf[:n])
	}
}

func createCreateBody(scenario *TestScenario) io.ReadCloser {
	model := &models.ServerCreation{
		Server: pufferpanel.Server{},
		NodeId: 0,
		Name:   scenario.Name,
	}

	_ = json.NewDecoder(bytes.NewReader(scenario.Test.Template)).Decode(&model.Server)

	for k, v := range scenario.Test.Variables {
		val, ok := model.Variables[k]
		if ok {
			val.Value = v
			model.Variables[k] = val
		}
	}

	for k, v := range scenario.Test.Environment {
		if k == "type" {
			model.Environment.Type = v.(string)

		}
		model.Environment.Metadata[k] = v
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(model)
	return io.NopCloser(buf)
}
