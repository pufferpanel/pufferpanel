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

type PermissionView struct {
	Username         string `json:"username,omitempty"`
	Email            string `json:"email,omitempty"`
	ServerIdentifier string `json:"serverIdentifier,omitempty"`

	EditServerData    bool `json:"editServerData,omitempty"`
	EditServerUsers   bool `json:"editServerUsers,omitempty"`
	InstallServer     bool `json:"installServer,omitempty"`
	UpdateServer      bool `json:"-"` //this is unused currently
	ViewServerConsole bool `json:"viewServerConsole,omitempty"`
	SendServerConsole bool `json:"sendServerConsole,omitempty"`
	StopServer        bool `json:"stopServer,omitempty"`
	StartServer       bool `json:"startServer,omitempty"`
	ViewServerStats   bool `json:"viewServerStats,omitempty"`
	ViewServerFiles   bool `json:"viewServerFiles,omitempty"`
	SFTPServer        bool `json:"sftpServer,omitempty"`
	PutServerFiles    bool `json:"putServerFiles,omitempty"`

	Admin           bool `json:"admin,omitempty"`
	ViewServer      bool `json:"viewServers,omitempty"`
	CreateServer    bool `json:"createServers,omitempty"`
	ViewNodes       bool `json:"viewNodes,omitempty"`
	EditNodes       bool `json:"editNodes,omitempty"`
	DeployNodes     bool `json:"deployNodes,omitempty"`
	ViewTemplates   bool `json:"viewTemplates,omitempty"`
	EditTemplates   bool `json:"editTemplates,omitempty"`
	EditUsers       bool `json:"editUsers,omitempty"`
	ViewUsers       bool `json:"viewUsers,omitempty"`
	EditServerAdmin bool `json:"editServerAdmin,omitempty"`
	DeleteServer    bool `json:"deleteServers,omitempty"`
}

func FromPermission(p *Permissions) *PermissionView {
	model := &PermissionView{
		Username: p.User.Username,
		Email: p.User.Email,
	}

	//only show server specific perms
	if p.ServerIdentifier != nil {
		model.ServerIdentifier = *p.ServerIdentifier
		model.EditServerData = p.EditServerData
		model.EditServerUsers = p.EditServerUsers
		model.InstallServer = p.InstallServer
		model.UpdateServer = p.UpdateServer
		model.ViewServerConsole = p.ViewServerConsole
		model.SendServerConsole = p.SendServerConsole
		model.StopServer = p.StopServer
		model.StartServer = p.StartServer
		model.ViewServerStats = p.ViewServerStats
		model.ViewServerFiles = p.ViewServerFiles
		model.SFTPServer = p.SFTPServer
		model.PutServerFiles = p.PutServerFiles
	} else {
		model.Admin = p.Admin
		model.ViewServer = p.ViewServer
		model.CreateServer = p.CreateServer
		model.ViewNodes = p.ViewNodes
		model.EditNodes = p.EditNodes
		model.DeployNodes = p.DeployNodes
		model.ViewTemplates = p.ViewTemplates
		model.EditTemplates = p.EditTemplates
		model.EditUsers = p.EditUsers
		model.ViewUsers = p.ViewUsers
		model.EditServerAdmin = p.EditServerAdmin
		model.DeleteServer = p.DeleteServer
	}

	return model
}

//Copies perms from the view to the model
//This will only copy what it knows about the server
func (p *PermissionView) CopyTo(model *Permissions, copyAdminFlags bool) {
	if model.ServerIdentifier != nil {
		model.EditServerData = p.EditServerData
		model.EditServerUsers = p.EditServerUsers
		model.InstallServer = p.InstallServer
		model.UpdateServer = p.UpdateServer
		model.ViewServerConsole = p.ViewServerConsole
		model.SendServerConsole = p.SendServerConsole
		model.StopServer = p.StopServer
		model.StartServer = p.StartServer
		model.ViewServerStats = p.ViewServerStats
		model.ViewServerFiles = p.ViewServerFiles
		model.SFTPServer = p.SFTPServer
		model.PutServerFiles = p.PutServerFiles
	} else if copyAdminFlags {
		model.Admin = p.Admin
		model.ViewServer = p.ViewServer
		model.CreateServer = p.CreateServer
		model.ViewNodes = p.ViewNodes
		model.EditNodes = p.EditNodes
		model.DeployNodes = p.DeployNodes
		model.ViewTemplates = p.ViewTemplates
		model.EditTemplates = p.EditTemplates
		model.EditUsers = p.EditUsers
		model.ViewUsers = p.ViewUsers
		model.EditServerAdmin = p.EditServerAdmin
		model.DeleteServer = p.DeleteServer
	}
}
