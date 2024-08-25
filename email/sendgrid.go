package email

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendgridProvider struct {
	Provider
}

func init() {
	providers["sendgrid"] = sendgridProvider{}
}

func (sendgridProvider) Send(to, subject, body string) error {
	from := config.EmailFrom.Value()
	if from == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailFrom.Key())
	}

	key := config.EmailKey.Value()
	if key == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailKey.Key())
	}

	fromEmail := mail.NewEmail("", from)
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(fromEmail, subject, toEmail, "", body)

	client := sendgrid.NewSendClient(key)
	_, err := client.Send(message)
	return err
}
