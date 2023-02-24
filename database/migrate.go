package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func migrate(dbConn *gorm.DB) error {

	m := gormigrate.New(dbConn, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1626910428",
			Migrate: func(db *gorm.DB) error {
				_ = db.Migrator().DropIndex(&models.Server{}, "uix_servers_name")
				return nil
			},
			Rollback: nil,
		},
		{
			ID: "1658926619",
			Migrate: func(db *gorm.DB) error {
				return db.Create(&models.TemplateRepo{
					Name:   "community",
					Url:    "https://github.com/pufferpanel/templates",
					Branch: "v2.6",
				}).Error
			},
		},
		{
			ID: "1665609381",
			Migrate: func(db *gorm.DB) error {
				var nodes []*models.Node
				err := db.Find(&nodes).Error
				if err != nil {
					return err
				}

				var local *models.Node
				for _, v := range nodes {
					if v.Name == "LocalNode" {
						local = v
					}
				}

				if local == nil {
					return nil
				}

				err = db.Table("servers").Where("node_id = ?", local.ID).Update("node_id", nil).Error
				if err != nil {
					return err
				}
				err = db.Delete(local).Error
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			ID: "1677250619",
			Migrate: func(db *gorm.DB) error {
				var templates []*models.Template
				err := db.Find(&templates).Error
				if err != nil {
					return err
				}

				for _, v := range templates {
					rawMap := make(map[string]interface{})
					err = pufferpanel.UnmarshalTo(v.Environment, &rawMap)
					if err != nil {
						return err
					}
					declaredEnv := cast.ToString(rawMap["type"])
					if declaredEnv == "tty" || declaredEnv == "standard" {
						rawMap["type"] = "host"
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
	})

	return m.Migrate()
}
