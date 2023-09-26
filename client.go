package pufferpanel

import (
	"io"
	"net/http"
	"os"
)

var httpClient = &http.Client{}

func Http() *http.Client {
	return httpClient
}

func HttpGet(url string) (*http.Response, error) {
	return httpClient.Get(url)
}

func HttpGetTarGz(url, directory string) error {
	response, err := HttpGet(url)
	defer CloseResponse(response)
	if err != nil {
		return err
	}

	err = ExtractTarGz(response.Body, directory)
	return err
}

func HttpGetZip(url, directory string) error {
	//we will write this to temp so we can not keep so much in memory
	file, err := os.CreateTemp("", "pufferpanel-dl-*")
	if err != nil {
		return err
	}

	defer os.Remove(file.Name())

	response, err := HttpGet(url)
	defer CloseResponse(response)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return ExtractZip(file.Name(), directory)
}
