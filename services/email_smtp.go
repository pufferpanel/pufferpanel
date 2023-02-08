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
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"net/smtp"
	"strings"
)

func SendEmailViaSMTP(to, subject, body string, async bool) error {
	from := config.EmailFrom.Value()
	if from == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailFrom.Key())
	}

	host := config.EmailHost.Value()
	if host == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailHost.Key())
	}

	var auth smtp.Auth = nil

	if username := config.EmailUsername.Value(); username != "" {
		auth = smtp.PlainAuth("", username, config.EmailPassword.Value(), strings.Split(host, ":")[0])
	}

	data := []byte("Subject: " + subject + "\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" + body)

	if async {
		go func(host string, auth smtp.Auth, from, to string, data []byte) {
			err := smtp.SendMail(host, auth, from, []string{to}, data)
			if err != nil {
				logging.Error.Printf("Error sending email: %s", err)
			}
		}(host, auth, from, to, data)
		return nil
	} else {
		return smtp.SendMail(host, auth, from, []string{to}, data)
	}
}
