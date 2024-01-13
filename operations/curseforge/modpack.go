package curseforge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadModpack(file File) error {
	var cacheZipFolder = getCacheFolderForFile(file)
	var cacheZipFileLocation = getCacheFilePath(file)

	//see if the file already exists, if so, use it instead
	if fi, err := os.Lstat(cacheZipFileLocation); err == nil && !fi.IsDir() && fi.Size() > 0 {
		return nil
	}

	err := os.MkdirAll(cacheZipFolder, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	tmpFile, err := os.CreateTemp(cacheZipFolder, "tmp-*.zip")
	if err != nil {
		return err
	}
	defer pufferpanel.Close(tmpFile)

	response, err := pufferpanel.Http().Get(file.DownloadUrl)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusNotFound {
		return errNoFile
	}
	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		return err
	}
	pufferpanel.Close(tmpFile)
	pufferpanel.CloseResponse(response)

	err = os.Rename(tmpFile.Name(), cacheZipFileLocation)
	if err != nil {
		return err
	}

	return nil
}

func getCacheFolderForFile(file File) string {
	return filepath.Join(config.CacheFolder.Value(), "curseforge", fmt.Sprintf("%d", file.Id))
}

func getCacheFilePath(file File) string {
	return filepath.Join(getCacheFolderForFile(file), "download.zip")
}

func getManifest(clientFile File) (Manifest, error) {
	if clientFile.Id == 0 {
		return Manifest{}, os.ErrNotExist
	}
	manifestFile, err := extractFile(getCacheFilePath(clientFile), "manifest.json")
	defer pufferpanel.Close(manifestFile)

	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	err = json.NewDecoder(manifestFile).Decode(&manifest)
	return manifest, err
}

func extractFile(zipFile, fileName string) (*os.File, error) {
	folder := filepath.Dir(zipFile)

	file, err := os.Open(filepath.Join(folder, fileName))
	if err != nil && os.IsNotExist(err) {
		err = pufferpanel.ExtractFileFromZip(zipFile, folder, fileName)
		if err != nil {
			return nil, err
		}
		//re-open file
		file, err = os.Open(filepath.Join(folder, fileName))
	}
	return file, err
}

func readVariableFile(serverFile File) (map[string]string, error) {
	varFile, err := extractFile(getCacheFilePath(serverFile), "variables.txt")
	defer pufferpanel.Close(varFile)

	if err != nil {
		return nil, err
	}

	data := make(map[string]string)

	scanner := bufio.NewScanner(varFile)
	var txt string
	for scanner.Scan() {
		txt = scanner.Text()
		parts := strings.SplitN(txt, "=", 2)
		data[parts[0]] = strings.Trim(parts[1], "\"")
	}
	return data, scanner.Err()
}
