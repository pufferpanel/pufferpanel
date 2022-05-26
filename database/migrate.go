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
	})

	return m.Migrate()
}
