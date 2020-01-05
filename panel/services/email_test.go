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