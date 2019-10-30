package models

type PermissionView struct {
	//Username         string `json:"username,omitempty"`
	Email            string `json:"email,omitempty"`
	ServerIdentifier string `json:"serverIdentifier,omitempty"`

	EditServerData    bool `json:"editServerData,omitempty,string"`
	EditServerUsers   bool `json:"editServerUsers,omitempty,string"`
	InstallServer     bool `json:"installServer,omitempty,string"`
	UpdateServer      bool `json:"-"` //this is unused currently
	ViewServerConsole bool `json:"viewServerConsole,omitempty,string"`
	SendServerConsole bool `json:"sendServerConsole,omitempty,string"`
	StopServer        bool `json:"stopServer,omitempty,string"`
	StartServer       bool `json:"startServer,omitempty,string"`
	ViewServerStats   bool `json:"viewServerStats,omitempty,string"`
	ViewServerFiles   bool `json:"viewServerFiles,omitempty,string"`
	SFTPServer        bool `json:"sftpServer,omitempty,string"`
	PutServerFiles    bool `json:"putServerFiles,omitempty,string"`

	Admin           bool `json:"admin,omitempty,string"`
	ViewServer      bool `json:"viewServers,omitempty,string"`
	CreateServer    bool `json:"createServers,omitempty,string"`
	ViewNodes       bool `json:"viewNodes,omitempty,string"`
	EditNodes       bool `json:"editNodes,omitempty,string"`
	DeployNodes     bool `json:"deployNodes,omitempty,string"`
	ViewTemplates   bool `json:"viewTemplates,omitempty,string"`
	EditUsers       bool `json:"editUsers,omitempty,string"`
	ViewUsers       bool `json:"viewUsers,omitempty,string"`
	EditServerAdmin bool `json:"editServerAdmin,omitempty,string"`
	DeleteServer    bool `json:"deleteServers,omitempty,string"`
}

func FromPermission(p *Permissions) *PermissionView {
	model := &PermissionView{
		//Username: p.User.Username,
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
		model.EditUsers = p.EditUsers
		model.ViewUsers = p.ViewUsers
		model.EditServerAdmin = p.EditServerAdmin
		model.DeleteServer = p.DeleteServer
	}
}
