/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

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
