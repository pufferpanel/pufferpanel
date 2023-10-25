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

	err := filepath.Walk(workingDir, reformatFile)
	if err != nil {
		panic(err)
	}
}

func reformatFile(path string, info fs.FileInfo, err error) error {
	nonAbsPath := strings.TrimPrefix(path, workingDir+string(filepath.Separator))

	if info.IsDir() {
		return nil
	}
	if strings.HasPrefix(nonAbsPath, ".") {
		return nil
	}
	if info.Name() == "data.json" {
		return nil
	}
	if !strings.HasSuffix(info.Name(), ".json") {
		return nil
	}

	fmt.Printf("Reformatting %s\n", nonAbsPath)

	template := &pufferpanel.Server{}
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&template)
	if err != nil {
		return err
	}

	_ = file.Close()

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	return err
}
