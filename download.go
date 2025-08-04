package pufferpanel

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DownloadFile(url, fileName string, env *Environment) error {
	target, err := os.Create(filepath.Join(env.GetRootDirectory(), fileName))
	defer utils.Close(target)
	if err != nil {
		return err
	}

	env.DisplayToConsole(true, "Downloading: "+url+"\n")

	response, err := HttpGet(url)
	defer utils.CloseResponse(response)
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
	defer utils.Close(target)
	if err != nil {
		return err
	}

	logging.Info.Printf("Downloading: %s\n", url)

	response, err := HttpGet(url)
	defer utils.CloseResponse(response)
	if err != nil {
		return err
	}

	_, err = io.Copy(target, response.Body)
	return err
}

type FileHash int

const (
	FileHashSHA1 = iota
	FileHashSHA256
)

func FileFromCacheOrDownload(downloadUrl, expectedHash string, algorithm FileHash, env *Environment) (string, error) {
	localPath := filepath.Join(config.CacheFolder.Value(), strings.TrimPrefix(strings.TrimPrefix(downloadUrl, "http://"), "https://"))

	if os.PathSeparator != '/' {
		localPath = strings.Replace(localPath, "/", string(os.PathSeparator), -1)
	}

	if env != nil {
		env.DisplayToConsole(true, "Downloading: %s\n", downloadUrl)
	}

	useCache := true
	f, err := os.Open(localPath)
	defer utils.Close(f)
	//cache was readable, so validate
	if err == nil {
		var h hash.Hash
		if algorithm == FileHashSHA1 {
			h = sha1.New()
		} else if algorithm == FileHashSHA256 {
			h = sha256.New()
		}

		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}
		utils.Close(f)

		actualHash := fmt.Sprintf("%x", h.Sum(nil))

		if strings.HasPrefix(expectedHash, "https://") || strings.HasPrefix(expectedHash, "http://") {
			logging.Info.Printf("Downloading hash from %s", expectedHash)
			response, err := HttpGet(expectedHash)
			defer utils.CloseResponse(response)
			if err != nil {
				useCache = false
			} else {
				data := make([]byte, 40)
				_, err := response.Body.Read(data)
				expectedHash = string(data)

				if err != nil {
					logging.Info.Printf("Failed downloading hash, not using cache")
					useCache = false
				}
			}
		}
		
		if expectedHash != actualHash {
			logging.Info.Printf("Cache expected %s but was actually %s", expectedHash, actualHash)
			useCache = false
		}

		if useCache {
			if env != nil {
				logging.Info.Printf("Using cached copy of file: %s\n", downloadUrl)
			}
		}
	} else {
		logging.Info.Printf("Cached file is not readable, will download (%s)", localPath)
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

func DownloadViaMaven(downloadUrl string, env *Environment) (string, error) {
	return FileFromCacheOrDownload(downloadUrl, downloadUrl + ".sha1", FileHashSHA1, env)
}
