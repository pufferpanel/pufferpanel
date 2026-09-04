package curseforge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

func downloadModpack(file File) error {
	cacheFS := files.CacheFS

	var cacheZipFolder = getCacheFolderForFile(file)
	var cacheZipFileLocation = getCacheFilePath(file)

	//see if the file already exists, if so, use it instead
	if fi, err := cacheFS.Lstat(cacheZipFileLocation); err == nil && !fi.IsDir() && fi.Size() > 0 {
		return nil
	}

	err := cacheFS.MkdirAll(cacheZipFolder, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	tmpFile, err := cacheFS.CreateTemp(cacheZipFolder, "tmp-*.zip")
	if err != nil {
		return err
	}
	defer utils.Close(tmpFile)

	response, err := callCurseForge(file.DownloadUrl)
	defer utils.CloseResponse(response)
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
	utils.Close(tmpFile)
	utils.CloseResponse(response)

	err = cacheFS.Rename(tmpFile.Name(), cacheZipFileLocation)
	if err != nil {
		return err
	}

	return nil
}

func getCacheFolderForFile(file File) string {
	return filepath.Join("curseforge", fmt.Sprintf("%d", file.Id))
}

func getCacheFilePath(file File) string {
	return filepath.Join(getCacheFolderForFile(file), "download.zip")
}

func getManifest(clientFile File) (Manifest, error) {
	if clientFile.Id == 0 {
		return Manifest{}, os.ErrNotExist
	}
	manifestFile, err := extractFile(getCacheFilePath(clientFile), "manifest.json")
	defer utils.Close(manifestFile)

	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	err = json.NewDecoder(manifestFile).Decode(&manifest)
	return manifest, err
}

func extractFile(zipFile, fileName string) (fs.File, error) {
	folder := filepath.Dir(zipFile)

	file, err := files.CacheFS.Open(filepath.Join(folder, fileName))
	if err != nil && os.IsNotExist(err) {
		err = files.Extract(files.CacheFS, zipFile, files.CacheFS, folder, fileName, false, nil)
		if err != nil {
			return nil, err
		}
		//re-open file
		file, err = files.CacheFS.Open(filepath.Join(folder, fileName))
	}
	return file, err
}

func readVariableFile(serverFile File) (map[string]string, error) {
	varFile, err := extractFile(getCacheFilePath(serverFile), "variables.txt")
	defer utils.Close(varFile)

	if err != nil {
		return nil, err
	}

	data := make(map[string]string)

	scanner := bufio.NewScanner(varFile)
	var txt string
	for scanner.Scan() {
		txt = scanner.Text()
		if strings.HasPrefix(txt, "#") {
			continue
		}
		parts := strings.SplitN(txt, "=", 2)
		if len(parts) != 2 {
			continue
		}
		data[parts[0]] = strings.Trim(parts[1], "\"")
	}
	return data, scanner.Err()
}
