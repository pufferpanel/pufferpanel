package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"gorm.io/gorm"
	"os"
	"path/filepath"
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
		&models.Backup{},
	}

	session := dbConn.Session(&gorm.Session{})
	migrator := session.Migrator()

	z := gormigrate.New(session, &gormigrate.Options{TableName: "migrations", IDColumnName: "id", IDColumnSize: 255, UseTransaction: true, ValidateUnknownMigrations: false}, []*gormigrate.Migration{
		{
			ID: "1726675832-mysql",
			Migrate: func(db *gorm.DB) error {
				logging.Info.Printf("Migrate id:1726675832-mysql")

				if config.DatabaseDialect.Value() != "mysql" {
					return nil
				}

				//at this point for mysql, just manually do the queries...

				type FKs struct {
					Table string `gorm:"column:TABLE_NAME"`
					Name  string `gorm:"column:CONSTRAINT_NAME"`
				}
				var results []FKs
				err := db.Raw("SELECT TABLE_NAME, CONSTRAINT_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = (SELECT DATABASE())").Scan(&results).Error
				if err != nil {
					return err
				}

				for _, v := range results {
					if err = db.Exec("alter table " + v.Table + " drop foreign key " + v.Name).Error; err != nil {
						return err
					}
				}

				return nil
			},
		},
		{
			ID: "1726675832",
			Migrate: func(db *gorm.DB) error {
				//delete all indices, because we need to just recreate them
				logging.Info.Printf("Migrate id:1726675832")

				for _, v := range dbObjects {
					m := db.Migrator()
					indices, err := m.GetIndexes(v)
					if err != nil {
						return err
					}

					for _, z := range indices {
						if isPrim, ok := z.PrimaryKey(); ok && isPrim {
							continue
						}

						if err = m.DropIndex(v, z.Name()); err != nil {
							return err
						}
					}
				}

				return nil
			},
		},
	})

	if err := z.Migrate(); err != nil {
		return err
	}

	if err := migrator.AutoMigrate(dbObjects...); err != nil {
		return err
	}

	m := gormigrate.New(session, &gormigrate.Options{TableName: "migrations", IDColumnName: "id", IDColumnSize: 255, UseTransaction: true, ValidateUnknownMigrations: false}, []*gormigrate.Migration{
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
					err = utils.UnmarshalTo(v.Environment, &rawMap)
					if err != nil {
						logging.Error.Printf("Failed to migrate template %s, template saved off. %s", v.Name, err)
						err = saveToFile("template-"+v.Name+".json", []byte(v.RawValue))
						if err != nil {
							return err
						}
						//return err
						err = db.Delete(&v).Error
						if err != nil {
							return err
						}
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
						Scopes: []*scopes.Scope{
							scopes.ScopeLogin,
							scopes.ScopeSelfEdit,
							scopes.ScopeSelfClients,
						},
					}

					//now... map all the perms to the new scopes
					if v.Admin {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeAdmin)
					}

					if v.CreateServer {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerCreate)
					}

					if v.ViewNodes {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeNodesView)
					}
					if v.EditNodes {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeNodesCreate)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeNodesDelete)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeNodesEdit)
					}
					if v.DeployNodes {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeNodesDeploy)
					}

					if v.ViewTemplates {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeTemplatesView)
					}
					if v.EditTemplates {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeTemplatesLocalEdit)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeTemplatesRepoCreate)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeTemplatesRepoDelete)
					}

					if v.EditUsers {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeUserInfoEdit)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeUserPermsEdit)
					}
					if v.ViewUsers {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeUserInfoSearch)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeUserInfoView)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeUserPermsView)
					}

					if v.PanelSettings {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeSettingsEdit)
					}

					if v.ServerIdentifier != nil && *v.ServerIdentifier != "" {
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerClientView)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerClientEdit)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerClientCreate)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerClientDelete)
						newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerStatus)

						if v.EditServerData {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerEditData)
						}
						if v.EditServerUsers {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerUserCreate)
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerUserEdit)
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerUserDelete)
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerUserView)
						}

						if v.InstallServer {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerInstall)
						}
						if v.ViewServerConsole {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerConsole)
						}
						if v.SendServerConsole {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerSendCommand)
						}

						if v.StartServer {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerStart)
						}
						if v.StopServer {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerStop)
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerKill)
						}

						if v.ViewServerStats {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerStats)
						}

						if v.SFTPServer {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerSftp)
						}
						if v.ViewServerFiles {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerFileView)
						}
						if v.PutServerFiles {
							newPerms.Scopes = scopes.AddScope(newPerms.Scopes, scopes.ScopeServerFileEdit)
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

func saveToFile(filename string, data []byte) error {
	//just dump it into working dir
	err := os.MkdirAll("migrations", 0755)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join("migrations", filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer utils.Close(file)
	_, err = file.Write(data)
	return err
}
