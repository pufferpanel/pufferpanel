package services

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/models/view"
	o2 "github.com/pufferpanel/pufferpanel/oauth2"
	oauth "gopkg.in/oauth2.v3"
	oauthErrors "gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"log"
	"net/http"
	"time"
)

type OAuthService interface {
	HandleHTTPTokenRequest(writer http.ResponseWriter, request *http.Request)

	GetInfo(token string) (info *view.OAuthTokenInfoViewModel, valid bool, err error)

	Create(user *models.User, server *models.Server, clientId string) (clientSecret string, existing bool, err error)

	UpdateScopes(client *models.ClientInfo, server *models.Server, scopes ...string) (err error)

	Delete(clientId string) (err error)

	GetByClientId(clientId string) (client *models.ClientInfo, exists bool, err error)

	GetByUser(user *models.User) (client *models.ClientInfo, exists bool, err error)

	HasRights(accessToken string, serverId *uint, scope string) (client *models.ClientInfo, allowed bool, err error)

	HasTokenExpired(info oauth.TokenInfo) (expired bool)

	ValidationBearerToken(r *http.Request) (ti oauth.TokenInfo, err error)
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

	srv.SetInternalErrorHandler(func(err error) (re *oauthErrors .Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *oauthErrors.Response) {
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

func (oauth2 *oauthService) ValidationBearerToken(r *http.Request) (ti oauth.TokenInfo, err error){
	return oauth2.server.ValidationBearerToken(r)
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

func (oauth2 *oauthService) UpdateScopes(client *models.ClientInfo, server *models.Server, scopes ...string) (err error) {
	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	if server != nil && server.ID == 0 {
		res := db.Where(server).First(server)
		if res.Error != nil {
			return res.Error
		}
	}

	if client.ID == 0 {
		var query *gorm.DB
		if server != nil && server.ID != 0 {
			query = db.Preload("ServerScopes", "server_id = ?", server.ID)
		} else {
			query = db.Preload("ServerScopes", "server_id IS NULL")
		}
		res := query.Where(client).First(client)
		if res.Error != nil {
			return res.Error
		}
		if client.ID == 0 {
			return errors.New("no client with given information")
		}
	}

	//delete ones which don't exist on the new list
	for _, v := range client.ServerScopes {
		toDelete := true
		for _, s := range scopes {
			if s == v.Scope {
				toDelete = false
				break
			}
		}
		if toDelete {
			db.Delete(v)
		}
	}

	//add new values
	for _, v := range scopes {
		toAdd := true
		for _, s := range client.ServerScopes {
			if v == s.Scope {
				toAdd = false
				break
			}
		}
		if toAdd {
			replacement := &models.ClientServerScopes{
				Scope: v,
				ClientInfoID: client.ID,
			}
			if server != nil && server.ID != 0{
				replacement.ServerId = &server.ID
			}

			db.Create(replacement)
		}
	}

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

func (oauth2 *oauthService) HasRights(accessToken string, serverId *uint, scope string) (client *models.ClientInfo, allowed bool, err error) {
	ts := o2.TokenStore{}

	var ti oauth.TokenInfo
	if  ti, err = ts.GetByAccess(accessToken); err != nil {
		return
	}
	if oauth2.HasTokenExpired(ti) {
		return
	}

	converted, ok := ti.(*models.TokenInfo)
	if !ok {
		err = errors.New("token info state was invalid")
		return
	}

	for _, v := range converted.ClientInfo.ServerScopes {
		if v.ServerId == nil || v.ServerId == serverId {
			if v.Scope == scope {
				return &converted.ClientInfo, true, nil
			}
		}
	}

	return &converted.ClientInfo, false, nil
}

func (oauth2 *oauthService) HasTokenExpired(info oauth.TokenInfo) (expired bool) {
	if info == nil {
		return true
	}

	if info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Before(time.Now()) {
		return true
	}

	return false
}