package steamgamedl

import (
	"bufio"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/spf13/cast"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var downloader sync.Mutex

const SteamMetadataServerLink = "https://media.steampowered.com/client/"

func init() {
}

type SteamGameDl struct {
	AppId     string
	Username  string
	Password  string
	ExtraArgs []string
}

func (c SteamGameDl) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	env.DisplayToConsole(true, "Downloading game from Steam")

	rootBinaryFolder := config.BinariesFolder.Value()

	err := downloadBinaries(rootBinaryFolder)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	err = downloadMetadata(env)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	//generate a login id
	//this is a 32-bit id, which Steam derives from private IP
	//as such, we can kinda send anything we want
	//our approach will be we hash the server id
	loginId := cast.ToString(rand.Int31())

	manifestFolder := filepath.Join(env.GetRootDirectory(), ".manifest")
	_ = os.RemoveAll(manifestFolder)

	args := []string{"-app", c.AppId, "-dir", manifestFolder, "-loginid", loginId, "-manifest-only"}
	if c.Username != "" {
		args = append(args, "-username", c.Username, "-remember-password")
		if c.Password != "" {
			args = append(args, "-password", c.Password)
		}
	}
	args = append(args, c.ExtraArgs...)

	ch := make(chan int, 1)
	steps := pufferpanel.ExecutionData{
		Command:   filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary),
		Arguments: args,
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
	args = []string{"-app", c.AppId, "-dir", env.GetRootDirectory(), "-loginid", loginId, "-validate"}
	if c.Username != "" {
		args = append(args, "-username", c.Username, "-remember-password")
		if c.Password != "" {
			args = append(args, "-password", c.Password)
		}
	}

	if c.ExtraArgs != nil && len(c.ExtraArgs) > 0 {
		args = append(args, c.ExtraArgs...)
	}

	steps = pufferpanel.ExecutionData{
		Command:   filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary),
		Arguments: args,
		Callback: func(exitCode int) {
			ch <- exitCode
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
	manifests, err := os.ReadDir(manifestFolder)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	for _, manifest := range manifests {
		if manifest.Type().IsDir() || !strings.HasSuffix(manifest.Name(), ".txt") {
			continue
		}
		err = walkManifest(env.GetRootDirectory(), manifest.Name())
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	return pufferpanel.OperationResult{Error: nil}
}

func downloadBinaries(rootBinaryFolder string) error {
	downloader.Lock()
	defer downloader.Unlock()

	fi, err := os.Stat(filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary))
	if err == nil && fi.Size() > 0 {
		return nil
	}

	link := DepotDownloaderLink
	arch := "x64"
	if runtime.GOOS == "arm64" {
		arch = "arm64"
	}
	link = strings.Replace(link, "${arch}", arch, 1)

	err = pufferpanel.HttpGetZip(link, filepath.Join(rootBinaryFolder, "depotdownloader"))
	if err != nil {
		return err
	}

	_ = os.Chmod(filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary), 0755)
	return nil
}

func downloadMetadata(env pufferpanel.Environment) error {
	response, err := pufferpanel.HttpGet(SteamMetadataLink)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return err
	}

	metadataName, err := Parse(DownloadOs, response.Body)
	pufferpanel.CloseResponse(response)

	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(env.GetRootDirectory(), ".steam"))
	if err != nil {
		return err
	}

	err = pufferpanel.HttpGetZip(SteamMetadataServerLink+metadataName, filepath.Join(env.GetRootDirectory(), ".steam"))
	if err != nil {
		return err
	}

	return err
}

func walkManifest(folder, filename string) error {
	file, err := os.Open(filepath.Join(folder, ".manifest", filename))
	defer pufferpanel.Close(file)
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
			_ = os.Chmod(filepath.Join(folder, fileToUpdate), 0755)
		}
	}

	return nil
}
