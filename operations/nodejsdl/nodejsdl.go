package nodejsdl

import (
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/hashicorp/go-version"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

var downloader sync.Mutex
var VersionMeta = "https://nodejs.org/dist/index.json"
var DownloadLink = "https://nodejs.org/dist/v${version}/node-v${version}-${os}-${arch}.${ext}"
var VersionSlug = "node-v${version}-${os}-${arch}"

type NodejsDl struct {
	Version string
}

func (op NodejsDl) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment

	env.DisplayToConsole(true, "Downloading Node.js "+op.Version)

	downloader.Lock()
	defer downloader.Unlock()

	mainNodeCommand := "node" + op.Version
	mainNpmCommand := "npm" + op.Version

	_, err := exec.LookPath("node" + op.Version)

	if errors.Is(err, exec.ErrNotFound) {
		var release ReleaseInfo
		release, err = op.getRelease()
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}

		//cleanup the existing dir
		err = files.BinaryFS.RemoveAll(release.Slug)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}

		logging.Debug.Println("Calling " + release.Url)
		err = pufferpanel.HttpExtract(release.Url, files.BinaryFS, "/", nil)

		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}

		_ = files.BinaryFS.Remove(mainNodeCommand)
		_ = files.BinaryFS.Remove(mainNpmCommand)

		logging.Debug.Printf("Adding to path: %s\n", mainNodeCommand)
		err = files.BinaryFS.Symlink(filepath.Join(release.Slug, "bin", "node"), mainNodeCommand)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}

		logging.Debug.Printf("Adding to path: %s\n", mainNpmCommand)
		err = files.BinaryFS.Symlink(filepath.Join(release.Slug, "bin", "npm"), mainNpmCommand)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	return pufferpanel.OperationResult{Error: err}
}

func (op NodejsDl) getRelease() (ReleaseInfo, error) {
	logging.Debug.Println("Calling " + VersionMeta)
	response, err := pufferpanel.HttpGet(VersionMeta)
	defer utils.CloseResponse(response)
	if err != nil {
		return ReleaseInfo{}, err
	}

	var releases []Release
	err = json.NewDecoder(response.Body).Decode(&releases)
	if err != nil {
		return ReleaseInfo{}, err
	}

	var bestMatch, _ = version.NewVersion("0")
	for _, release := range releases {
		if !strings.HasPrefix(release.Version, "v"+op.Version) {
			continue
		}
		if ver, err := version.NewVersion(strings.TrimPrefix(release.Version, "v")); err == nil && bestMatch.LessThan(ver) {
			bestMatch = ver
		} else if err != nil {
			logging.Info.Printf("failed to parse version '%s', %s", release, err)
		}
	}

	replacements := map[string]interface{}{
		"version": bestMatch.Original(),
	}
	if runtime.GOOS == "windows" {
		replacements["os"] = "windows"
		replacements["ext"] = "zip"
	} else {
		replacements["os"] = "linux"
		replacements["ext"] = "tar.xz"
	}

	switch runtime.GOARCH {
	case "arm64":
		{
			replacements["arch"] = "arm64"
		}
	case "arm":
		{
			replacements["arch"] = "armv7l"
		}
	default:
		{
			replacements["arch"] = "x64"
		}
	}

	release := ReleaseInfo{
		Url:  utils.ReplaceTokens(DownloadLink, replacements, utils.PlainReplace),
		Slug: utils.ReplaceTokens(VersionSlug, replacements, utils.PlainReplace),
	}
	return release, nil
}

type Release struct {
	Version string `json:"version"`
}

type ReleaseInfo struct {
	Url  string
	Slug string
}
