package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
)

type ServerService struct {
	db *gorm.DB
}

func GetServerService() (*ServerService, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	service := &ServerService{
		db: db,
	}

	return service, nil
}

func (ss *ServerService) GetAll() (models.Server, error) {
	servers := models.Server{}

	res := ss.db.Find(&servers)

	return servers, res.Error
}

//TODO: Waiting on user objects with rights to implement correctly
func (ss *ServerService) GetForUser(userId int) (models.Server, error) {
	servers := models.Server{}

	res := ss.db.Find(&servers)

	return servers, res.Error
}

func (ss *ServerService) Get(id int) (models.Server, bool, error) {
	model := models.Server{}

	res := ss.db.First(&model, id)

	return model, model.Id != 0, res.Error
}

func (ss *ServerService) Update(node models.Server) error {
	res := ss.db.Update(&node)
	return res.Error
}

func (ss *ServerService) Delete(id int) error {
	model := models.Server{
		Id: id,
	}

	res := ss.db.Delete(&model)
	return res.Error
}