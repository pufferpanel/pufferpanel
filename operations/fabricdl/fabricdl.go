package fabricdl

import (
	"encoding/json"
	"errors"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

const FabricMetadataUrl = "https://meta.fabricmc.net/v2/versions/installer"

type Fabricdl struct {
}

type FabricMetadata struct {
	Url string `json:"url"`
}

func (f *Fabricdl) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment

	env.DisplayToConsole(true, "Downloading metadata from %s\n", FabricMetadataUrl)
	response, err := pufferpanel.HttpGet(FabricMetadataUrl)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	defer utils.Close(response.Body)

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
	defer utils.Close(file)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	err = files.WriteFile(args.Server.GetFileServer(), file, "fabric-installer.jar")
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	return pufferpanel.OperationResult{Error: nil}
}
