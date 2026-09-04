package pufferpanel

import (
	"net/http"
	"path/filepath"

	"github.com/cavaliergopher/grab/v3"
	"github.com/mholt/archiver/v3"
	"github.com/pufferpanel/pufferpanel/v3/files"
)

var httpClient = &http.Client{}

func Http() *http.Client {
	return httpClient
}

func HttpGet(requestUrl string) (*http.Response, error) {
	return httpClient.Get(requestUrl)
}

func HttpExtract(requestUrl string, fs files.FileServer, directory string, archiveType archiver.Walker) error {
	//we will write this to temp so we can not keep so much in memory
	response, err := grab.Get(fs.Prefix(), requestUrl)
	if err != nil {
		return err
	}

	file := filepath.Base(response.Filename)

	defer fs.Remove(file)

	err = files.Extract(fs, file, fs, directory, "*", false, archiveType)
	return err
}
