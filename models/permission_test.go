/*
 Copyright 2020 Padduck, LLC
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

package models

import "testing"

func TestPermissions_ShouldDelete(t *testing.T) {
	type fields struct {
		Admin             bool
		ViewServer        bool
		CreateServer      bool
		ViewNodes         bool
		EditNodes         bool
		DeployNodes       bool
		ViewTemplates     bool
		EditUsers         bool
		ViewUsers         bool
		EditServerAdmin   bool
		DeleteServer      bool
		EditServerData    bool
		EditServerUsers   bool
		InstallServer     bool
		UpdateServer      bool
		ViewServerConsole bool
		SendServerConsole bool
		StopServer        bool
		StartServer       bool
		ViewServerStats   bool
		ViewServerFiles   bool
		SFTPServer        bool
		PutServerFiles    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "One set",
			fields: fields{Admin: true},
			want:   false,
		},
		{
			name:   "All false",
			fields: fields{},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Permissions{
				Admin:             tt.fields.Admin,
				ViewServer:        tt.fields.ViewServer,
				CreateServer:      tt.fields.CreateServer,
				ViewNodes:         tt.fields.ViewNodes,
				EditNodes:         tt.fields.EditNodes,
				DeployNodes:       tt.fields.DeployNodes,
				ViewTemplates:     tt.fields.ViewTemplates,
				EditUsers:         tt.fields.EditUsers,
				ViewUsers:         tt.fields.ViewUsers,
				EditServerAdmin:   tt.fields.EditServerAdmin,
				DeleteServer:      tt.fields.DeleteServer,
				EditServerData:    tt.fields.EditServerData,
				EditServerUsers:   tt.fields.EditServerUsers,
				InstallServer:     tt.fields.InstallServer,
				UpdateServer:      tt.fields.UpdateServer,
				ViewServerConsole: tt.fields.ViewServerConsole,
				SendServerConsole: tt.fields.SendServerConsole,
				StopServer:        tt.fields.StopServer,
				StartServer:       tt.fields.StartServer,
				ViewServerStats:   tt.fields.ViewServerStats,
				ViewServerFiles:   tt.fields.ViewServerFiles,
				SFTPServer:        tt.fields.SFTPServer,
				PutServerFiles:    tt.fields.PutServerFiles,
			}
			if got := p.ShouldDelete(); got != tt.want {
				t.Errorf("ShouldDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}
