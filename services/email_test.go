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
	"github.com/spf13/viper"
	"testing"
)

func Test_emailService_SendEmail(t *testing.T) {
	type args struct {
		to       string
		template string
		data     map[string]interface{}
		async    bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Account Creation",
			args: args{
				to: "test@example.com",
				template: "accountCreation",
				data: nil,
				async: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("email.templates", "assets/email/emails.json")
			viper.Set("email.provider", "debug")
			viper.Set("settings.companyName", "PufferPanel")
			viper.Set("settings.masterUrl", "http://localhost:8080")

			LoadEmailService()
			es := GetEmailService()
			if err := es.SendEmail(tt.args.to, tt.args.template, tt.args.data, tt.args.async); (err != nil) != tt.wantErr {
				t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}