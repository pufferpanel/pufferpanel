package neoforgedl

import (
	"encoding/xml"
	"errors"
	"github.com/hashicorp/go-version"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"path"
	"strings"
)

const InstallerUrl = "https://maven.neoforged.net/releases/net/neoforged/neoforge/${version}/neoforge-${version}-installer.jar"
const MetadataUrl = "https://maven.neoforged.net/releases/net/neoforged/neoforge/maven-metadata.xml"

type NeoforgeDL struct {
	Version          string
	Filename         string
	MinecraftVersion string
	OutputVariable   string
}

func (op NeoforgeDL) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	if op.Version == "" {
		neoVersion, err := getLatestForMCVersion(op.MinecraftVersion)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
		op.Version = neoVersion
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
	response, err := pufferpanel.HttpGet(MetadataUrl)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return "", err
	}

	var metadata Metadata
	err = xml.NewDecoder(response.Body).Decode(&metadata)
	if err != nil {
		return "", err
	}
	splitVersion := strings.TrimPrefix(minecraftVersion, "1.")

	var topVersion *version.Version

	for _, v := range metadata.Versions {
		if strings.HasPrefix(v, splitVersion) {
			newVersion, err := version.NewVersion(v)
			if err != nil {
				logging.Debug.Printf("Failed to parse version for Neoforge: %s -> %s", v, err.Error())
				continue
			}
			if topVersion == nil {
				topVersion = newVersion
			} else if newVersion.GreaterThan(topVersion) {
				topVersion = newVersion
			}
		}
	}

	if topVersion == nil {
		return "", errors.New("failed to find neoforge version for " + minecraftVersion)
	}

	return topVersion.Original(), nil
}

type Metadata struct {
	Versions []string `xml:"versioning>versions>version"`
	Latest   string   `xml:"latest"`
	Release  string   `xml:"release"`
}
