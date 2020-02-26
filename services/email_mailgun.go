package services

import (
	"github.com/mailgun/mailgun-go"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/spf13/viper"
)

func SendEmailViaMailgun(to, subject, body string, async bool) error {
	domain := viper.GetString("panel.email.domain")
	if domain == "" {
		return pufferpanel.ErrSettingNotConfigured("domain")
	}

	from := viper.GetString("panel.email.from")
	if from == "" {
		return pufferpanel.ErrSettingNotConfigured("panel.email.from")
	}

	key := viper.GetString("panel.email.key")
	if key == "" {
		return pufferpanel.ErrSettingNotConfigured("panel.email.key")
	}

	mg := mailgun.NewMailgun(domain, key)
	message := mg.NewMessage(from, subject, "", to)
	message.SetHtml(body)

	if async {
		go func(mgI *mailgun.MailgunImpl, messageI *mailgun.Message) {
			_, _, err := mgI.Send(messageI)
			if err != nil {
				logging.Error().Printf("Error sending email: %s", err)
			}
		}(mg, message)
		return nil
	} else {
		_, _, err := mg.Send(message)
		return err
	}
}
