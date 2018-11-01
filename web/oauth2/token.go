package oauth2

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/oauth2"
	"github.com/pufferpanel/pufferpanel/shared"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"log"
	"net/http"
)

func registerTokens(g *gin.RouterGroup) {
	handle()
	g.POST("/token", handle())
	g.OPTIONS("/token", shared.CreateOptions("POST"))

	g.POST("/info", shared.NotImplemented)
	g.OPTIONS("/info", shared.CreateOptions("POST"))
}

func handle() func(*gin.Context){
	manager := manage.NewDefaultManager()
	manager.MapClientStorage(&oauth2.ClientStore{})
	manager.MapTokenStorage(&oauth2.TokenStore{})
	manager.MapAccessGenerate(oauth2.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))

	db, err := database.GetConnection()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&oauth2.ClientInfo{}, &oauth2.TokenInfo{})

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetClientScopeHandler(handleScopes)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	return func(c *gin.Context) {
		handleTokenRequest(srv, c)
	}
}

func handleTokenRequest(srv *server.Server, c *gin.Context) {
	err := srv.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func handleScopes(clientID, scope string) (allowed bool, err error) {
	return true, nil
}