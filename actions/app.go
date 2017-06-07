package actions

import (
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/i18n"

	"github.com/pufferpanel/pufferpanel/models"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
	"github.com/pufferpanel/pufferpanel/actions/api"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env:         ENV,
			SessionName: "_pufferpanel_session",
		})
		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}
		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(middleware.CSRF)

		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		T, err = i18n.New(packr.NewBox("../locales"), "en-US")
		if err != nil {
			log.Fatal(err)
		}
		app.Use(T.Middleware())

		//register root level
		Register(app)

		//register API
		api.Register(app)

		app.ServeFiles("/assets", packr.NewBox("../public/assets"))
	}

	return app
}
