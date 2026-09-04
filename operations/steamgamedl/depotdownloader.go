package steamgamedl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v3"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

const DownloadBaseUrl = "https://github.com/SteamRE/DepotDownloader/releases/download/DepotDownloader_${version}/"
const RepoReleases = "https://api.github.com/repos/SteamRE/DepotDownloader/releases?per_page=1"

func downloadDD(rootBinaryFolder string, version string) error {
	downloader.Lock()
	defer downloader.Unlock()

	fi, err := os.Stat(filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary))
	if err == nil && fi.Size() > 0 {
		return nil
	}

	var link string

	if version == "" || version == "latest" {
		link, err = getLatestVersion()
		if err != nil {
			return err
		}
	} else {
		link = strings.ReplaceAll(DownloadBaseUrl+AssetName, "${version}", version)
		arch := "x64"
		if runtime.GOOS == "arm64" {
			arch = "arm64"
		}
		link = strings.Replace(link, "${arch}", arch, 1)
	}

	err = pufferpanel.HttpExtract(link, files.BinaryFS, "depotdownloader", archiver.DefaultZip)
	if err != nil {
		return err
	}

	_ = os.Chmod(filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary), 0755)
	return nil
}

func getLatestVersion() (string, error) {
	client := pufferpanel.Http()
	request := &http.Request{}

	var err error
	request.URL, err = url.Parse(RepoReleases)
	if err != nil {
		return "", err
	}

	request.Header = make(http.Header)
	request.Header.Add("Accept", "application/vnd.github+json")

	response, err := client.Do(request)
	defer utils.CloseResponse(response)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("github api status code: %d", response.StatusCode)
	}

	var data []GithubRelease
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	for _, release := range data {
		for _, asset := range release.Assets {
			assetName := strings.ToLower(asset.Name)
			if strings.Contains(assetName, strings.ToLower(runtime.GOOS)) {
				if runtime.GOARCH == "amd64" {
					if strings.Contains(assetName, "x64") {
						return asset.DownloadUrl, nil
					}
				} else if strings.Contains(assetName, strings.ToLower(runtime.GOARCH)) {
					return asset.DownloadUrl, nil
				}
			}
		}
	}
	return "", errors.New("failed to find latest version for DepotDownloader")
}

type GithubRelease struct {
	Assets []GithubAsset `json:"assets"`
}

type GithubAsset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}
