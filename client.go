package pufferpanel

import (
	"net/http"

	"github.com/cavaliergopher/grab/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/files"
)

var httpClient = &http.Client{}

func Http() *http.Client {
	return httpClient
}

func HttpGet(requestUrl string) (*http.Response, error) {
	return httpClient.Get(requestUrl)
}

func HttpExtract(serverFS files.FileServer, requestUrl, directory string) error {
	//we will write this to temp so we can not keep so much in memory
	response, err := grab.Get(config.CacheFolder.Value(), requestUrl)
	if err != nil {
		return err
	}
	defer files.CacheFS.Remove(response.Filename)

	err = files.Extract(files.CacheFS, response.Filename, serverFS, directory, "*")
	return err
}
