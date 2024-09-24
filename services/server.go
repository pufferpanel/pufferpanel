package services

import (
	"github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (ss *Server) Search(searchCriteria ServerSearch) (records []*models.Server, total int64, err error) {
	query := ss.DB

	if searchCriteria.NodeId != 0 {
		query = query.Where(&models.Server{NodeID: searchCriteria.NodeId})
	} else if searchCriteria.NodeName != "" {
		if searchCriteria.NodeName == "LocalNode" {
			query = query.Where(&models.Server{NodeID: 0, RawNodeID: nil})
		} else {
			query = query.Joins("JOIN nodes n ON servers.node_id = n.id AND n.name = ?", searchCriteria.NodeName)
		}
	}

	if searchCriteria.Username != "" {
		query = query.Joins("JOIN permissions p ON servers.identifier = p.server_identifier")
		query = query.Joins("JOIN users ON p.user_id = users.id")
		query = query.Where("users.username = ?", searchCriteria.Username)
	}

	nameFilter := strings.Replace(searchCriteria.Name, "*", "%", -1)

	if nameFilter != "" && nameFilter != "%" {
		query = query.Where("name LIKE ?", nameFilter)
	}

	err = query.Model(&records).Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	err = query.Preload(clause.Associations).Offset(int((searchCriteria.Page - 1) * searchCriteria.PageSize)).Limit(int(searchCriteria.PageSize)).Order("servers.name").Find(&records).Error

	return
}

func (ss *Server) Get(id string) (*models.Server, error) {
	if id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	model := &models.Server{
		Identifier: id,
	}

	err := ss.DB.Preload(clause.Associations).First(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (ss *Server) Update(model *models.Server) error {
	res := ss.DB.Omit(clause.Associations).Save(model)
	return res.Error
}

// Delete a server by ID, This is _not_ ran in a transaction automatically to allow for more flexibility
// Callers should set the DB to be a transaction if needed
// (Because Gorm V2 has removed `RollbackUnlessCommitted1)
func (ss *Server) Delete(id string) error {
	model := &models.Server{
		Identifier: id,
	}

	err := ss.DB.Delete(models.Permissions{}, "server_identifier = ?", id).Error
	if err != nil {
		return err
	}

	err = ss.DB.Delete(models.Client{}, "server_id = ?", id).Error
	if err != nil {
		return err
	}

	err = ss.DB.Delete(model).Error
	if err != nil {
		return err
	}

	return nil
}

func (ss *Server) Create(model *models.Server) error {
	if model.Identifier == "" {
		uniqueId, err := uuid.NewV4()
		if err != nil {
			return err
		}
		generatedId := strings.ToUpper(uniqueId.String())[0:8]
		model.Identifier = generatedId
	}

	res := ss.DB.Omit(clause.Associations).Create(model)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
