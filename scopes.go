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

package pufferpanel

type Scope string

const (
	ScopeAdmin            = Scope("admin")
	ScopeLogin            = Scope("login")        //can you log in
	ScopeOAuth2Auth       = Scope("oauth2.auth")  //can you validate user credentials over OAuth2
	ScopeNodesView        = Scope("nodes.view")   //can you globally view nodes
	ScopeNodesCreate      = Scope("nodes.create") //can you create nodes
	ScopeNodesEdit        = Scope("nodes.edit")   //can you edit an existing node
	ScopeNodesDelete      = Scope("nodes.delete") //can you delete a node
	ScopeNodesDeploy      = Scope("nodes.deploy") //can you deploy the node (this has secret info, which is why it's special)
	ScopeSelfEdit         = Scope("self.edit")    //can you manage your own account
	ScopeSelfClients      = Scope("self.clients") //can the user create and manage OAuth2 clients for their own account
	ScopeSelfSettingsView = Scope("self.settings.view")
	ScopeSelfSettingsEdit = Scope("self.settings.edit")

	ScopeServerList         = Scope("servers.list")
	ScopeServerAdmin        = Scope("server.admin")
	ScopeServerCreate       = Scope("server.create")
	ScopeServerDelete       = Scope("server.delete")
	ScopeServerEditAdmin    = Scope("server.edit.admin")
	ScopeServerEditData     = Scope("server.edit.data")
	ScopeServerEditFlags    = Scope("server.edit.flags")
	ScopeServerEditName     = Scope("server.edit.name")
	ScopeServerViewAdmin    = Scope("server.view.admin")
	ScopeServerViewData     = Scope("server.view.data")
	ScopeServerClientView   = Scope("server.client.view")
	ScopeServerClientEdit   = Scope("server.client.edit")
	ScopeServerClientAdd    = Scope("server.client.add")
	ScopeServerClientDelete = Scope("server.client.delete")
	ScopeServerUserView     = Scope("server.users.view")
	ScopeServerUserCreate   = Scope("server.users.add")
	ScopeServerUserEdit     = Scope("server.users.edit")
	ScopeServerUserDelete   = Scope("server.users.delete")
	ScopeServerTaskView     = Scope("server.tasks.view")
	ScopeServerTaskRun      = Scope("server.tasks.run")
	ScopeServerTaskCreate   = Scope("server.tasks.create")
	ScopeServerTaskDelete   = Scope("server.tasks.delete")
	ScopeServerReload       = Scope("server.reload")
	ScopeServerStart        = Scope("server.start")
	ScopeServerStop         = Scope("server.stop")
	ScopeServerKill         = Scope("server.kill")
	ScopeServerInstall      = Scope("server.install")
	ScopeServerFileGet      = Scope("server.files.get")
	ScopeServerFileEdit     = Scope("server.files.edit")
	ScopeServerSftp         = Scope("server.sftp")
	ScopeServerLogs         = Scope("server.logs")
	ScopeServerSendCommand  = Scope("server.console.send")
	ScopeServerStat         = Scope("server.stat")
	ScopeServerStatus       = Scope("server.status")

	ScopeSettingsEdit        = Scope("settings.edit")
	ScopeTemplatesView       = Scope("templates.view")
	ScopeTemplatesEdit       = Scope("templates.local.edit")
	ScopeTemplatesRepoView   = Scope("templates.repo.view")
	ScopeTemplatesRepoAdd    = Scope("templates.repo.add")
	ScopeTemplatesRepoDelete = Scope("templates.repo.remove")

	ScopeUserInfoSearch = Scope("users.info.search")
	ScopeUserInfoView   = Scope("users.info.view")
	ScopeUserInfoEdit   = Scope("users.info.edit")
	ScopeUserPermsView  = Scope("users.perms.view")
	ScopeUserPermsEdit  = Scope("users.perms.edit")

	ScopePanel = Scope("panel")
)

func AllKnownScopes() []Scope {
	return []Scope{
		ScopeLogin,
		ScopeOAuth2Auth,
		ScopeNodesView,
		ScopeNodesCreate,
		ScopeNodesEdit,
		ScopeNodesDelete,
		ScopeNodesDeploy,
		ScopeSelfEdit,
		ScopeSelfClients,
		ScopeSelfSettingsView,
		ScopeSelfSettingsEdit,
		ScopeServerList,
		ScopeServerAdmin,
		ScopeServerCreate,
		ScopeServerDelete,
		ScopeServerEditAdmin,
		ScopeServerEditData,
		ScopeServerEditFlags,
		ScopeServerEditName,
		ScopeServerViewAdmin,
		ScopeServerViewData,
		ScopeServerClientView,
		ScopeServerClientEdit,
		ScopeServerClientAdd,
		ScopeServerClientDelete,
		ScopeServerUserView,
		ScopeServerUserCreate,
		ScopeServerUserEdit,
		ScopeServerUserDelete,
		ScopeServerTaskView,
		ScopeServerTaskRun,
		ScopeServerTaskCreate,
		ScopeServerTaskDelete,
		ScopeServerReload,
		ScopeServerStart,
		ScopeServerStop,
		ScopeServerKill,
		ScopeServerInstall,
		ScopeServerFileGet,
		ScopeServerFileEdit,
		ScopeServerSftp,
		ScopeServerLogs,
		ScopeServerSendCommand,
		ScopeServerStat,
		ScopeServerStatus,
		ScopeSettingsEdit,
		ScopeTemplatesView,
		ScopeTemplatesEdit,
		ScopeTemplatesRepoView,
		ScopeTemplatesRepoAdd,
		ScopeTemplatesRepoDelete,
		ScopeUserInfoSearch,
		ScopeUserInfoView,
		ScopeUserInfoEdit,
		ScopeUserPermsView,
		ScopeUserPermsEdit,
	}
}

func (s Scope) String() string {
	return string(s)
}

func (s Scope) Matches(string string) bool {
	return string == s.String()
}

func ContainsScope(arr []Scope, value Scope) bool {
	return containsScope(arr, value, ScopeAdmin)
}

func ContainsServerScope(arr []Scope, value Scope) bool {
	return containsScope(arr, value, ScopeServerAdmin, ScopeAdmin)
}

func AddScope(source []Scope, addition Scope) []Scope {
	for _, v := range source {
		if v == addition {
			return source
		}
	}
	return append(source, addition)
}

func RemoveScope(source []Scope, removal Scope) []Scope {
	replacement := make([]Scope, 0)
	for _, v := range source {
		if v != removal {
			replacement = append(replacement, v)
		}
	}
	return replacement
}

func containsScope(arr []Scope, desired ...Scope) bool {
	for _, v := range arr {
		for _, z := range desired {
			if v == z {
				return true
			}
		}
	}

	return false
}
