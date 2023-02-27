package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v2/models"
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
			ID: "1665609381",
			Migrate: func(db *gorm.DB) error {
				var nodes []*models.Node
				err := db.Find(&nodes).Error
				if err != nil {
					return err
				}

				var local *models.Node
				for _, n := range nodes {
					if (n.PrivateHost == "localhost" || n.PrivateHost == "127.0.0.1") && (n.PublicHost == "localhost" || n.PublicHost == "127.0.0.1") {
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
		{
			ID: "client-add-nullable",
			Migrate: func(db *gorm.DB) error {
				err := db.Migrator().AlterColumn(&models.Client{}, "ServerId")
				if err != nil {
					return err
				}
				err = db.Table("clients").Where("server_id = ?", "").Update("server_id", nil).Error
				return err
			},
		},
		{
			ID: "servers-add-nullable",
			Migrate: func(db *gorm.DB) error {
				err := db.Migrator().AlterColumn(&models.Server{}, "RawNodeID")
				if err != nil {
					return err
				}

				err = db.Table("servers").Where("node_id = ?", 0).Update("node_id", nil).Error
				return err
			},
		},
	})

	return m.Migrate()
}
