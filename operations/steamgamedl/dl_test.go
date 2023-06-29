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
	"github.com/pufferpanel/pufferpanel/v2/config"
	"testing"
)

func Test_downloadSteamcmd(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		version string
		args    []string
	}{
		{
			name:    "download steamcmd",
			wantErr: false,
			args:    []string{},
		},
	}

	_ = config.BinariesFolder.Set("C:\\Temp\\pufferpanel\\binaries", false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := downloadBinaries(config.BinariesFolder.Value()); (err != nil) != tt.wantErr {
				t.Errorf("downloadBinaries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
