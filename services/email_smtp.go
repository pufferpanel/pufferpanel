package services

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"io"
	"strings"
)

func SendEmailViaSMTP(to, subject, body string, async bool) error {
	from := config.EmailFrom.Value()
	if from == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailFrom.Key())
	}

	host := config.EmailHost.Value()
	if host == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailHost.Key())
	}

	var auth sasl.Client

	if username := config.EmailUsername.Value(); username != "" {
		auth = sasl.NewPlainClient("", username, config.EmailPassword.Value())
	} else {
		auth = sasl.NewAnonymousClient("")
	}

	data := strings.NewReader("Subject: " + subject + "\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" + body)

	if async {
		go func(host string, auth sasl.Client, from, to string, data io.Reader) {
			err := smtp.SendMail(host, auth, from, []string{to}, data)
			if err != nil {
				logging.Error.Printf("Error sending email: %s", err)
			}
		}(host, auth, from, to, data)
		return nil
	} else {
		return smtp.SendMail(host, auth, from, []string{to}, data)
	}
}
