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

package scope

type Scope string

const (
	//generic
	Login = Scope("login")

	OAuth2Auth = Scope("oauth2.auth")

	//server
	ServersAdmin       = Scope("servers.admin")
	ServersView        = Scope("servers.view")
	ServersEdit        = Scope("servers.edit")
	ServersEditAdmin   = Scope("servers.edit.admin")
	ServersEditUsers   = Scope("servers.edit.users")
	ServersCreate      = Scope("servers.create")
	ServersDelete      = Scope("servers.delete")
	ServersInstall     = Scope("servers.install")
	ServersUpdate      = Scope("servers.update")
	ServersConsole     = Scope("servers.console")
	ServersConsoleSend = Scope("servers.console.send")
	ServersStop        = Scope("servers.stop")
	ServersStart       = Scope("servers.start")
	ServersStat        = Scope("servers.stats")
	ServersSFTP        = Scope("servers.sftp")
	ServersFilesGet    = Scope("servers.files.get")
	ServersFilesPut    = Scope("servers.files.put")

	//node
	NodesView   = Scope("nodes.view")
	NodesEdit   = Scope("nodes.edit")
	NodesDeploy = Scope("nodes.deploy")

	//template
	TemplatesView = Scope("templates.view")

	//user
	UsersView = Scope("users.view")
	UsersEdit = Scope("users.edit")
)
