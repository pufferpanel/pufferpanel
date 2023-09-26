package services

import (
	"github.com/mailgun/mailgun-go"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

func SendEmailViaMailgun(to, subject, body string, async bool) error {
	domain := config.EmailDomain.Value()
	if domain == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailDomain.Key())
	}

	from := config.EmailFrom.Value()
	if from == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailFrom.Key())
	}

	key := config.EmailKey.Value()
	if key == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailKey.Key())
	}

	mg := mailgun.NewMailgun(domain, key)
	message := mg.NewMessage(from, subject, "", to)
	message.SetHtml(body)

	if async {
		go func(mgI *mailgun.MailgunImpl, messageI *mailgun.Message) {
			_, _, err := mgI.Send(messageI)
			if err != nil {
				logging.Error.Printf("Error sending email: %s", err)
			}
		}(mg, message)
		return nil
	} else {
		_, _, err := mg.Send(message)
		return err
	}
}
