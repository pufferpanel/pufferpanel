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
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/environments"
	"github.com/pufferpanel/pufferpanel/v2/environments/envs"
	"path"
	"strings"
)

const InstallerUrl = "https://files.minecraftforge.net/maven/net/minecraftforge/forge/${version}/forge-${version}-installer.jar"

type ForgeDl struct {
	Version  string
	Filename string
}

func (op ForgeDl) Run(env envs.Environment) error {
	jarDownload := strings.Replace(InstallerUrl, "${version}", op.Version, -1)

	localFile, err := environments.DownloadViaMaven(jarDownload, env)
	if err != nil {
		return err
	}

	//copy from the cache
	return pufferpanel.CopyFile(localFile, path.Join(env.GetRootDirectory(), op.Filename))
}
