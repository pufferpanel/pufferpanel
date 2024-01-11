package curseforge

import (
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"io"
	"net/http"
	"net/url"
)

func getLatestFiles(projectId uint) ([]File, error) {
	u := fmt.Sprintf("https://api.curseforge.com/v1/mods/%d", projectId)

	response, err := callCurseForge(u)
	if err != nil {
		return nil, err
	}
	defer pufferpanel.CloseResponse(response)

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if response.StatusCode != http.StatusOK {
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

func getFileById(projectId, fileId uint) (*File, error) {
	u := fmt.Sprintf("https://api.curseforge.com/v1/mods/%d/files/%d", projectId, fileId)

	response, err := callCurseForge(u)
	if err != nil {
		return nil, err
	}
	defer pufferpanel.CloseResponse(response)

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("file id %d not found", fileId)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code from CurseForge: %s", response.Status)
	}

	var res FileResponse
	err = json.NewDecoder(response.Body).Decode(&res)
	return &res.Data, err
}

func callCurseForge(u string) (*http.Response, error) {
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

	response, err := pufferpanel.Http().Do(request)
	return response, err
}
