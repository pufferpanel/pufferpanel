package tests

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
)

var loginNoLoginUser = &models.User{
	Username:  "loginNoLoginUser",
	Email:     "noscope@example.com",
	OtpActive: false,
}

const loginNoLoginUserPassword = "dontletmein"

var loginNoServerViewUser = &models.User{
	Username:  "loginNoServerViewUser",
	Email:     "test@example.com",
	OtpActive: false,
}

const loginNoServerViewUserPassword = "testing123"

var loginAdminUser = &models.User{
	Username:  "loginAdminUser",
	Email:     "admin@example.com",
	OtpActive: false,
}

const loginAdminUserPassword = "asdfasdf"

var loginNoAdminWithServersUser = &models.User{
	Username:  "loginNoAdminWithServersUser",
	Email:     "notadmin@example.com",
	OtpActive: false,
}

const loginNoAdminWithServersUserPassword = "dowiuzlaslf"

func init() {
	_ = loginNoLoginUser.SetPassword(loginNoLoginUserPassword)
	_ = loginNoServerViewUser.SetPassword(loginNoServerViewUserPassword)
	_ = loginAdminUser.SetPassword(loginAdminUserPassword)
	_ = loginNoAdminWithServersUser.SetPassword(loginNoAdminWithServersUserPassword)
}

func prepareUsers(db *gorm.DB) error {
	err := initNoLoginUser(db)
	if err != nil {
		return err
	}

	err = initLoginNoServersUser(db)
	if err != nil {
		return err
	}

	err = initLoginAdminUser(db)
	if err != nil {
		return err
	}

	err = initLoginNoAdminWithServersUser(db)
	if err != nil {
		return err
	}

	return nil
}

func initNoLoginUser(db *gorm.DB) error {
	return db.Create(loginNoLoginUser).Error
}

func initLoginNoServersUser(db *gorm.DB) error {
	err := db.Create(loginNoServerViewUser).Error
	if err != nil {
		return err
	}

	perms := &models.Permissions{
		UserId: &loginNoServerViewUser.ID,
		Scopes: []*pufferpanel.Scope{pufferpanel.ScopeLogin},
	}
	err = db.Create(perms).Error
	return err
}

func initLoginAdminUser(db *gorm.DB) error {
	err := db.Create(loginAdminUser).Error
	if err != nil {
		return err
	}

	perms := &models.Permissions{
		UserId: &loginAdminUser.ID,
		Scopes: []*pufferpanel.Scope{pufferpanel.ScopeAdmin},
	}
	err = db.Create(perms).Error
	return err
}

func initLoginNoAdminWithServersUser(db *gorm.DB) error {
	return db.Create(loginNoAdminWithServersUser).Error
}

func createSession(db *gorm.DB, user *models.User) (string, error) {
	ss := &services.Session{DB: db}
	return ss.CreateForUser(user)
}

func createSessionAdmin() (string, error) {
	db, err := database.GetConnection()
	if err != nil {
		return "", err
	}
	return createSession(db, loginAdminUser)
}

const CreateServerData = `{
  "type": "testing",
  "display": "API Test Server",
  "name": "Test",
  "data": {
    "ip": {
      "value": "127.0.0.1",
      "required": true,
      "desc": "What IP to bind the server to",
      "display": "IP"
    },
    "port": {
      "value": "123",
      "required": true,
      "desc": "What port to bind the server to",
      "display": "Port",
      "type": "integer"
    }
  },
  "install": [
    {
      "type": "writefile",
      "text": "install should not be doing this",
      "target": "bad.txt"
    }
  ],
  "run": {
    "command": [
      "echo nothing'"      
    ],
    "stop": "norun"
  },
  "environment": {
    "type": "tty"
  }
}`

const TestServerData = `{
  "type": "testing",
  "display": "API Test Server",
  "name": "Test Var",
  "data": {
    "ip": {
      "value": "0.0.0.0",
      "required": true,
      "desc": "What IP to bind the server to",
      "display": "IP"
    },
    "port": {
      "value": "25565",
      "required": true,
      "desc": "What port to bind the server to",
      "display": "Port",
      "type": "integer"
    }
  },
  "install": [
    {
      "type": "writefile",
      "text": "install successed",
      "target": "installed.txt"
    }
  ],
  "run": {
    "command": [
      "echo started"      
    ],
    "stop": "stop"
  },
  "environment": {
    "type": "standard"
  }
}`

const TemplateData = `{
  "type": "minecraft-java",
  "display": "Vanilla - Minecraft",
  "data": {
    "version": {
      "value": "latest",
      "required": true,
      "desc": "Version of Minecraft to install",
      "display": "Version",
      "internal": false
    },
    "memory": {
      "value": "1024",
      "required": true,
      "desc": "How much memory in MB to allocate to the Java Heap",
      "display": "Memory (MB)",
      "internal": false,
      "type": "integer"
    },
    "ip": {
      "value": "0.0.0.0",
      "required": true,
      "desc": "What IP to bind the server to",
      "display": "IP",
      "internal": false
    },
    "port": {
      "value": "25565",
      "required": true,
      "desc": "What port to bind the server to",
      "display": "Port",
      "internal": false,
      "type": "integer"
    },
    "eula": {
      "value": "false",
      "required": true,
      "desc": "Do you (or the server owner) agree to the <a href='https://account.mojang.com/documents/minecraft_eula'>Minecraft EULA?</a>",
      "display": "EULA Agreement",
      "internal": false,
      "type": "boolean"
    },
    "motd": {
      "value": "A Minecraft Server\\n\\u00A79 hosted on PufferPanel",
      "required": true,
      "desc": "This is the message that is displayed in the server list of the client, below the name. The MOTD does support <a href='https://minecraft.wiki/w/Formatting_codes' target='_blank'>color and formatting codes</a>.",
      "display": "MOTD message of the day",
      "internal": false
    },
    "javaversion": {
      "type": "integer",
      "desc": "Version of Java to use",
      "display": "Java Version",
      "value": "17",
      "required": true
    }
  },
  "install": [
    {
      "type": "javadl",
      "version": "${javaversion}"
    },
    {
      "type": "mojangdl",
      "version": "${version}",
      "target": "server.jar"
    },
    {
      "type": "writefile",
      "text": "server-ip=${ip}\nserver-port=${port}\nmotd=${motd}\n",
      "target": "server.properties"
    },
    {
      "type": "writefile",
      "text": "eula=${eula}",
      "target": "eula.txt"
    }
  ],
  "run": {
    "command": "java${javaversion} -Xmx${memory}M -Dlog4j2.formatMsgNoLookups=true -jar server.jar nogui",
    "stop": "stop"
  },
  "environment": {
    "type": "standard"
  }
}`
