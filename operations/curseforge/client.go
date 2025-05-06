package curseforge

import (
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"io"
	"net/http"
	"net/url"
)

func getAddonData(projectId uint) (AddonResponse, error) {
	u := fmt.Sprintf("https://api.curseforge.com/v1/mods/%d", projectId)

	response, err := callCurseForge(u)
	if err != nil {
		return AddonResponse{}, err
	}
	defer utils.CloseResponse(response)

	if response.StatusCode == http.StatusNotFound {
		return AddonResponse{}, nil
	}

	if response.StatusCode != http.StatusOK {
		return AddonResponse{}, pufferpanel.ErrCurseForgeStatus(response.Status)
	}

	d, err := io.ReadAll(response.Body)
	if err != nil {
		return AddonResponse{}, err
	}

	var addon AddonResponse
	err = json.Unmarshal(d, &addon)
	if err != nil {
		return AddonResponse{}, err
	}
	return addon, nil
}

func getAddonFileData(projectId uint, fileId uint) (FileResponse, error) {
	u := fmt.Sprintf("https://api.curseforge.com/v1/mods/%d/files/%d", projectId, fileId)

	response, err := callCurseForge(u)
	if err != nil {
		return FileResponse{}, err
	}
	defer utils.CloseResponse(response)

	if response.StatusCode == http.StatusNotFound {
		return FileResponse{}, pufferpanel.ErrCurseForgeFile(projectId, fileId)
	}

	if response.StatusCode != http.StatusOK {
		return FileResponse{}, pufferpanel.ErrCurseForgeStatus(response.Status)
	}

	var res FileResponse
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return FileResponse{}, err
	}
	return res, nil
}

func getLatestFiles(projectId uint) ([]File, error) {
	addon, err := getAddonData(projectId)
	if err != nil {
		return nil, err
	}

	if !addon.Data.AllowModDistribution {
		return nil, pufferpanel.ErrCurseForgeDistribution(projectId)
	}
				       
	return addon.Data.LatestFiles, err
}

func getFileById(projectId uint, fileId uint) (File, error) {
	addon, addonErr := getAddonData(projectId, fileId)

	if addonErr != nil {
		return File{}, addonErr
	}

	if !addon.Data.AllowModDistribution {
		return File{}, pufferpanel.ErrCurseForgeDistribution(projectId)
	}
	
	file, fileErr := getAddonFileData(projectId, fileId)

	if fileErr != nil {
		return File{}, fileErr
	}
				       
	return file.Data, nil
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

	logging.Debug.Printf("Calling %s\n", request.URL.String())
	response, err := pufferpanel.Http().Do(request)
	return response, err
}
