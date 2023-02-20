package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v3/config"
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
			ID: "1676911364",
			Migrate: func(db *gorm.DB) error {
				if !config.DaemonEnabled.Value() {
					return nil
				}

				var nodes []*models.Node
				err := db.Find(&nodes).Error
				if err != nil {
					return err
				}

				var local *models.Node
				for _, n := range nodes {
					if n.ID == 1 {
						local = n
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
	})

	return m.Migrate()
}
