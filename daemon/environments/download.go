/*
 Copyright 2018 Padduck, LLC

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

package environments

import (
	"crypto/sha1"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/daemon/environments/envs"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func DownloadFile(url, fileName string, env envs.Environment) error {
	target, err := os.Create(path.Join(env.GetRootDirectory(), fileName))
	defer pufferpanel.Close(target)
	if err != nil {
		return err
	}

	client := &http.Client{}

	logging.Info().Printf("Downloading: %s", url)
	env.DisplayToConsole(true, "Downloading: "+url+"\n")

	response, err := client.Get(url)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return err
	}

	_, err = io.Copy(target, response.Body)
	return err
}

func DownloadFileToCache(url, fileName string) error {
	parent := filepath.Dir(fileName)
	err := os.MkdirAll(parent, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	target, err := os.Create(fileName)
	defer pufferpanel.Close(target)
	if err != nil {
		return err
	}

	client := &http.Client{}

	logging.Info().Printf("Downloading: " + url)

	response, err := client.Get(url)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return err
	}

	_, err = io.Copy(target, response.Body)
	return err
}

func DownloadViaMaven(downloadUrl string, env envs.Environment) (string, error) {
	localPath := path.Join(viper.GetString("data.cache"), strings.TrimPrefix(strings.TrimPrefix(downloadUrl, "http://"), "https://"))

	if os.PathSeparator != '/' {
		localPath = strings.Replace(localPath, "/", string(os.PathSeparator), -1)
	}

	sha1Url := downloadUrl + ".sha1"

	useCache := true
	f, err := os.Open(localPath)
	defer pufferpanel.Close(f)
	//cache was readable, so validate
	if err == nil {
		h := sha1.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}
		pufferpanel.Close(f)

		actualHash := fmt.Sprintf("%x", h.Sum(nil))

		client := &http.Client{}
		logging.Info().Printf("Downloading hash from %s", sha1Url)
		response, err := client.Get(sha1Url)
		defer pufferpanel.CloseResponse(response)
		if err != nil {
			useCache = false
		} else {
			data := make([]byte, 40)
			_, err := response.Body.Read(data)
			expectedHash := string(data)

			if err != nil {
				useCache = false
			} else if expectedHash != actualHash {
				logging.Info().Printf("Cache expected %s but was actually %s", expectedHash, actualHash)
				useCache = false
			}
		}
	} else if !os.IsNotExist(err) {
		logging.Info().Printf("Cached file is not readable, will download (%s)", localPath)
	} else {
		useCache = false
	}

	//if we can't use cache, redownload it to the cache
	if !useCache {
		logging.Info().Printf("Downloading new version and caching to %s", localPath)
		if env != nil {
			env.DisplayToConsole(true, "Downloading:"+downloadUrl)
		}
		err = DownloadFileToCache(downloadUrl, localPath)
	}
	if err == nil {
		return localPath, err
	} else {
		return "", err
	}
}
