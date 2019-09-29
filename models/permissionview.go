package models

type PermissionView struct {
	Username         string `json:"username,omitempty"`
	ServerIdentifier string `json:"serverIdentifier,omitempty"`

	EditServerData    bool `json:"editServerData"`
	EditServerUsers   bool `json:"editServerUsers"`
	InstallServer     bool `json:"installServer"`
	UpdateServer      bool `json:"-"` //this is unused currently
	ViewServerConsole bool `json:"viewServerConsole"`
	SendServerConsole bool `json:"sendServerConsole"`
	StopServer        bool `json:"stopServer"`
	StartServer       bool `json:"startServer"`
	ViewServerStats   bool `json:"viewServerStats"`
	ViewServerFiles   bool `json:"viewServerFiles"`
	SFTPServer        bool `json:"sftpServer"`
	PutServerFiles    bool `json:"putServerFiles"`
}

func FromPermission(p *Permissions) *PermissionView {
	var serverIdentifier string

	if p.ServerIdentifier != nil {
		serverIdentifier = *p.ServerIdentifier
	}

	return &PermissionView{
		Username:          p.User.Username,
		ServerIdentifier:  serverIdentifier,
		EditServerData:    p.EditServerData,
		EditServerUsers:   p.EditServerUsers,
		InstallServer:     p.InstallServer,
		UpdateServer:      p.UpdateServer,
		ViewServerConsole: p.ViewServerConsole,
		SendServerConsole: p.SendServerConsole,
		StopServer:        p.StopServer,
		StartServer:       p.StartServer,
		ViewServerStats:   p.ViewServerStats,
		ViewServerFiles:   p.ViewServerFiles,
		SFTPServer:        p.SFTPServer,
		PutServerFiles:    p.PutServerFiles,
	}
}

func (p *PermissionView) CopyTo(model *Permissions) {
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
}
