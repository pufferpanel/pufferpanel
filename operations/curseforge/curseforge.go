package curseforge

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/environments"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

const PageSize = 10

type CurseForge struct {
	ProjectId  uint
	FileId     uint
	JavaBinary string
}

func (c CurseForge) Run(env pufferpanel.Environment) error {
	client := pufferpanel.Http()

	var file *File
	var err error
	if c.FileId == 0 {
		//we need to get the latest file id to do our calls
		files, err := c.getLatestFiles(client, c.ProjectId)
		if err != nil {
			return err
		}

		for _, v := range files {
			if !IsAllowedFile(v.FileStatus) {
				continue
			}
			if file == nil {
				file = &v
				continue
			}
			if file.FileDate.Before(v.FileDate) {
				file = &v
				continue
			}
		}

		if file == nil {
			return errors.New("no files available on CurseForge")
		}
	} else {
		file, err = c.getFileById(client, c.ProjectId, c.FileId)
		if err != nil {
			return err
		}
	}

	if !file.IsServerPack && file.ServerPackFileId != 0 {
		file, err = c.getFileById(client, c.ProjectId, file.ServerPackFileId)
		if err != nil {
			return err
		}
	}

	err = environments.DownloadFile(file.DownloadUrl, "download.zip", env)
	if err != nil {
		return err
	}
	env.DisplayToConsole(true, "Extracting %s", filepath.Join(env.GetRootDirectory(), "download.zip"))
	err = pufferpanel.ExtractZipIgnoreSingleDir(filepath.Join(env.GetRootDirectory(), "download.zip"), env.GetRootDirectory())
	if err != nil {
		return err
	}

	//err = os.Remove(filepath.Join(env.GetRootDirectory(), "download.zip"))
	//if err != nil {
	//	return err
	//}

	//now... set up the mod launcher in the event it's not there
	//use the file to work out what we need to do, because it's going to be a mess of a thing

	entries, err := os.ReadDir(env.GetRootDirectory())
	if err != nil {
		return err
	}

	forgeInstallerRegex := regexp.MustCompile("forge-.*-installer.jar")
	for _, v := range entries {
		if v.IsDir() {
			continue
		}
		if forgeInstallerRegex.MatchString(v.Name()) {
			//forge installer found, we will run this one
			result := make(chan int, 1)
			err = env.Execute(pufferpanel.ExecutionData{
				Command:   c.JavaBinary,
				Arguments: []string{"-jar", v.Name(), "--installServer"},
				Callback: func(exitCode int) {
					result <- exitCode
					env.DisplayToConsole(true, "Installer exit code: %d", exitCode)
				},
			})
			if err != nil {
				return err
			}
			if <-result != 0 {
				return errors.New("failed to run forge installer")
			}
		}
	}

	return nil
}

func (c CurseForge) getLatestFiles(client *http.Client, projectId uint) ([]File, error) {
	u := fmt.Sprintf("https://api.curseforge.com/v1/mods/%d", projectId)

	response, err := c.call(client, u)
	if err != nil {
		return nil, err
	}
	defer pufferpanel.CloseResponse(response)

	if response.StatusCode == 404 {
		return nil, nil
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code from CurseForge: %s", response.Status)
	}

	d, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var addon AddonResponse
	err = json.Unmarshal(d, &addon)
	if err != nil {
		return nil, err
	}
	return addon.Data.LatestFiles, err
}

func (c CurseForge) getFileById(client *http.Client, projectId, fileId uint) (*File, error) {
	u := fmt.Sprintf("https://api.curseforge.com/v1/mods/%d/files/%d", projectId, fileId)

	response, err := c.call(client, u)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return nil, fmt.Errorf("file id %d not found", fileId)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code from CurseForge: %s", response.Status)
	}

	var res FileResponse
	err = json.NewDecoder(response.Body).Decode(&res)
	return &res.Data, err
}

func (c CurseForge) call(client *http.Client, u string) (*http.Response, error) {
	path, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method: "GET",
		URL:    path,
		Header: http.Header{},
	}
	request.Header.Add("x-api-key", config.CurseForgeKey.Value())

	response, err := client.Do(request)
	return response, err
}
