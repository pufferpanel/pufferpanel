package services

import (
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/models/view"
	o2 "github.com/pufferpanel/pufferpanel/oauth2"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"log"
	"net/http"
)

type OAuthService interface {
	HandleHTTPTokenRequest(writer http.ResponseWriter, request *http.Request)

	GetInfo(token string) (info *view.OAuthTokenInfoViewModel, valid bool, err error)

	Create(user *models.User, server *models.Server, clientId string) (clientSecret string, existing bool, err error)

	UpdateScopes(client *models.ClientInfo, server *models.Server, scopes []string) (err error)

	Delete(clientId string) (err error)

	GetByClientId(clientId string) (client *models.ClientInfo, exists bool, err error)

	GetByUser(user *models.User) (client *models.ClientInfo, exists bool, err error)
}

type oauthService struct {
	server *server.Server
}

var _oauthService *oauthService

func GetOAuthService() (service OAuthService, err error) {
	if _oauthService == nil {
		err = configureServer()
	}
	return _oauthService, err
}

func configureServer() error {
	manager := manage.NewDefaultManager()
	manager.MapClientStorage(&o2.ClientStore{})
	manager.MapTokenStorage(&o2.TokenStore{})

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	_oauthService = &oauthService{server: srv}
	return nil
}

func (oauth2 *oauthService) HandleHTTPTokenRequest(writer http.ResponseWriter, request *http.Request) {
	err := oauth2.server.HandleTokenRequest(writer, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (oauth2 *oauthService) GetInfo(token string) (info *view.OAuthTokenInfoViewModel, valid bool, err error) {
	ts := &o2.TokenStore{}
	info = &view.OAuthTokenInfoViewModel{Active: false}

	item, err := ts.GetByAccess(token)

	if err != nil {
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		return
	}

	client := &models.ClientInfo{
		ClientID: item.GetClientID(),
	}
	err = db.Set("gorm:auto_preload", true).Where(client).First(client).Error
	if err != nil {
		return
	}

	//see if the access token expiration is after now
	info = view.FromTokenInfo(item, client)
	valid = info.Active

	return
}

func (oauth2 *oauthService) Create(user *models.User, server *models.Server, clientId string) (clientSecret string, existing bool, err error) {
	return
}

func (oauth2 *oauthService) UpdateScopes(client *models.ClientInfo, server *models.Server, scopes []string) (err error) {
	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	deleteIds := make([]int, 0)
	for k, v := range client.ServerScopes {
		if v.Server.ID == server.ID || (v.Server.ID == 0 && server == nil){
			deleteIds = append(deleteIds, k)
		}
	}

	for index := len(deleteIds) - 1; index >= 0; index-- {
		client.ServerScopes = append(client.ServerScopes[:index], client.ServerScopes[index+1:]...)
	}

	//re-add new values
	for _, v := range scopes {
		replacement := &models.ClientServerScopes{
			Scope: v,
			ClientInfoID: client.ID,
		}
		if server != nil {
			replacement.ServerId = server.ID
		}

		db.Create(replacement)
		client.ServerScopes = append(client.ServerScopes, *replacement)
	}

	db.Save(client)

	return
}

func (oauth2 *oauthService) Delete(clientId string) (err error) {
	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	model := &models.ClientInfo{ClientID: clientId}

	return db.Delete(model).Error
}

func (oauth2 *oauthService) GetByClientId(clientId string) (client *models.ClientInfo, exists bool, err error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, false, err
	}

	model := &models.ClientInfo{ClientID: clientId}

	res := db.Where(model).First(model)

	return model, model.ID != 0, res.Error
}

func (oauth2 *oauthService) GetByUser(user *models.User) (client *models.ClientInfo, exists bool, err error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, false, err
	}

	model := &models.ClientInfo{UserID: user.ID}

	res := db.Where(model).First(model)

	return model, model.ID != 0, res.Error
}
