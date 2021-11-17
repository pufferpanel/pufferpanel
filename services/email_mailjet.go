package services

import (
	"github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/logging"
)

func SendEmailViaMailjet(to, subject, body string, async bool) error {
	domain := config.GetString("panel.email.domain")
	if domain == "" {
		return pufferpanel.ErrSettingNotConfigured("domain")
	}

	from := config.GetString("panel.email.from")
	if from == "" {
		return pufferpanel.ErrSettingNotConfigured("panel.email.from")
	}

	key := config.GetString("panel.email.key")
	if key == "" {
		return pufferpanel.ErrSettingNotConfigured("panel.email.key")
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
