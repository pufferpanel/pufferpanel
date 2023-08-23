package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
)

func migrate(dbConn *gorm.DB) error {

	m := gormigrate.New(dbConn, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1658926619",
			Migrate: func(db *gorm.DB) error {
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
	})

	return m.Migrate()
}
