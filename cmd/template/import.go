package template

import (
	"archive/zip"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/panel/database"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const ReleaseUrl = "https://github.com/PufferPanel/templates/archive/v2.zip"

var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports templates defined at https://github.com/PufferPanel/templates",
	Run:   runImport,
	Args:  cobra.NoArgs,
}

func runImport(cmd *cobra.Command, args []string) {
	err := pufferpanel.LoadConfig("")
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err.Error())
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		fmt.Printf("Error connecting to database: %s\n", err.Error())
		return
	}

	dir, err := ioutil.TempDir("", "pufferpaneltemplates")
	if err != nil {
		fmt.Printf("Error creating temp directory: %s\n", err.Error())
		return
	}
	defer func(d string) {
		err := os.RemoveAll(d)
		if err != nil {
			fmt.Printf("Error deleting temp directory %s: %s\n", d, err.Error())
		}
	}(dir)

	targetFile, err := ioutil.TempFile("", "pufferpaneltemplates*.zip")
	if err != nil {
		fmt.Printf("Error creating temp file: %s\n", err.Error())
		os.Exit(1)
		return
	}

	defer func(f *os.File) {
		pufferpanel.Close(targetFile)
		err := os.Remove(f.Name())
		if err != nil {
			fmt.Printf("Error deleting file %s: %s\n", f.Name(), err.Error())
		}
	}(targetFile)

	fmt.Printf("Downloading %s\n", ReleaseUrl)
	response, err := http.Get(ReleaseUrl)
	if err != nil {
		fmt.Printf("Error downloading: %s\n", err.Error())
		return
	}
	defer pufferpanel.Close(response.Body)

	_, err = io.Copy(targetFile, response.Body)
	if err != nil {
		fmt.Printf("Error downloading: %s\n", err.Error())
		return
	}
	_ = response.Body.Close()

	err = unzip(targetFile.Name(), dir)
	if err != nil {
		fmt.Printf("Error extracting zip: %s\n", err.Error())
		return
	}

	ts := &services.Template{DB: db}
	err = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}

		readmePath := filepath.Join(filepath.Dir(p), "README.md")

		fmt.Printf("Importing %s\n", p)
		name := strings.TrimSuffix(filepath.Base(info.Name()), filepath.Ext(info.Name()))
		err = importTemplate(name, p, readmePath, ts)
		return err
	})

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func writeFile(source *zip.File, target string) error {
	s, err := source.Open()
	defer pufferpanel.Close(s)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	defer pufferpanel.Close(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, s)
	return err
}

func unzip(sourceZip, targetDir string) error {
	zipFile, err := zip.OpenReader(sourceZip)
	defer pufferpanel.Close(zipFile)
	if err != nil {
		return err
	}

	for _, f := range zipFile.File {
		if f.FileInfo().IsDir() {
			continue
		}
		fmt.Printf("Extracting %s\n", f.Name)
		exportPath := filepath.Join(targetDir, f.Name)
		err := os.MkdirAll(filepath.Dir(exportPath), 0644)
		if err != nil {
			return err
		}
		err = writeFile(f, exportPath)
		if err != nil {
			return err
		}
	}
	return nil
}
