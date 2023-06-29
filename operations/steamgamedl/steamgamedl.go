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
	"errors"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
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

	var args = []string{"-app", c.AppId, "-dir", "."}
	if c.Username != "" {
		args = append(args, "-username", c.Username, "-remember-password")
		if c.Password != "" {
			args = append(args, "-password", c.Password)
		}
	}
	args = append(args, c.ExtraArgs...)

	cmd, err := filepath.Abs(filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary))
	if err != nil {
		return err
	}

	ch := make(chan bool, 1)
	steps := pufferpanel.ExecutionData{
		//Command:          fmt.Sprintf("%s%c%s", ".", filepath.Separator, "dotnet"),
		Command:   cmd,
		Arguments: args,
		Callback: func(exitCode bool) {
			ch <- exitCode
		},
	}
	err = env.Execute(steps)
	if err != nil {
		return err
	}
	success := <-ch
	if !success {
		return errors.New("depotdownloader exited with non-zero code")
	}

	//for some steam games, there's a binary we can instant-mark
	if fi, err := os.Stat(filepath.Join(env.GetRootDirectory(), "srcds_run")); err == nil && !fi.IsDir() {
		_ = os.Chmod(filepath.Join(env.GetRootDirectory(), "srcds_run"), 0755)
	}

	return nil
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
		arch = "amd64"
	}
	link = strings.Replace(link, "${arch}", arch, 1)

	err = pufferpanel.HttpGetZip(link, filepath.Join(rootBinaryFolder, "depotdownloader"))
	if err != nil {
		return err
	}

	err = os.Chmod(filepath.Join(rootBinaryFolder, "depotdownloader", DepotDownloaderBinary), 0755)
	return err
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

	for source, target := range RenameFolders {
		err = os.Rename(filepath.Join(env.GetRootDirectory(), ".steam", source), filepath.Join(env.GetRootDirectory(), ".steam", target))
		if err != nil {
			return err
		}
	}
	return err
}
