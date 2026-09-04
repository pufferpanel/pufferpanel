package pufferpanel

import (
	"crypto"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

func DownloadFile(url string, fs files.FileServer, fileName string, env *Environment) error {
	target, err := fs.Create(fileName)
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
	err := files.CacheFS.MkdirAll(parent, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	target, err := files.CacheFS.Create(fileName)
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

func downloadFile(url string) (io.ReadCloser, error) {
	logging.Info.Printf("Downloading: %s", url)
	response, err := HttpGet(url)
	if err != nil {
		return nil, err
	}
	return response.Body, err
}

func cacheFile(downloadUrl, localPath string) (io.ReadCloser, error) {
	dl, err := downloadFile(downloadUrl)
	if err != nil {
		utils.Close(dl)
		return nil, err
	}
	parent := filepath.Dir(localPath)
	err = files.CacheFS.MkdirAll(parent, 0755)
	if err != nil && !os.IsExist(err) {
		logging.Info.Printf("Failed directories in cache: %s", err)
		return dl, nil
	}
	f, err := files.CacheFS.Create(localPath)
	if err != nil {
		utils.Close(f)
		logging.Info.Printf("Failed creating file in cache: %s", err)
		return dl, nil
	}
	_, err = io.Copy(f, dl)
	utils.Close(dl)
	if err != nil {
		// failed actually writing to the successfully created file
		utils.Close(f)
		return nil, err
	}
	err = f.Sync()
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Download(downloadUrl, hash string, algorithm crypto.Hash, cache bool, env *Environment) (io.ReadCloser, error) {
	if env != nil {
		env.DisplayToConsole(true, "Downloading: %s\n", downloadUrl)
	}

	if !cache {
		// don't interact with cache, directly return download response
		return downloadFile(downloadUrl)
	} else {
		// caching allowed
		localPath := filepath.Join(strings.TrimPrefix(strings.TrimPrefix(downloadUrl, "http://"), "https://"))

		if os.PathSeparator != '/' {
			localPath = strings.Replace(localPath, "/", string(os.PathSeparator), -1)
		}

		// try to open existing cached file
		f, err := files.CacheFS.Open(localPath)
		if os.IsNotExist(err) {
			// cache miss, need to download
			return cacheFile(downloadUrl, localPath)
		} else if err != nil {
			logging.Info.Printf("Failed opening cached file despite it existing: %s", err)
			return downloadFile(downloadUrl)
		}

		defer utils.Close(f)

		toReturn, err := files.CacheFS.Open(localPath)

		h := algorithm.New()
		if _, err := io.Copy(h, f); err != nil {
			utils.Close(f)
			logging.Info.Printf("Cached file is not readable, will download (%s)", localPath)
			return downloadFile(downloadUrl)
		}
		actualHash := fmt.Sprintf("%x", h.Sum(nil))
		if hash == actualHash {
			logging.Info.Printf("Using cached copy of file: %s\n", downloadUrl)
			return toReturn, nil
		} else {
			logging.Info.Printf("Cache expected %s but was actually %s, downloading new version and caching to %s", hash, actualHash, localPath)
			utils.Close(toReturn)
			return cacheFile(downloadUrl, localPath)
		}
	}
}

func DownloadHash(hashUrl string, algorithm crypto.Hash) (string, error) {
	logging.Info.Printf("Downloading hash from %s", hashUrl)
	response, err := HttpGet(hashUrl)
	defer utils.CloseResponse(response)
	if err != nil {
		return "", err
	} else {
		data := make([]byte, algorithm.Size()*2)
		_, err := response.Body.Read(data)
		if err != nil {
			return "", err
		}

		return string(data), nil
	}
}

func DownloadViaMaven(downloadUrl string, env *Environment) (io.ReadCloser, error) {
	hashUrl := downloadUrl + ".sha1"
	expectedHash, err := DownloadHash(hashUrl, crypto.SHA1)
	if err != nil {
		logging.Info.Printf("Failed downloading hash, not using cache")
		return Download(downloadUrl, "", crypto.SHA1, false, env)
	}

	return Download(downloadUrl, expectedHash, crypto.SHA1, true, env)
}
