package steamgamedl

import (
	"bufio"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"github.com/spf13/cast"
)

var downloader sync.Mutex

const SteamMetadataServerLink = "https://media.steampowered.com/client/"

func init() {
}

type SteamGameDl struct {
	AppId          string
	Username       string
	Password       string
	Branch         string
	BranchPassword string
	ExtraArgs      []string
}

func (c SteamGameDl) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment
	serverFs := args.Server.GetFileServer()

	env.DisplayToConsole(true, "Downloading game from Steam")

	rootBinaryFolder := config.BinariesFolder.Value()

	err := downloadDD(config.DepotDownloaderVersion.Value())
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	err = downloadMetadata(serverFs, env)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	//generate a login id
	//this is a 32-bit id, which Steam derives from private IP
	//as such, we can kinda send anything we want
	//our approach will be we hash the server id
	loginId := cast.ToString(rand.Int31())

	manifestFolder := ".manifest"
	_ = serverFs.RemoveAll(manifestFolder)

	cmdArgs := []string{filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary), "-app", c.AppId, "-dir", ".manifest", "-loginid", loginId, "-manifest-only"}
	if c.Username != "" {
		cmdArgs = append(cmdArgs, "-username", c.Username, "-remember-password")
		if c.Password != "" {
			cmdArgs = append(cmdArgs, "-password", c.Password)
		}
	}
	if !config.DepotDownloaderDisableLancache.Value() {
		cmdArgs = append(cmdArgs, "-use-lancache")
	}

	if c.Branch != "" {
		cmdArgs = append(cmdArgs, "-branch", c.Branch)
		if c.BranchPassword != "" {
			cmdArgs = append(cmdArgs, "-branchpassword", c.BranchPassword)
		}
	}

	cmdArgs = append(cmdArgs, c.ExtraArgs...)

	ch := make(chan int, 1)
	steps := pufferpanel.ExecutionData{
		Command: utils.MergeArguments(cmdArgs),
		Callback: func(exitCode int) {
			ch <- exitCode
		},
	}
	err = env.Execute(steps)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	exitCode := <-ch
	if exitCode != 0 {
		err = fmt.Errorf("depotdownloader exited with non-zero code %d", exitCode)
		return pufferpanel.OperationResult{Error: err}
	}

	//download game itself now
	cmdArgs = []string{filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary), "-app", c.AppId, "-dir", ".", "-loginid", loginId, "-validate"}
	if c.Username != "" {
		cmdArgs = append(cmdArgs, "-username", c.Username, "-remember-password")
		if c.Password != "" {
			cmdArgs = append(cmdArgs, "-password", c.Password)
		}
	}

	if !config.DepotDownloaderDisableLancache.Value() {
		cmdArgs = append(cmdArgs, "-use-lancache")
	}

	if len(c.ExtraArgs) > 0 {
		cmdArgs = append(cmdArgs, c.ExtraArgs...)
	}

	steps = pufferpanel.ExecutionData{
		Command: utils.MergeArguments(cmdArgs),
		Callback: func(exitCode int) {
			ch <- exitCode
		},
		Environment: map[string]string{
			"TERM": "pufferpanel", //we use a fake TERM because DD will use a display that is not supported by us directly
		},
	}
	err = env.Execute(steps)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	exitCode = <-ch
	if exitCode != 0 {
		err = fmt.Errorf("depotdownloader exited with non-zero code %d", exitCode)
		return pufferpanel.OperationResult{Error: err}
	}

	//for each file we download, we need to just... chmod +x the files
	//we rely on the manifests for this
	manifests, err := args.Server.GetFileServer().ReadDir(manifestFolder)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	for _, manifest := range manifests {
		if manifest.Type().IsDir() || !strings.HasSuffix(manifest.Name(), ".txt") {
			continue
		}
		err = walkManifest(args.Server.GetFileServer(), manifest.Name())
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	return pufferpanel.OperationResult{Error: nil}
}

func downloadMetadata(serverFs files.FileServer, env *pufferpanel.Environment) error {
	response, err := pufferpanel.HttpGet(SteamMetadataLink)
	defer utils.CloseResponse(response)
	if err != nil {
		return err
	}

	metadataName, err := Parse(DownloadOs, response.Body)
	utils.CloseResponse(response)

	if err != nil {
		return err
	}

	err = serverFs.RemoveAll(".steam")
	if err != nil {
		return err
	}

	err = pufferpanel.HttpExtract(serverFs, SteamMetadataServerLink+metadataName, ".steam")
	if err != nil {
		return err
	}

	return err
}

func walkManifest(fs files.FileServer, filename string) error {
	file, err := fs.Open(filepath.Join(".manifest", filename))
	defer utils.Close(file)
	if err != nil {
		return err
	}
	data := bufio.NewScanner(file)
	skipCounter := 8
	for data.Scan() {
		line := data.Text()
		if skipCounter > 0 {
			skipCounter--
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 5 || parts[0] == "Size" {
			continue
		}
		if len(parts) > 5 {
			//the filename at the end has spaces, we need to consolidate
			parts[4] = strings.Join(parts[5:], " ")
			parts = parts[0:5]
		}

		//we will only work on 0 files, because this mean no other flags were told
		if parts[3] == "0" {
			fileToUpdate := parts[4]
			_ = fs.Chmod(filepath.Join(fileToUpdate), 0755)
		}
	}

	return nil
}
