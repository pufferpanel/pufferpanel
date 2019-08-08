package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/models"
	uuid2 "github.com/satori/go.uuid"
	"strings"
)

type Server struct {
	DB *gorm.DB
}

type ServerSearch struct {
	Username string
	NodeId   uint
	NodeName string
	Name     string
	PageSize uint
	Page     uint
}

func (ss *Server) Search(searchCriteria ServerSearch) (*models.Servers, uint, error) {
	servers := &models.Servers{}

	query := ss.DB

	if searchCriteria.NodeId != 0 {
		query = query.Where(&models.Server{NodeID: searchCriteria.NodeId})
	} else if searchCriteria.NodeName != "" {
		query = query.Joins("JOIN nodes n ON servers.node_id = n.id AND n.name = ?", searchCriteria.NodeName)
	}

	if searchCriteria.Username != "" {
		query = query.Joins("JOIN client_server_scopes css ON css.server_id = servers.id AND css.scope = 'servers.view'")
		query = query.Joins("JOIN client_infos ci ON ci.id = css.client_info_id")
		query = query.Joins("JOIN users ON users.id = ci.user_id")
		query = query.Where("users.username = ?", searchCriteria.Username)
	}

	nameFilter := strings.Replace(searchCriteria.Name, "*", "%", -1)

	if nameFilter != "" && nameFilter != "%" {
		query = query.Where("name LIKE ?", nameFilter)
	}

	var count uint
	err := query.Model(&servers).Count(&count).Error

	if err != nil {
		return nil, 0, err
	}

	res := query.Preload("Node").Offset((searchCriteria.Page - 1) * searchCriteria.PageSize).Limit(searchCriteria.PageSize).Find(servers)

	return servers, count, res.Error
}

func (ss *Server) Get(id string) (*models.Server, bool, error) {
	if id == "" {
		return nil, false, nil
	}
	model := &models.Server{
		Identifier: id,
	}

	res := ss.DB.Where(model).Preload("Node").First(model)

	return model, model.ID != 0, res.Error
}

func (ss *Server) Update(model *models.Server) error {
	res := ss.DB.Save(model)
	return res.Error
}

func (ss *Server) Delete(id uint) error {
	model := &models.Server{
		ID: id,
	}

	res := ss.DB.Delete(model)
	return res.Error
}

func (ss *Server) Create(model *models.Server) (err error) {
	uuid := uuid2.NewV4()
	generatedId := strings.ToUpper(uuid.String())[0:8]
	model.Identifier = generatedId

	res := ss.DB.Create(model)
	if res.Error != nil {
		err = res.Error
		return
	}
	return
}
