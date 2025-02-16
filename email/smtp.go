package email

import (
	"github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"github.com/wneessen/go-mail"
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

	client, err := mail.NewClient(config.EmailHost.Value(),
		mail.WithSSLPort(true),
		mail.WithUsername(config.EmailUsername.Value()),
		mail.WithPassword(config.EmailPassword.Value()),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
	)
	if err != nil {
		return err
	}
	defer utils.Close(client)

	refId, _ := uuid.NewV4()
	refIdStr := strings.ReplaceAll(refId.String(), "-", "")
	msg := mail.NewMsg()
	if err = msg.From(from); err != nil {
		return err
	}
	if err = msg.To(to); err != nil {
		return err
	}

	msg.SetMessageIDWithValue(refIdStr)
	msg.Subject(subject)
	msg.SetBodyString("text/html", body)

	return client.DialAndSend(msg)
}
