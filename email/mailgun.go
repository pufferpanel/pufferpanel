package email

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
)

type mailgunProvider struct {
	Provider
}

func init() {
	providers["mailgun"] = mailgunProvider{}
}

func (mailgunProvider) Send(to, subject, body string) error {
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

	mgapi := mailgun.NewMailgun(domain, key)
	message := mgapi.NewMessage(from, subject, "", to)
	message.SetHtml(body)

	_, _, err := mgapi.Send(context.Background(), message)
	return err
}
