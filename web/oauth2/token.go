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

var oauth2Server *server.Server
var jwtService *oauth2.JWTAccessGenerate

func registerTokens(g *gin.RouterGroup) {
	configureServer()

	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", shared.CreateOptions("POST"))

	g.POST("/validate", shared.NotImplemented)
	g.OPTIONS("/validate", shared.CreateOptions("POST"))
}

func configureServer() {

	jwtService = oauth2.NewJWTAccessGenerate([]byte(oauth2.GetJWTSecret()), jwt.SigningMethodHS512)

	manager := manage.NewDefaultManager()
	manager.MapClientStorage(&oauth2.ClientStore{})
	manager.MapTokenStorage(&oauth2.TokenStore{})
	manager.MapAccessGenerate(jwtService)

	db, err := database.GetConnection()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&oauth2.ClientInfo{}, &oauth2.TokenInfo{})

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	oauth2Server = srv
}

func handleTokenRequest(c *gin.Context) {
	err := oauth2Server.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func handleValidate(c *gin.Context) {
	type body struct {
		token string
	}
	msg := &body{}
	c.BindJSON(msg)
	c.JSON(200, jwtService.Validate(msg.token))
}
