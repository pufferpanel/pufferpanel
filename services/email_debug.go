package services

import (
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

func SendEmailViaDebug(to, subject, body string, async bool) error {
	logging.Debug.Println("DEBUG EMAIL TO " + to + "\n" + subject + "\n" + body)
	return nil
}
