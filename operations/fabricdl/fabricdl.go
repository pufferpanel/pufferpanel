package fabricdl

import (
	"encoding/json"
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"path"
)

const FabricMetadataUrl = "https://meta.fabricmc.net/v2/versions/installer"

type Fabricdl struct {
}

type FabricMetadata struct {
	Url string `json:"url"`
}

func (f *Fabricdl) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	env.DisplayToConsole(true, "Downloading metadata from %s\n", FabricMetadataUrl)
	response, err := pufferpanel.HttpGet(FabricMetadataUrl)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	defer pufferpanel.Close(response.Body)

	var metadata []FabricMetadata
	err = json.NewDecoder(response.Body).Decode(&metadata)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	if len(metadata) == 0 {
		err = errors.New("no metadata available from Fabric, unable to download installer")
		return pufferpanel.OperationResult{Error: err}
	}

	file, err := pufferpanel.DownloadViaMaven(metadata[0].Url, env)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	err = pufferpanel.CopyFile(file, path.Join(env.GetRootDirectory(), "fabric-installer.jar"))
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	return pufferpanel.OperationResult{Error: nil}
}
