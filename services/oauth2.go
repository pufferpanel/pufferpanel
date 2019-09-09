package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/apufferi/scope"
	"github.com/pufferpanel/pufferpanel"
	"github.com/pufferpanel/pufferpanel/models"
	o2 "github.com/pufferpanel/pufferpanel/oauth2"
	"github.com/satori/go.uuid"
	oauth "gopkg.in/oauth2.v3"
	oauthErrors "gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"net/http"
	"strings"
	"time"
)

var oauthServer *server.Server

type OAuth struct {
	server *server.Server
	DB     *gorm.DB
}

func GetOAuth(db *gorm.DB) *OAuth {
	if oauthServer == nil {
		configureServer()
	}
	return &OAuth{server: oauthServer, DB: db}
}

func configureServer() {
	manager := manage.NewDefaultManager()
	manager.MapClientStorage(&o2.ClientStore{})
	manager.MapTokenStorage(&o2.TokenStore{})

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *oauthErrors.Response) {
		logging.Build(logging.ERROR).WithMessage("internal error on oauth2 service").WithError(err).Log()
		return
	})

	srv.SetResponseErrorHandler(func(re *oauthErrors.Response) {
		logging.Build(logging.ERROR).WithMessage("response error on oauth2 service").WithError(re.Error).Log()
	})

	oauthServer = srv
}

func (oauth2 *OAuth) HandleHTTPTokenRequest(writer http.ResponseWriter, request *http.Request) {
	err := oauth2.server.HandleTokenRequest(writer, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (oauth2 *OAuth) ValidationBearerToken(r *http.Request) (ti oauth.TokenInfo, err error) {
	return oauth2.server.ValidationBearerToken(r)
}

func (oauth2 *OAuth) GetInfo(token string) (info *models.TokenInfoView, valid bool, err error) {
	tokenInfo, client, err := oauth2.GetByToken(token)

	if tokenInfo == nil || client == nil {
		return nil, false, nil
	}

	//see if the access token expiration is after now
	info = models.FromTokenInfo(tokenInfo, client)
	valid = oauth2.HasTokenExpired(tokenInfo)
	return
}

func (oauth2 *OAuth) Create(user *models.User, server *models.Server, clientId string, panel bool, scopes ...string) (clientSecret string, err error) {
	clientSecret = strings.Replace(uuid.NewV4().String(), "-", "", -1)

	if clientId == "" {
		panel = true
		clientId = CreateInternalClientId(user, server)
	}

	ci := &models.ClientInfo{
		ClientID: clientId,
		UserID:   user.ID,
		Secret:   clientSecret,
		Panel:    panel,
	}

	res := oauth2.DB.Create(ci)

	if res.Error != nil {
		return "", res.Error
	}

	for _, s := range scopes {
		scopeInfo := &models.ClientServerScopes{
			ClientInfoID: ci.ID,
			Scope:        s,
		}
		if server != nil {
			scopeInfo.ServerId = &server.ID
		}
		err = oauth2.DB.Create(scopeInfo).Error
		if err != nil {
			return "", err
		}
	}

	return
}

func (oauth2 *OAuth) RemoveScope(client *models.ClientInfo, server *models.Server, scope string) error {
	if server != nil && server.ID == 0 {
		res := oauth2.DB.Where(server).First(server)
		if res.Error != nil {
			return res.Error
		}
	}

	if client.ID == 0 {
		var query *gorm.DB
		if server != nil && server.ID != 0 {
			query = oauth2.DB.Preload("ServerScopes", "server_id = ?", server.ID)
		} else {
			query = oauth2.DB.Preload("ServerScopes", "server_id IS NULL")
		}
		res := query.Where(client).First(client)
		if res.Error != nil {
			return res.Error
		}
		if client.ID == 0 {
			return pufferpanel.ErrClientNotFound
		}
	}

	for _, v := range client.ServerScopes {
		if v.Scope == scope {
			res := oauth2.DB.Delete(v)
			return res.Error
		}
	}

	return nil
}

func (oauth2 *OAuth) AddScope(client *models.ClientInfo, server *models.Server, scopes ...string) error {
	trans := oauth2.DB.Begin()
	commit := false
	defer func() {
		if commit {
			trans.Commit()
		} else {
			trans.Rollback()
		}
	}()

	if server != nil && server.ID == 0 {
		res := trans.Where(server).First(server)
		if res.Error != nil {
			return res.Error
		}
	}

	if client.ID == 0 {
		var query *gorm.DB
		if server != nil && server.ID != 0 {
			query = trans.Preload("ServerScopes", "server_id = ?", server.ID)
		} else {
			query = trans.Preload("ServerScopes", "server_id IS NULL")
		}
		res := query.Where(client).First(client)
		if res.Error != nil {
			return res.Error
		}
		if client.ID == 0 {
			return pufferpanel.ErrClientNotFound
		}
	}

	for _, scope := range scopes {
		toAdd := true
		for _, s := range client.ServerScopes {
			if s.Scope == scope {
				toAdd = false
				break
			}
		}
		if toAdd {
			replacement := &models.ClientServerScopes{
				Scope:        scope,
				ClientInfoID: client.ID,
			}
			if server != nil && server.ID != 0 {
				replacement.ServerId = &server.ID
			}

			res := trans.Create(replacement)
			return res.Error
		}
	}

	commit = true
	return nil
}

func (oauth2 *OAuth) UpdateScopes(client *models.ClientInfo, server *models.Server, scopes ...string) (err error) {
	if server != nil && server.ID == 0 {
		res := oauth2.DB.Where(server).First(server)
		if res.Error != nil {
			return res.Error
		}
	}

	if client.ID == 0 {
		var query *gorm.DB
		if server != nil && server.ID != 0 {
			query = oauth2.DB.Preload("ServerScopes", "server_id = ?", server.ID)
		} else {
			query = oauth2.DB.Preload("ServerScopes", "server_id IS NULL")
		}
		res := query.Where(client).First(client)
		if res.Error != nil {
			return res.Error
		}
		if client.ID == 0 {
			return pufferpanel.ErrClientNotFound
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
			oauth2.DB.Delete(v)
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
				Scope:        v,
				ClientInfoID: client.ID,
			}
			if server != nil && server.ID != 0 {
				replacement.ServerId = &server.ID
			}

			oauth2.DB.Create(replacement)
		}
	}

	return
}

func (oauth2 *OAuth) Delete(clientId string) (err error) {
	model := &models.ClientInfo{ClientID: clientId}

	return oauth2.DB.Where(model).Delete(model).Error
}

func (oauth2 *OAuth) GetByClientId(clientId string) (client *models.ClientInfo, exists bool, err error) {
	model := &models.ClientInfo{ClientID: clientId}

	res := oauth2.DB.Set("gorm:auto_preload", true).Where(model).First(model)

	return model, model.ID != 0, res.Error
}

func (oauth2 *OAuth) GetByUser(user *models.User) (client *models.ClientInfo, exists bool, err error) {
	model := &models.ClientInfo{UserID: user.ID, Panel: true}

	res := oauth2.DB.Set("gorm:auto_preload", true).Where(model).First(model)

	return model, model.ID != 0, res.Error
}

func (oauth2 *OAuth) GetByScope(scope string, userId, serverId *uint, internal bool) (*models.ClientInfos, error) {
	model := &models.ClientInfos{}

	query := oauth2.DB.Set("gorm:auto_preload", true)
	query = query.Joins("JOIN client_server_scopes ON client_server_scopes.client_info_id = client_infos.id")

	scopes := &models.ClientServerScopes{Scope: scope}
	if serverId != nil {
		query = query.Joins("JOIN servers ON servers.id = client_server_scopes.server_id")
		query = query.Where(&models.Server{ID: *serverId})
	}
	if userId != nil {
		query = query.Joins("JOIN users ON users.id = client_infos.user_id")
		query = query.Where(&models.User{ID: *userId})
	}

	err := query.Where(scopes).Find(model).Error
	if err != nil {
		return nil, err
	}

	return model, err
}

func (oauth2 *OAuth) HasRights(accessToken string, serverId *uint, scope string) (client *models.ClientInfo, allowed bool, err error) {
	ts := o2.TokenStore{}

	var ti oauth.TokenInfo
	if ti, err = ts.GetByAccess(accessToken); err != nil {
		return
	}
	if oauth2.HasTokenExpired(ti) {
		return
	}

	converted, ok := ti.(*models.TokenInfo)
	if !ok {
		err = pufferpanel.ErrInvalidTokenState
		return
	}

	for _, v := range converted.ClientInfo.ServerScopes {
		if (v.ServerId == nil && serverId != nil) || (v.ServerId != nil && serverId == nil) {
			continue
		}
		if ((v.ServerId == nil && serverId == nil) || (*v.ServerId == *serverId)) && v.Scope == scope {
			return &converted.ClientInfo, true, nil
		}
	}

	return &converted.ClientInfo, false, nil
}

func (oauth2 *OAuth) HasTokenExpired(info oauth.TokenInfo) (expired bool) {
	if info == nil {
		return true
	}

	if info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Before(time.Now()) {
		return true
	}

	return false
}

func (oauth2 *OAuth) GetByToken(token string) (tokenInfo oauth.TokenInfo, client *models.ClientInfo, err error) {
	ts := &o2.TokenStore{}
	tokenInfo, err = ts.GetByAccess(token)

	if err != nil {
		return
	}

	if oauth2.HasTokenExpired(tokenInfo) {
		return nil, nil, nil
	}

	client = &models.ClientInfo{
		ClientID: tokenInfo.GetClientID(),
	}
	err = oauth2.DB.Set("gorm:auto_preload", true).Where(client).First(client).Error
	if err != nil {
		return
	}

	return tokenInfo, client, err
}

func (oauth2 *OAuth) UpdateExpirationTime(tokenInfo oauth.TokenInfo, duration time.Duration) (err error) {
	res := oauth2.DB.Set("gorm:association_save_reference", false).Model(tokenInfo).Update("access_create_at", time.Now())
	return res.Error
}

func (oauth2 *OAuth) CreateSession(user *models.User) (string, error) {
	ci, _, err := oauth2.GetByUser(user)

	if err != nil {
		return "", err
	}
	if ci == nil || ci.ID == 0 {
		return "", pufferpanel.ErrUserNotFound
	}

	valid := false
	for _, v := range ci.ServerScopes {
		if v.ServerId == nil && v.Scope == scope.Login {
			valid = true
			break
		}
	}

	if !valid {
		return "", pufferpanel.ErrLoginNotPermitted
	}

	ti := &models.TokenInfo{
		ClientInfoID:    ci.ID,
		AccessCreateAt:  time.Now(),
		AccessExpiresIn: 1 * time.Hour,
		Access:          strings.Replace(uuid.NewV4().String(), "-", "", -1),
	}

	err = oauth2.DB.Create(ti).Error

	if err != nil {
		ti.Access = ""
	}

	return ti.Access, err
}

func (oauth2 *OAuth) GetForServer(serverId uint, includeAdmin bool) (*models.ClientInfos, error) {
	model := &models.ClientInfos{}

	exclude := &models.ClientInfos{}
	//exclusion list
	if !includeAdmin {
		query := oauth2.DB.Set("gorm:auto_preload", true)
		query = query.Table("client_infos")
		query = query.Joins("JOIN client_server_scopes ON client_server_scopes.client_info_id = client_infos.id")
		err := query.Select("client_id").Where("scope = 'servers.admin'").Find(&exclude).Error
		if err != nil {
			return nil, err
		}
	}

	query := oauth2.DB.Set("gorm:auto_preload", true)
	query = query.Joins("JOIN client_server_scopes ON client_server_scopes.client_info_id = client_infos.id")
	query = query.Joins("JOIN servers ON servers.id = client_server_scopes.server_id")
	query = query.Where(&models.Server{ID: serverId})
	query = query.Where("client_infos.panel = 1")
	if len(*exclude) != 0 {
		ids := make([]string, len(*exclude))
		for k, v := range *exclude {
			ids[k] = v.ClientID
		}
		query = query.Not("client_id", ids)
	}

	err := query.Find(model).Error
	if err != nil {
		return nil, err
	}

	return model, err
}

func CreateInternalClientId(user *models.User, server *models.Server) string {
	return fmt.Sprintf(".internal_%d_%d", user.ID, server.ID)
}
