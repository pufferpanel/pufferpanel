package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
	"reflect"
)

func Migrate(dbConn *gorm.DB) error {
	dbObjects := []interface{}{
		&models.Node{},
		&models.Server{},
		&models.User{},
		&models.Template{},
		&models.Permissions{},
		&models.Client{},
		&models.UserSetting{},
		&models.Session{},
		&models.TemplateRepo{},
	}

	for _, v := range dbObjects {
		logging.Info.Printf("Migrating model: " + reflect.TypeOf(v).Elem().Name())
		if err := dbConn.AutoMigrate(v); err != nil {
			return err
		}
	}

	m := gormigrate.New(dbConn, &gormigrate.Options{TableName: "migrations", IDColumnName: "id", IDColumnSize: 255, UseTransaction: true, ValidateUnknownMigrations: false}, []*gormigrate.Migration{
		{
			ID: "1658926619",
			Migrate: func(db *gorm.DB) error {
				logging.Info.Printf("Migrate id:1658926619")

				return db.Create(&models.TemplateRepo{
					Name:   "community",
					Url:    "https://github.com/pufferpanel/templates",
					Branch: "v3",
				}).Error
			},
		},
		{
			ID: "1677250619",
			Migrate: func(db *gorm.DB) error {
				logging.Info.Printf("Migrate id:1677250619")

				var templates []*models.Template
				err := db.Find(&templates).Error
				if err != nil {
					return err
				}

				for _, v := range templates {
					var rawMap pufferpanel.MetadataType
					err = pufferpanel.UnmarshalTo(v.Environment, &rawMap)
					if err != nil {
						return err
					}
					if rawMap.Type == "tty" || rawMap.Type == "standard" {
						rawMap.Type = "host"
						v.Environment = rawMap
						err = db.Save(&v).Error
						if err != nil {
							return err
						}
					}
				}

				return nil
			},
		},
		{
			ID: "permissions-from-v2",
			Migrate: func(db *gorm.DB) error {
				//this is going to be a nightmare
				logging.Info.Printf("Migrate id:permissions-from-v2")

				//go ahead and migrate the table, so that the columns we need are there
				err := db.AutoMigrate(&models.Permissions{})
				if err != nil {
					return err
				}

				if !db.Migrator().HasColumn(&models.Permissions{}, "admin") {
					logging.Info.Printf("No admin column exists, assuming no migration needed")
					return nil
				}

				type permissions struct {
					ID uint `gorm:"primaryKey,autoIncrement" json:"-"`

					//owners of this permission set
					UserId *uint `json:"-"`

					ClientId *uint `json:"-"`

					//if this set is for a server, what server
					ServerIdentifier *string `json:"-"`

					//and here are all the perms we support
					Admin           bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					ViewServer      bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					CreateServer    bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					ViewNodes       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					EditNodes       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					DeployNodes     bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					ViewTemplates   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					EditTemplates   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					EditUsers       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					ViewUsers       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					EditServerAdmin bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					DeleteServer    bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					PanelSettings   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`

					//these only will exist if tied to a server, and for a user
					EditServerData    bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					EditServerUsers   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					InstallServer     bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					UpdateServer      bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""` //this is unused currently
					ViewServerConsole bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					SendServerConsole bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					StopServer        bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					StartServer       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					ViewServerStats   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					ViewServerFiles   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					SFTPServer        bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
					PutServerFiles    bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
				}

				var allPerms []*permissions
				err = db.Find(&allPerms).Error
				if err != nil {
					return err
				}

				for _, v := range allPerms {
					newPerms := &models.Permissions{
						ID:               v.ID,
						UserId:           v.UserId,
						ClientId:         v.ClientId,
						ServerIdentifier: v.ServerIdentifier,
						Scopes: []*pufferpanel.Scope{
							pufferpanel.ScopeLogin,
							pufferpanel.ScopeSelfEdit,
							pufferpanel.ScopeSelfClients,
						},
					}

					//now... map all the perms to the new scopes
					if v.Admin {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeAdmin)
					}

					if v.CreateServer {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerCreate)
					}

					if v.ViewNodes {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeNodesView)
					}
					if v.EditNodes {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeNodesCreate)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeNodesDelete)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeNodesEdit)
					}
					if v.DeployNodes {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeNodesDeploy)
					}

					if v.ViewTemplates {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeTemplatesView)
					}
					if v.EditTemplates {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeTemplatesLocalEdit)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeTemplatesRepoCreate)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeTemplatesRepoDelete)
					}

					if v.EditUsers {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeUserInfoEdit)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeUserPermsEdit)
					}
					if v.ViewUsers {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeUserInfoSearch)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeUserInfoView)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeUserPermsView)
					}

					if v.PanelSettings {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeSettingsEdit)
					}

					if v.ServerIdentifier != nil && *v.ServerIdentifier != "" {
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerClientView)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerClientEdit)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerClientCreate)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerClientDelete)
						newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerStatus)

						if v.EditServerData {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerEditData)
						}
						if v.EditServerUsers {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerUserCreate)
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerUserEdit)
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerUserDelete)
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerUserView)
						}

						if v.InstallServer {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerInstall)
						}
						if v.ViewServerConsole {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerConsole)
						}
						if v.SendServerConsole {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerSendCommand)
						}

						if v.StartServer {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerStart)
						}
						if v.StopServer {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerStop)
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerKill)
						}

						if v.ViewServerStats {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerStats)
						}

						if v.SFTPServer {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerSftp)
						}
						if v.ViewServerFiles {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerFileView)
						}
						if v.PutServerFiles {
							newPerms.Scopes = pufferpanel.AddScope(newPerms.Scopes, pufferpanel.ScopeServerFileEdit)
						}
					}

					err = db.Table("permissions").Save(newPerms).Error
					if err != nil {
						return err
					}
				}

				//now... nuke the old columns
				p := &permissions{}
				for _, v := range []string{"Admin", "ViewServer", "CreateServer", "ViewNodes", "EditNodes",
					"DeployNodes", "ViewTemplates", "EditTemplates", "EditUsers", "ViewUsers", "EditServerAdmin",
					"DeleteServer", "PanelSettings", "EditServerData", "EditServerUsers", "InstallServer",
					"UpdateServer", "ViewServerConsole", "SendServerConsole", "StopServer", "StartServer",
					"ViewServerStats", "ViewServerFiles", "SFTPServer", "PutServerFiles",
				} {
					err = db.Migrator().DropColumn(p, v)
					if err != nil {
						return err
					}
				}

				return nil
			},
		},
	})

	return m.Migrate()
}
