/*
 Copyright 2023 PufferPanel

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

package curseforge

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	if config.CurseForgeKey.Value() == "" {
		return nil, errors.New("CurseForge key is required to use this module")
	}

	projectId, err := cast.ToUintE(op.OperationArgs["projectId"])
	if err != nil {
		return nil, err
	}
	fileId, err := cast.ToUintE(op.OperationArgs["fileId"])
	if op.OperationArgs["fileId"] != "" && err != nil {
		return nil, err
	} else if op.OperationArgs["fileId"] == "" {
		fileId = 0
	}

	javaBinary := cast.ToString(op.OperationArgs["java"])
	if javaBinary == "" {
		javaBinary = "java"
	}

	return &CurseForge{ProjectId: projectId, FileId: fileId, JavaBinary: javaBinary}, nil
}

func (of OperationFactory) Key() string {
	return "curseforge"
}

var Factory OperationFactory
