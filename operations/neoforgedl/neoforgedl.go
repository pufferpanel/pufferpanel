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

package neoforgedl

import (
	"encoding/json"
  "encoding/xml"
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"path"
	"strings"
)

const InstallerUrl = "https://maven.neoforged.net/releases/net/neoforged/neoforge/${version}/neoforge-${version}-installer.jar"
const PromoUrl = "https://maven.neoforged.net/releases/net/neoforged/forge/maven-metadata.xml"

// const PromoUrl = "https://files.minecraftforge.net/net/minecraftforge/forge/promotions_slim.json"

type ForgeDl struct {
	Version          string
	Filename         string
	MinecraftVersion string
	OutputVariable   string
}

func (op ForgeDl) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment

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

  /*
	var promos ForgePromos
	err = json.NewDecoder(response.Body).Decode(&promos)
	if err != nil {
		return "", err
	}
	version := promos.VersionMap[minecraftVersion+"-latest"]
	if version == "" {
		return "", errors.New("no neoforge available for mc version")
	}
	return version, nil
}

type ForgePromos struct {
	Homepage   string            `json:"homepage"`
	VersionMap map[string]string `json:"promos"`
}
*/

func getLatestForMCVersion(minecraftVersion string) (string, error) {
	response, err := http.Get(PromoUrl)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch data: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var metadata Metadata
	err = xml.Unmarshal(body, &metadata)
	if err != nil {
		return "", err
	}

	version := getVersionFromMetadata(metadata, minecraftVersion)
	if version == "" {
		return "", errors.New("no neoforge available for mc version")
	}
	return version, nil
}

func getVersionFromMetadata(metadata Metadata, minecraftVersion string) string {
	for _, version := range metadata.Versioning.Versions.Version {
		if version == minecraftVersion+"-latest" {
			return version
		}
	}
	return ""
}

type Metadata struct {
	XMLName    xml.Name   `xml:"metadata"`
	GroupId    string     `xml:"groupId"`
	ArtifactId string     `xml:"artifactId"`
	Versioning Versioning `xml:"versioning"`
}

type Versioning struct {
	Latest     string    `xml:"latest"`
	Release    string    `xml:"release"`
	Versions   Versions  `xml:"versions"`
	LastUpdated string   `xml:"lastUpdated"`
}

type Versions struct {
	Version []string `xml:"version"`
}
