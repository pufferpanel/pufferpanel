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

package resolveforgeversion

import (
	"errors"
	"github.com/hashicorp/go-version"
	"github.com/pufferpanel/pufferpanel/v3"
	"path/filepath"
	"strings"
)

type ResolveForgeVersion struct {
	Version          string
	MinecraftVersion string
	OutputVariable   string
}

func (op ResolveForgeVersion) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment
	fs := args.Server.GetFileServer()

	//if a specific version wasn't specified, we have to dig around through the files....
	if op.Version == "" {
		dir := filepath.Join("libraries", "net", "minecraftforge", "forge")
		folders, err := fs.ReadDir(dir)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
		if len(folders) == 0 {
			return pufferpanel.OperationResult{Error: errors.New("forge not installed")}
		}

		var ver *version.Version
		for _, v := range folders {
			//look for folders
			if v.IsDir() {
				folderName := v.Name()
				//look for the unix file to accurately confirm this to be supported
				desiredFile := filepath.Join(dir, folderName, "unix_args.txt")
				if _, err = fs.Stat(desiredFile); err != nil {
					continue
				}
				if op.Version == "" {
					op.Version = v.Name()
					ver, _ = version.NewVersion(op.Version)
				} else if !strings.HasPrefix(folderName, op.MinecraftVersion) {
					//we need a different version of MC
					continue
				} else if ver != nil {
					//time to check to see if this a newer version
					if ver2, _ := version.NewVersion(op.Version); ver2 != nil && ver.LessThan(ver2) {
						op.Version = v.Name()
						ver = ver2
					}
				}
			}
		}
	}

	env.DisplayToConsole(true, "Resolved Forge Version: %s", op.Version)

	return pufferpanel.OperationResult{VariableOverrides: map[string]interface{}{
		op.OutputVariable: op.Version,
	}}
}
