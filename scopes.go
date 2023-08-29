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

import "encoding/json"

type Scope struct {
	Value     string
	ForServer bool
}

var (
	ScopeAdmin            = registerNonServerScope("admin")
	ScopeLogin            = registerNonServerScope("login")        //can you log in
	ScopeOAuth2Auth       = registerNonServerScope("oauth2.auth")  //can you validate user credentials over OAuth2
	ScopeNodesView        = registerNonServerScope("nodes.view")   //can you globally view nodes
	ScopeNodesCreate      = registerNonServerScope("nodes.create") //can you create nodes
	ScopeNodesEdit        = registerNonServerScope("nodes.edit")   //can you edit an existing node
	ScopeNodesDelete      = registerNonServerScope("nodes.delete") //can you delete a node
	ScopeNodesDeploy      = registerNonServerScope("nodes.deploy") //can you deploy the node (this has secret info, which is why it's special)
	ScopeSelfEdit         = registerNonServerScope("self.edit")    //can you manage your own account
	ScopeSelfClients      = registerNonServerScope("self.clients") //can the user create and manage OAuth2 clients for their own account
	ScopeSelfSettingsView = registerNonServerScope("self.settings.view")
	ScopeSelfSettingsEdit = registerNonServerScope("self.settings.edit")

	ScopeServerList         = registerNonServerScope("servers.list")
	ScopeServerAdmin        = registerServerScope("server.admin")
	ScopeServerCreate       = registerServerScope("server.create")
	ScopeServerDelete       = registerServerScope("server.delete")
	ScopeServerEditAdmin    = registerServerScope("server.edit.admin")
	ScopeServerEditData     = registerServerScope("server.edit.data")
	ScopeServerEditFlags    = registerServerScope("server.edit.flags")
	ScopeServerEditName     = registerServerScope("server.edit.name")
	ScopeServerViewAdmin    = registerServerScope("server.view.admin")
	ScopeServerViewData     = registerServerScope("server.view.data")
	ScopeServerClientView   = registerServerScope("server.client.view")
	ScopeServerClientEdit   = registerServerScope("server.client.edit")
	ScopeServerClientAdd    = registerServerScope("server.client.add")
	ScopeServerClientDelete = registerServerScope("server.client.delete")
	ScopeServerUserView     = registerServerScope("server.users.view")
	ScopeServerUserCreate   = registerServerScope("server.users.add")
	ScopeServerUserEdit     = registerServerScope("server.users.edit")
	ScopeServerUserDelete   = registerServerScope("server.users.delete")
	ScopeServerTaskView     = registerServerScope("server.tasks.view")
	ScopeServerTaskRun      = registerServerScope("server.tasks.run")
	ScopeServerTaskCreate   = registerServerScope("server.tasks.create")
	ScopeServerTaskDelete   = registerServerScope("server.tasks.delete")
	ScopeServerReload       = registerServerScope("server.reload")
	ScopeServerStart        = registerServerScope("server.start")
	ScopeServerStop         = registerServerScope("server.stop")
	ScopeServerKill         = registerServerScope("server.kill")
	ScopeServerInstall      = registerServerScope("server.install")
	ScopeServerFileGet      = registerServerScope("server.files.get")
	ScopeServerFileEdit     = registerServerScope("server.files.edit")
	ScopeServerSftp         = registerServerScope("server.sftp")
	ScopeServerLogs         = registerServerScope("server.logs")
	ScopeServerSendCommand  = registerServerScope("server.console.send")
	ScopeServerStat         = registerServerScope("server.stat")
	ScopeServerStatus       = registerServerScope("server.status")

	ScopeSettingsEdit        = registerNonServerScope("settings.edit")
	ScopeTemplatesView       = registerNonServerScope("templates.view")
	ScopeTemplatesEdit       = registerNonServerScope("templates.local.edit")
	ScopeTemplatesRepoView   = registerNonServerScope("templates.repo.view")
	ScopeTemplatesRepoAdd    = registerNonServerScope("templates.repo.add")
	ScopeTemplatesRepoDelete = registerNonServerScope("templates.repo.remove")

	ScopeUserInfoSearch = registerNonServerScope("users.info.search")
	ScopeUserInfoView   = registerNonServerScope("users.info.view")
	ScopeUserInfoEdit   = registerNonServerScope("users.info.edit")
	ScopeUserPermsView  = registerNonServerScope("users.perms.view")
	ScopeUserPermsEdit  = registerNonServerScope("users.perms.edit")

	ScopePanel = registerNonServerScope("panel")
)

func (s Scope) String() string {
	return s.Value
}

func (s Scope) Matches(string string) bool {
	return string == s.String()
}

func (s Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}

var allScopes []Scope

func registerScope(s Scope) Scope {
	allScopes = append(allScopes, s)
	return s
}
func registerNonServerScope(s string) Scope {
	return registerScope(Scope{Value: s})
}
func registerServerScope(s string) Scope {
	return registerScope(Scope{Value: s, ForServer: true})
}

func GetScope(str string) Scope {
	for _, v := range allScopes {
		if v.Matches(str) {
			return v
		}
	}
	return Scope{Value: str}
}

func ContainsScope(arr []Scope, value Scope) bool {
	return containsScope(arr, value, ScopeAdmin)
}

func ContainsServerScope(arr []Scope, value Scope) bool {
	return containsScope(arr, value, ScopeServerAdmin, ScopeAdmin)
}

func AddScope(source []Scope, addition Scope) []Scope {
	for _, v := range source {
		if v.Matches(addition.String()) {
			return source
		}
	}
	return append(source, addition)
}

func RemoveScope(source []Scope, removal Scope) []Scope {
	replacement := make([]Scope, 0)
	for _, v := range source {
		if !v.Matches(removal.String()) {
			replacement = append(replacement, v)
		}
	}
	return replacement
}

func containsScope(arr []Scope, desired ...Scope) bool {
	for _, v := range arr {
		for _, z := range desired {
			if v.Matches(z.String()) {
				return true
			}
		}
	}

	return false
}
