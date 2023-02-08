package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v3/models"
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
					Branch: "v2",
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

				err = db.Table("servers").Where("node_id = ?", local.ID).Update("node_id", 0).Error
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
	})

	return m.Migrate()
}
