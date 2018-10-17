package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
)

type LocationService struct {
	db *gorm.DB
}

func GetLocationService() (*LocationService, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	service := &LocationService{}
	service.db = db

	return service, nil
}

func (ls LocationService) GetAll() models.Locations {
	locations := models.Locations{}

	ls.db.Find(&locations)

	return locations
}

func (ls LocationService) Get(id int) (models.Location, bool) {
	model := models.Location{}

	ls.db.First(&model, id)

	return model, model.Id != 0
}

func (ls LocationService) GetByCode(code string) (models.Location, bool) {
	model := models.Location{}
	model.Code = code

	ls.db.First(&model)

	return model, model.Id != 0
}

func (ls LocationService) Delete(id int) {
	model, exist := ls.Get(id)
	if !exist {
		return
	}

	ls.db.Delete(&model)
}

func (ls LocationService) Update(location models.Location) {
	model, exist := ls.Get(location.Id)
	if !exist {
		return
	}

	ls.db.Model(&model).Updates(models.Location{Name: location.Name})
}