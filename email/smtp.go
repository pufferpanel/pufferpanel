package email

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"strings"
)

type smtpProvider struct {
	Provider
}

func init() {
	providers["smtp"] = smtpProvider{}
}

func (smtpProvider) Send(to, subject, body string) error {
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
	return smtp.SendMail(host, auth, from, []string{to}, data)
}
