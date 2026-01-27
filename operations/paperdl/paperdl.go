package paperdl

import (
	"crypto"
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-version"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const VersionsUrl = "https://fill.papermc.io/v3/projects/${project}/versions"
const BuildUrl = "https://fill.papermc.io/v3/projects/${project}/versions/${mcVersion}/builds/${build}"
var UserAgent = pufferpanel.Display + " https://github.com/pufferpanel/pufferpanel"

type PaperDl struct {
	MinecraftVersion string
	Build            string
	Filename         string
	Project          string
}

func (op PaperDl) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment

	if op.MinecraftVersion == "latest" {
		logging.Info.Printf("PaperDL got Minecraft version 'latest', looking up latest version supported by Paper")
		mcVersion, err := op.getLatestMCVersion()
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
		op.MinecraftVersion = mcVersion
	}

	dlUrl, hash, err := op.getDownloadUrlAndHash(env)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	dl, err := pufferpanel.Download(dlUrl, hash, crypto.SHA256, true, env)
	defer utils.Close(dl)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	err = files.WriteFile(dl, path.Join(env.GetRootDirectory(), op.Filename))
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	return pufferpanel.OperationResult{Error: nil}
}

func (op PaperDl) getLatestMCVersion() (string, error) {
	path, err := url.Parse(strings.Replace(VersionsUrl, "${project}", op.Project, -1))
	if err != nil {
		return "", err
	}

	request := &http.Request{
		Method: "GET",
		URL: path,
		Header: http.Header{},
	}
	request.Header.Add("user-agent", UserAgent)

	response, err := pufferpanel.Http().Do(request)
	defer utils.CloseResponse(response)
	if err != nil {
		return "", err
	}

	var versions PaperVersionsResponse
	err = json.NewDecoder(response.Body).Decode(&versions)
	if err != nil {
		return "", err
	}

	latest, _ := version.NewVersion("0.0")
	for _, v := range versions.Versions {
		if ver, err := version.NewVersion(v.VersionInfo.Id); err == nil && latest.LessThan(ver) {
			latest = ver
		} else if err != nil {
			logging.Info.Printf("failed to parse version '%s', %s", v, err)
		}
	}

	logging.Info.Printf("Latest Minecraft version supported by Paper is %s", latest.Original())
	return latest.Original(), nil
}

func (op PaperDl) getDownloadUrlAndHash(env *pufferpanel.Environment) (string, string, error) {
	buildUrl := strings.Replace(BuildUrl, "${mcVersion}", op.MinecraftVersion, -1)
	buildUrl = strings.Replace(buildUrl, "${build}", op.Build, -1)
	buildUrl = strings.Replace(buildUrl, "${project}", op.Project, -1)
	path, err := url.Parse(buildUrl)
	if err != nil {
		return "", "", err
	}

	request := &http.Request{
		Method: "GET",
		URL: path,
		Header: http.Header{},
	}
	request.Header.Add("User-Agent", UserAgent)

	response, err := pufferpanel.Http().Do(request)
	defer utils.CloseResponse(response)
	if err != nil {
		return "", "", err
	}

	if response.StatusCode == 404 {
		env.DisplayToConsole(true, "Invalid Minecraft version or Paper build\n")
		return "", "", errors.New("Invalid minecraft version or paper build")
	}

	var build PaperBuild
	err = json.NewDecoder(response.Body).Decode(&build)
	if err != nil {
		return "", "", err
	}

	return build.Downloads.Server.Url, build.Downloads.Server.Checksums.Sha256, nil
}

type PaperVersionInfo struct {
	Id string `json:"id"`
}

type PaperVersion struct {
	VersionInfo PaperVersionInfo `json:"version"`
}

type PaperVersionsResponse struct {
	Versions []PaperVersion `json:"versions"`
}

type PaperChecksums struct {
	Sha256 string `json:"sha256"`
}

type PaperServer struct {
	Checksums PaperChecksums `json:"checksums"`
	Url       string         `json:"url"`
}

type PaperDownload struct {
	Server PaperServer `json:"server:default"`
}

type PaperBuild struct {
	Downloads PaperDownload `json:"downloads"`
}
