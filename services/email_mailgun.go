package services

import (
	"github.com/mailgun/mailgun-go"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/errors"
	"github.com/spf13/viper"
)

func sendEmailViaMailgun(to, subject, body string, async bool) error {
	domain := viper.GetString("email.domain")
	if domain == "" {
		return errors.NewEmailNotConfigured("no domain defined")
	}

	from := viper.GetString("email.from")
	if from == "" {
		return errors.NewEmailNotConfigured("no from email defined")
	}

	key := viper.GetString("email.key")
	if key == "" {
		return errors.NewEmailNotConfigured("no api key defined")
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
