package services

import (
	"errors"
	"strings"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"golang.org/x/crypto/ssh"
)

type DatabaseSFTPAuthorization struct {
}

func (s *DatabaseSFTPAuthorization) ValidatePassword(username, password string) (perms *ssh.Permissions, err error) {
	parts := strings.Split(username, "#")
	if len(parts) != 2 {
		return nil, errors.New("incorrect username or password")
	}

	email := parts[0]
	serverId := parts[1]

	db, err := database.GetConnection()
	if err != nil {
		return nil, pufferpanel.ErrDatabaseNotAvailable
	}

	us := &User{DB: db}
	user, err := us.GetByEmail(email)
	if user == nil || err != nil || !us.IsValidCredentials(user, password) {
		return nil, errors.New("incorrect username or password")
	}

	ss := &Permission{DB: db}
	allowed, err := ss.HasPermission(user.ID, serverId, scopes.ScopeServerSftp)
	if err != nil || !allowed {
		return nil, errors.New("incorrect username or password")
	}

	perms = &ssh.Permissions{}
	perms.Extensions = make(map[string]string)
	perms.Extensions["server_id"] = serverId
	return perms, nil
}

func (s *DatabaseSFTPAuthorization) ValidatePublicKey(username string, key ssh.PublicKey) (*ssh.Permissions, error) {
	parts := strings.Split(username, "#")
	if len(parts) != 2 {
		return nil, errors.New("incorrect username or password")
	}

	email := parts[0]
	serverId := parts[1]

	db, err := database.GetConnection()
	if err != nil {
		return nil, pufferpanel.ErrDatabaseNotAvailable
	}

	us := &User{DB: db}
	user, err := us.GetByEmail(email)
	if user == nil || err != nil {
		return nil, errors.New("incorrect username or password")
	}

	keys, err := us.GetPublicKeys(email)
	if keys == nil || err != nil {
		return nil, errors.New("incorrect username or password")
	}

	if !keys.ContainsKey(key) {
		return nil, errors.New("incorrect username or password")
	}

	ss := &Permission{DB: db}
	allowed, err := ss.HasPermission(user.ID, serverId, scopes.ScopeServerSftp)
	if err != nil || !allowed {
		return nil, errors.New("incorrect username or password")
	}

	perms := &ssh.Permissions{}
	perms.Extensions = make(map[string]string)
	perms.Extensions["server_id"] = serverId
	return perms, nil
}
