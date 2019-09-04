package pufferpanel

//generic
const ScopeLogin = "login"

//oauth2
const ScopeOauth2Info = "oauth2.info"

//server
const ScopeServerAdmin = "servers.admin"
const ScopeViewServers = "servers.view"
const ScopeEditServers = "servers.edit"
const ScopeEditServerAsUser = "servers.edit.user"
const ScopeEditServerUsers = "servers.edit.users"
const ScopeCreateServers = "servers.create"

const ScopeServerConsole = "servers.console"
const ScopeStopServers = "servers.stop"
const ScopeStartServers = "servers.start"
const ScopeKillServers = "servers.kill"
const ScopeStatServers = "servers.stats"
const ScopeFilesServers = "servers.files"
const ScopeGetFilesServers = "servers.files.get"
const ScopePutFilesServers = "servers.files.put"

//node
const ScopeViewNodes = "nodes.view"
const ScopeEditNode = "nodes.edit"
const ScopeDeployNode = "nodes.deploy"

//template
const ScopeViewTemplates = "templates.view"

//user
const ScopeViewUsers = "users.view"
const ScopeEditUsers = "users.edit"

func GetDefaultUserServerScopes() []string {
	return []string{
		ScopeViewServers,
		ScopeServerConsole,
		ScopeStopServers,
		ScopeStartServers,
		ScopeKillServers,
		ScopeStatServers,
		ScopeFilesServers,
		ScopeGetFilesServers,
		ScopePutFilesServers,
	}
}
