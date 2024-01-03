package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var workingDir string

func main() {
	flag.StringVar(&workingDir, "workDir", "", "")
	flag.Parse()

	err := filepath.WalkDir(workingDir, reformatFile)
	if err != nil {
		panic(err)
	}
}

func reformatFile(path string, info fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	nonAbsPath := strings.TrimPrefix(path, workingDir+string(filepath.Separator))

	if info.IsDir() {
		return nil
	}
	if strings.HasPrefix(nonAbsPath, ".") {
		return nil
	}
	if info.Name() == "data.json" || info.Name() == "spec.json" {
		return nil
	}
	if !strings.HasSuffix(info.Name(), ".json") {
		return nil
	}

	fmt.Printf("Reformatting %s\n", nonAbsPath)

	template, err := readFile(path)
	if err != nil {
		return err
	}
	return writeFile(path, template)
}

func readFile(path string) (*pufferpanel.Server, error) {
	template := &pufferpanel.Server{}
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&template)
	return template, err
}

func writeFile(path string, template *pufferpanel.Server) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(template)
	return err
}
