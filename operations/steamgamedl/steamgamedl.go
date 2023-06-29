/*
 Copyright 2022 PufferPanel
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

package steamgamedl

import (
	"bufio"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/spf13/cast"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var downloader sync.Mutex

const DepotDownloaderLink = "https://github.com/SteamRE/DepotDownloader/releases/download/DepotDownloader_2.4.7/depotdownloader-2.4.7.zip"
const SteamMetadataServerLink = "https://media.steampowered.com/client/"
const SteamMetadataLinuxLink = SteamMetadataServerLink + "steam_cmd_linux"

type SteamGameDl struct {
	AppId     string
	Username  string
	Password  string
	ExtraArgs []string
}

func (c SteamGameDl) Run(env pufferpanel.Environment) (err error) {
	env.DisplayToConsole(true, "Downloading game from Steam")
	rootBinaryFolder := config.BinariesFolder.Value()

	err = downloadBinaries(rootBinaryFolder)
	if err != nil {
		return err
	}

	err = downloadMetadata(env)
	if err != nil {
		return err
	}

	//generate a login id
	//this is a 32-bit id, which Steam derives from private IP
	//as such, we can kinda send anything we want
	//our approach will be we hash the server id
	loginId := cast.ToString(rand.Int31())

	manifestFolder := filepath.Join(env.GetRootDirectory(), ".manifest")
	_ = os.RemoveAll(manifestFolder)

	args := []string{filepath.Join(rootBinaryFolder, "depotdownloader", "DepotDownloader.dll"), "-app", c.AppId, "-dir", manifestFolder, "-loginid", loginId, "-manifest-only"}
	if c.Username != "" {
		args = append(args, "-username", c.Username, "-remember-password")
		if c.Password != "" {
			args = append(args, "-password", c.Password)
		}
	}

	ch := make(chan int, 1)
	steps := pufferpanel.ExecutionData{
		Command:   filepath.Join(rootBinaryFolder, "dotnet-runtime", "dotnet"),
		Arguments: args,
		Callback: func(exitCode int) {
			ch <- exitCode
		},
	}
	err = env.Execute(steps)
	if err != nil {
		return err
	}
	exitCode := <-ch
	if exitCode != 0 {
		return fmt.Errorf("depotdownloader exited with non-zero code %d", exitCode)
	}

	//download game itself now
	args = []string{filepath.Join(rootBinaryFolder, "depotdownloader", "DepotDownloader.dll"), "-app", c.AppId, "-dir", env.GetRootDirectory(), "-loginid", loginId, "-validate"}
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
		Command:   filepath.Join(rootBinaryFolder, "dotnet-runtime", "dotnet"),
		Arguments: args,
		Callback: func(exitCode int) {
			ch <- exitCode
		},
	}
	err = env.Execute(steps)
	if err != nil {
		return err
	}
	exitCode = <-ch
	if exitCode != 0 {
		return fmt.Errorf("depotdownloader exited with non-zero code %d", exitCode)
	}

	//for each file we download, we need to just... chmod +x the files
	//we rely on the manifests for this
	manifests, err := os.ReadDir(manifestFolder)
	if err != nil {
		return err
	}
	for _, manifest := range manifests {
		if manifest.Type().IsDir() || !strings.HasSuffix(manifest.Name(), ".txt") {
			continue
		}
		err = walkManifest(env.GetRootDirectory(), manifest.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadBinaries(rootBinaryFolder string) error {
	downloader.Lock()
	defer downloader.Unlock()

	fi, err := os.Stat(filepath.Join(rootBinaryFolder, "depotdownloader", "DepotDownloader.dll"))
	if err == nil && fi.Size() > 0 {
		return nil
	}

	err = downloadDotNet(rootBinaryFolder)
	if err != nil {
		return err
	}

	cmd := getDotNetInstallCommand()
	cmd.Dir = rootBinaryFolder

	err = cmd.Run()
	if err != nil {
		return err
	}

	if !cmd.ProcessState.Success() {
		out, _ := cmd.CombinedOutput()
		logging.Debug.Println(string(out))
		return fmt.Errorf("dotnet-install exited with non-zero code: %d", cmd.ProcessState.ExitCode())
	}

	_ = os.Remove(filepath.Join(rootBinaryFolder, DotNetScriptName))

	err = pufferpanel.HttpGetZip(DepotDownloaderLink, filepath.Join(rootBinaryFolder, "depotdownloader"))
	if err != nil {
		return err
	}

	return nil
}

func downloadDotNet(targetFolder string) error {
	target, err := os.Create(path.Join(targetFolder, DotNetScriptName))
	defer pufferpanel.Close(target)
	if err != nil {
		return err
	}

	response, err := pufferpanel.HttpGet(DotNetScriptDl)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return err
	}

	_, err = io.Copy(target, response.Body)
	return err
}

func downloadMetadata(env pufferpanel.Environment) error {
	if runtime.GOOS != "linux" {
		env.DisplayToConsole(true, "Node does not support linux, cannot download steam libraries. While this "+
			"will not stop the installation, the server may not fully work")
		return nil
	}

	if runtime.GOARCH != "amd64" {
		env.DisplayToConsole(true, "Node does not support amd64, cannot download steam libraries. While this "+
			"will not stop the installation, the server may not fully work")
		return nil
	}

	downloadOs := "linux"
	link := SteamMetadataLinuxLink

	response, err := pufferpanel.HttpGet(link)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return err
	}

	metadataName, err := Parse(downloadOs, response.Body)
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

	err = os.Rename(filepath.Join(env.GetRootDirectory(), ".steam", "linux32"), filepath.Join(env.GetRootDirectory(), ".steam", "sdk32"))
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(env.GetRootDirectory(), ".steam", "linux64"), filepath.Join(env.GetRootDirectory(), ".steam", "sdk64"))
	return err
}

func walkManifest(folder, filename string) error {
	file, err := os.Open(filename)
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
		if len(parts) > 5 {
			//the filename at the end has spaces, we need to consolidate
			parts[4] = strings.Join(parts[5:], " ")
			parts = parts[0:4]
		}

		//we will only work on 0 files, because this mean no other flags were told
		if parts[3] == "0" {
			fileToUpdate := parts[4]
			_ = os.Chmod(filepath.Join(folder, fileToUpdate), 0755)
		}
	}

	return nil
}
