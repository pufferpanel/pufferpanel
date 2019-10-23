package impl

import (
	"github.com/mailgun/mailgun-go"
	"github.com/pufferpanel/apufferi/v4/logging"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/spf13/viper"
)

func SendEmailViaMailgun(to, subject, body string, async bool) error {
	domain := viper.GetString("email.domain")
	if domain == "" {
		return pufferpanel.ErrSettingNotConfigured("domain")
	}

	from := viper.GetString("email.from")
	if from == "" {
		return pufferpanel.ErrSettingNotConfigured("email.from")
	}

	key := viper.GetString("email.key")
	if key == "" {
		return pufferpanel.ErrSettingNotConfigured("email.key")
	}

	mg := mailgun.NewMailgun(domain, key)
	message := mg.NewMessage(from, subject, body, to)

	if async {
		go func(mgI *mailgun.MailgunImpl, messageI *mailgun.Message) {
			_, _, err := mgI.Send(messageI)
			if err != nil {
				logging.Build(logging.ERROR).WithMessage("error sending email").WithError(err).Log()
			}
		}(mg, message)
		return nil
	} else {
		_, _, err := mg.Send(message)
		return err
	}
}
