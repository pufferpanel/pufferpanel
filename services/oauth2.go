package services

import (
	"github.com/pufferpanel/pufferpanel/database"
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

	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	db.AutoMigrate(&o2.ClientInfo{}, &o2.TokenInfo{}, &o2.ClientServerScopes{})

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

	item, err := ts.GetByAccess(token)

	if err != nil {
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		return
	}

	client := o2.ClientInfo{
		ClientID: item.GetClientID(),
	}
	err = db.Set("gorm:auto_preload", true).Where(&client).First(&client).Error
	if err != nil {
		return
	}

	//see if the access token expiration is after now
	info = view.FromTokenInfo(item, client)
	valid = info.Active

	return
}
