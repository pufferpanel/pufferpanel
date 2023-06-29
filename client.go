/*
 Copyright 2022 PufferPanel
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 	http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package pufferpanel

import (
	"bytes"
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

func HttpDownloadDeb(link, folder string) error {
	response, err := HttpGet(link)
	defer CloseResponse(response)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buff, response.Body)
	CloseResponse(response)

	if err != nil {
		return err
	}

	reader := bytes.NewReader(buff.Bytes())

	err = ExtractDeb(reader, folder)
	return err
}
