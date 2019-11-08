module github.com/pufferpanel/pufferpanel/v2

go 1.12

require (
	github.com/AlecAivazis/survey/v2 v2.0.4
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/denisenkom/go-mssqldb v0.0.0-20190906004059-62cf760a6c9e // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-openapi/spec v0.19.4 // indirect
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/gobuffalo/envy v1.7.0 // indirect
	github.com/gorilla/websocket v1.4.1
	github.com/jinzhu/gorm v1.9.11
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mailgun/mailgun-go v2.0.0+incompatible
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pufferpanel/apufferi/v4 v4.0.3
	github.com/pufferpanel/pufferd/v2 v2.0.0-00010101000000-000000000000
	github.com/rogpeppe/go-internal v1.3.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.4.0
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.3
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094 // indirect
	golang.org/x/sys v0.0.0-20191023151326-f89234f9a2c2 // indirect
	golang.org/x/tools v0.0.0-20191023143423-ff611c50cd12 // indirect
	google.golang.org/appengine v1.6.2 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/pufferpanel/pufferd/v2 => ../pufferd
