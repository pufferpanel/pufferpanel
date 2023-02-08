package services

import (
	"github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

func SendEmailViaMailjet(to, subject, body string, async bool) error {
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

	m := mailjet.NewMailjetClient(domain, key)

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: from,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
				},
			},
			Subject:  subject,
			HTMLPart: body,
		},
	}
	message := mailjet.MessagesV31{Info: messagesInfo}

	if async {
		go func(mgI *mailjet.Client, messageI *mailjet.MessagesV31) {
			_, err := m.SendMailV31(messageI)
			if err != nil {
				logging.Error.Printf("Error sending email: %s", err)
			}
		}(m, &message)
		return nil
	} else {
		_, err := m.SendMailV31(&message)
		return err
	}
}
