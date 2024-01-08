/*
 Copyright 2019 Padduck, LLC

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

package forgedl

import (
	"encoding/json"
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"path"
	"strings"
)

const InstallerUrl = "https://maven.minecraftforge.net/net/minecraftforge/forge/${version}/forge-${version}-installer.jar"
const PromoUrl = "https://files.minecraftforge.net/net/minecraftforge/forge/promotions_slim.json"

type ForgeDl struct {
	Version          string
	Filename         string
	MinecraftVersion string
	OutputVariable   string
}

func (op ForgeDl) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	if op.Version == "" {
		version, err := getLatestForMCVersion(op.MinecraftVersion)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
		op.Version = op.MinecraftVersion + "-" + version
	}

	jarDownload := strings.Replace(InstallerUrl, "${version}", op.Version, -1)

	localFile, err := pufferpanel.DownloadViaMaven(jarDownload, env)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	//copy from the cache
	err = pufferpanel.CopyFile(localFile, path.Join(env.GetRootDirectory(), op.Filename))
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	return pufferpanel.OperationResult{VariableOverrides: map[string]interface{}{
		op.OutputVariable: op.Version,
	}}
}

func getLatestForMCVersion(minecraftVersion string) (string, error) {
	response, err := pufferpanel.HttpGet(PromoUrl)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return "", err
	}

	var promos ForgePromos
	err = json.NewDecoder(response.Body).Decode(&promos)
	if err != nil {
		return "", err
	}
	version := promos.VersionMap[minecraftVersion+"-latest"]
	if version == "" {
		return "", errors.New("no forge available for mc version")
	}
	return version, nil
}

type ForgePromos struct {
	Homepage   string            `json:"homepage"`
	VersionMap map[string]string `json:"promos"`
}
