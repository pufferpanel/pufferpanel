package pufferpanel

import (
	"crypto/sha1"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DownloadFile(url, fileName string, env Environment) error {
	target, err := os.Create(filepath.Join(env.GetRootDirectory(), fileName))
	defer Close(target)
	if err != nil {
		return err
	}

	env.DisplayToConsole(true, "Downloading: "+url+"\n")

	response, err := HttpGet(url)
	defer CloseResponse(response)
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
	defer Close(target)
	if err != nil {
		return err
	}

	logging.Info.Printf("Downloading: " + url)

	response, err := HttpGet(url)
	defer CloseResponse(response)
	if err != nil {
		return err
	}

	_, err = io.Copy(target, response.Body)
	return err
}

func DownloadViaMaven(downloadUrl string, env Environment) (string, error) {
	localPath := filepath.Join(config.CacheFolder.Value(), strings.TrimPrefix(strings.TrimPrefix(downloadUrl, "http://"), "https://"))

	if os.PathSeparator != '/' {
		localPath = strings.Replace(localPath, "/", string(os.PathSeparator), -1)
	}

	sha1Url := downloadUrl + ".sha1"

	if env != nil {
		env.DisplayToConsole(true, "Downloading: %s\n", downloadUrl)
	}

	useCache := true
	f, err := os.Open(localPath)
	defer Close(f)
	//cache was readable, so validate
	if err == nil {
		h := sha1.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}
		Close(f)

		actualHash := fmt.Sprintf("%x", h.Sum(nil))

		logging.Info.Printf("Downloading hash from %s", sha1Url)
		response, err := HttpGet(sha1Url)
		defer CloseResponse(response)
		if err != nil {
			useCache = false
		} else {
			data := make([]byte, 40)
			_, err := response.Body.Read(data)
			expectedHash := string(data)

			if err != nil {
				useCache = false
			} else if expectedHash != actualHash {
				logging.Info.Printf("Cache expected %s but was actually %s", expectedHash, actualHash)
				useCache = false
			}
		}

		if useCache {
			if env != nil {
				logging.Info.Printf("Using cached copy of file: %s\n", downloadUrl)
			}
		}
	} else if !os.IsNotExist(err) {
		logging.Info.Printf("Cached file is not readable, will download (%s)", localPath)
	} else {
		useCache = false
	}

	//if we can't use cache, redownload it to the cache
	if !useCache {
		logging.Info.Printf("Downloading new version and caching to %s", localPath)
		err = DownloadFileToCache(downloadUrl, localPath)
	}
	if err == nil {
		return localPath, err
	} else {
		return "", err
	}
}
