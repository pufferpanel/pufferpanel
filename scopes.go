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

package pufferpanel

type Scope string

const (
	//none scope to allow defining that something doesn't need any specific permission
	ScopeNone = Scope("none")

	//generic
	ScopeOAuth2Auth = Scope("oauth2.auth")

	//server
	ScopeServersAdmin       = Scope("servers.admin")
	ScopeServersView        = Scope("servers.view")
	ScopeServersEdit        = Scope("servers.edit")
	ScopeServersEditAdmin   = Scope("servers.edit.admin")
	ScopeServersEditUsers   = Scope("servers.edit.users")
	ScopeServersCreate      = Scope("servers.create")
	ScopeServersDelete      = Scope("servers.delete")
	ScopeServersInstall     = Scope("servers.install")
	ScopeServersUpdate      = Scope("servers.update")
	ScopeServersConsole     = Scope("servers.console")
	ScopeServersConsoleSend = Scope("servers.console.send")
	ScopeServersStop        = Scope("servers.stop")
	ScopeServersStart       = Scope("servers.start")
	ScopeServersStat        = Scope("servers.stats")
	ScopeServersSFTP        = Scope("servers.sftp")
	ScopeServersFilesGet    = Scope("servers.files.get")
	ScopeServersFilesPut    = Scope("servers.files.put")

	//node
	ScopeNodesView   = Scope("nodes.view")
	ScopeNodesEdit   = Scope("nodes.edit")
	ScopeNodesDeploy = Scope("nodes.deploy")

	//template
	ScopeTemplatesView = Scope("templates.view")
	ScopeTemplatesEdit = Scope("templates.edit")

	//user
	ScopeUsersView = Scope("users.view")
	ScopeUsersEdit = Scope("users.edit")
)
