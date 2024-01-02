package mojangdl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

const VersionJsonUrl = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type MojangDl struct {
	Version string
	Target  string
}

func (op MojangDl) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	response, err := pufferpanel.HttpGet(VersionJsonUrl)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	var data LauncherJson
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	err = response.Body.Close()
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	var targetVersion string
	switch op.Version {
	case "release":
		targetVersion = data.Latest.Release
	case "latest":
		targetVersion = data.Latest.Release
	case "snapshot":
		targetVersion = data.Latest.Snapshot
	default:
		targetVersion = op.Version
	}

	for _, version := range data.Versions {
		if version.Id == targetVersion {
			logging.Info.Printf("Version %s json located, downloading from %s", version.Id, version.Url)
			env.DisplayToConsole(true, fmt.Sprintf("Version %s json located, downloading from %s\n", version.Id, version.Url))
			//now, get the version json for this one...
			err = downloadServerFromJson(version.Url, op.Target, env)
			return pufferpanel.OperationResult{Error: err}
		}
	}

	env.DisplayToConsole(true, "Could not locate version "+targetVersion+"\n")
	err = errors.New("Version not located: " + op.Version)
	return pufferpanel.OperationResult{Error: err}
}

func downloadServerFromJson(url, target string, env pufferpanel.Environment) error {
	response, err := pufferpanel.HttpGet(url)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return err
	}

	var data VersionJson
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return err
	}
	err = response.Body.Close()
	if err != nil {
		return err
	}

	serverBlock := data.Downloads["server"]

	logging.Info.Printf("Version jar located, downloading from %s", serverBlock.Url)
	env.DisplayToConsole(true, fmt.Sprintf("Version jar located, downloading from %s\n", serverBlock.Url))

	return pufferpanel.DownloadFile(serverBlock.Url, target, env)
}

type LauncherJson struct {
	Versions []LauncherVersion `json:"versions"`
	Latest   Latest            `json:"latest"`
}

type Latest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type LauncherVersion struct {
	Id   string `json:"id"`
	Url  string `json:"url"`
	Type string `json:"type"`
}

type VersionJson struct {
	Downloads map[string]DownloadType `json:"downloads"`
}

type DownloadType struct {
	Sha1 string `json:"sha1"`
	Size uint64 `json:"size"`
	Url  string `json:"url"`
}
