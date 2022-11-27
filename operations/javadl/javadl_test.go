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

package javadl

import (
	"github.com/pufferpanel/pufferpanel/v2/config"
	test "github.com/pufferpanel/pufferpanel/v2/testing"
	"os"
	"os/exec"
	"testing"
)

func Test_downloadJava(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		version string
	}{
		{
			version: "17",
			name:    "download java 17",
			wantErr: false,
		},
	}

	_ = os.Setenv("PATH", os.Getenv("PATH")+":"+config.BinariesFolder.Value())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := JavaDl{Version: tt.version}
			if err := op.Run(&test.Environment{}); (err != nil) != tt.wantErr {
				t.Errorf("downloadJava() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, err := exec.LookPath("java" + op.Version)
			if err != nil {
				t.Errorf("downloadJava() failed to add to path %v", err)
			}
		})
	}
}
