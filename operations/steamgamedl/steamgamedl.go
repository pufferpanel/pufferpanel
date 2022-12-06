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
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var downloader sync.Mutex

const DepotDownloaderLink = "https://github.com/SteamRE/DepotDownloader/releases/download/DepotDownloader_2.4.7/depotdownloader-2.4.7.zip"

func init() {
}

type SteamGameDl struct {
	AppId    string
	Username string
	Password string
}

func (c SteamGameDl) Run(env pufferpanel.Environment) (err error) {
	env.DisplayToConsole(true, "Downloading game from Steam")
	rootBinaryFolder := config.BinariesFolder.Value()

	err = downloadBinaries(rootBinaryFolder)
	if err != nil {
		return err
	}

	var args = []string{filepath.Join(rootBinaryFolder, "depotdownloader", "DepotDownloader.dll"), "-app", c.AppId, "-dir", env.GetRootDirectory()}
	if c.Username != "" {
		args = append(args, "-username", c.Username, "-remember-password")
		if c.Password != "" {
			args = append(args, "-password", c.Password)
		}
	}

	var success bool
	steps := pufferpanel.ExecutionData{
		//Command:          fmt.Sprintf("%s%c%s", ".", filepath.Separator, "dotnet"),
		Command:   filepath.Join(rootBinaryFolder, "dotnet-runtime", "dotnet"),
		Arguments: args,
		Callback: func(exitCode bool) {
			success = exitCode
		},
	}
	err = env.Execute(steps)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("depotdownloader exited with non-zero code")
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
		return errors.New(fmt.Sprintf("dotnet-install exited with non-zero code: %d", cmd.ProcessState.ExitCode()))
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
