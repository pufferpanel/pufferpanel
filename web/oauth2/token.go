package oauth2

import (
	"github.com/gin-gonic/gin"
	builder "github.com/pufferpanel/apufferi/http"
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

func registerTokens(g *gin.RouterGroup) {
	configureServer()

	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", shared.CreateOptions("POST"))

	g.POST("/validate", handleValidate)
	g.OPTIONS("/validate", shared.CreateOptions("POST"))
}

func configureServer() {
	manager := manage.NewDefaultManager()
	manager.MapClientStorage(&oauth2.ClientStore{})
	manager.MapTokenStorage(&oauth2.TokenStore{})

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

	if token := c.PostForm("token"); token != "" {
		msg.token = token
	} else {
		err := c.BindJSON(msg)
		if err != nil {
			builder.Respond(c).Status(http.StatusBadRequest).Fail().Message(err.Error()).Send()
			return
		}
	}

	builder.Respond(c).Data(true).Send()
}
