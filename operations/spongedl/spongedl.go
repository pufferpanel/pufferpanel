/*
 Copyright 2021 PufferPanel

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package spongedl

import (
	"encoding/json"
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/environments"
	"github.com/pufferpanel/pufferpanel/v3/operations/forgedl"
	"os"
	"path"
	"strings"
)

const SpongeApiBaseUrl = "https://dl-api-new.spongepowered.org/api/v2/groups/org.spongepowered/artifacts/"

type SpongeDl struct {
	Recommended      bool
	SpongeType       string
	SpongeVersion    string //version of sponge to download
	MinecraftVersion string //version of minecraft to download
}

type SpongeApiV2Versions struct {
	Artifacts map[string]interface{} `json:"artifacts"`
}

type SpongeApiV2Latest struct {
	Assets []SpongeApiV2Asset `json:"assets"`
}

type SpongeApiV2Asset struct {
	Classifier  string
	DownloadUrl string
	Extension   string
}

func (op SpongeDl) Run(env pufferpanel.Environment) error {
	//first, we need to get the build we need to get, if one isn't specified
	if op.SpongeVersion == "" {
		data, err := op.getLatestVersion(env)
		if err != nil {
			return err
		}

		if len(data.Artifacts) == 0 {
			env.DisplayToConsole(true, "No matching Sponge versions found")
			return errors.New("no valid sponge versions found")
		}

		for k, _ := range data.Artifacts {
			op.SpongeVersion = k
			break
		}
	}

	var key string
	if op.SpongeType == "vanilla" {
		key = ""
	} else {
		key = "universal"
	}

	data, err := op.getSpecificVersion(env, op.SpongeVersion)
	if err != nil {
		return err
	}

	var url string
	for _, v := range data.Assets {
		if v.Classifier == key && v.Extension == "jar" {
			url = v.DownloadUrl
		}
	}

	switch strings.ToLower(op.SpongeType) {
	case "spongeforge":
		{
			mapping := make(map[string]interface{})

			var version = ""

			mapping["version"] = version
			mapping["target"] = "forge-installer.jar"
			forgeDlOp, err := forgedl.Factory.Create(pufferpanel.CreateOperation{OperationArgs: mapping})
			if err != nil {
				return err
			}

			err = forgeDlOp.Run(env)
			if err != nil {
				return err
			}

			err = os.Mkdir(path.Join(env.GetRootDirectory(), "mods"), 0755)
			if err != nil && !os.IsExist(err) {
				return err
			}

			file, err := environments.DownloadViaMaven(url, env)
			if err != nil {
				return err
			}

			//going to stick the spongeforge rename in, to assist with those modpacks
			err = pufferpanel.CopyFile(file, path.Join(env.GetRootDirectory(), "mods", "spongeforge.jar"))
			if err != nil {
				return err
			}
		}
		break
	case "spongevanilla":
		{
			file, err := environments.DownloadViaMaven(url, env)
			if err != nil {
				return err
			}

			err = pufferpanel.CopyFile(file, path.Join(env.GetRootDirectory(), "server.jar"))
			if err != nil {
				return err
			}
		}
		break
	default:
		return errors.New("invalid sponge type")
	}

	return nil
}

func (op SpongeDl) getLatestVersion(env pufferpanel.Environment) (SpongeApiV2Versions, error) {
	var data SpongeApiV2Versions

	var params = "?limit=1"
	if op.MinecraftVersion != "" {
		params += "&tags=minecraft:" + op.MinecraftVersion
	}
	if op.Recommended {
		params += "&recommended=true"
	}

	var url = SpongeApiBaseUrl + op.SpongeType + "/versions" + params

	response, err := pufferpanel.HttpGet(url)
	if err != nil {
		return data, err
	}
	defer pufferpanel.CloseResponse(response)
	if response.StatusCode != 200 {
		env.DisplayToConsole(true, "Failed to get the Sponge information from %s: %s", url, response.Status)
		return data, errors.New(response.Status)
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	return data, err
}

func (op SpongeDl) getSpecificVersion(env pufferpanel.Environment, version string) (SpongeApiV2Latest, error) {
	var data SpongeApiV2Latest

	var url = SpongeApiBaseUrl + op.SpongeType + "/versions/" + version

	response, err := pufferpanel.HttpGet(url)
	if err != nil {
		return data, err
	}
	defer pufferpanel.CloseResponse(response)
	if response.StatusCode != 200 {
		env.DisplayToConsole(true, "Failed to get the Sponge information from %s: %s", url, response.Status)
		return data, errors.New(response.Status)
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	return data, err

	return data, nil
}
