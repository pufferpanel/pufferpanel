package claims

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"time"
)

type UserClaims struct {
	jwt.StandardClaims
	UserId uint `json:"user_id"`
}

type ServerClaims struct {
	UserClaims
	ServerId uint     `json:"server_id"`
	Scopes   []string `json:"scopes"`
}

func NewUserClaim(user *models.User) *UserClaims {
	return &UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.ID,
	}
}

func NewServerClaim(user *models.User, server *models.Server) *ServerClaims {
	return &ServerClaims{
		UserClaims: UserClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			UserId: user.ID,
		},
		ServerId: 0,
		Scopes:   nil,
	}
}
