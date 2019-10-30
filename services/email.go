package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pufferpanel/apufferi/v4"
	"github.com/pufferpanel/apufferi/v4/logging"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/services/impl"
	"github.com/spf13/viper"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type EmailService interface {
	SendEmail(to string, template string, data interface{}, async bool) error
}

type emailTemplate struct {
	Subject *template.Template
	Body    *template.Template
}

type emailDeclaration struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

var globalEmailService *emailService

type emailService struct {
	templates map[string]*emailTemplate
}

func LoadEmailService() {
	globalEmailService = &emailService{templates: make(map[string]*emailTemplate)}

	jsonPath := viper.GetString("email.templates")
	parentDir := filepath.Dir(jsonPath)
	emailDefinition, err := os.Open(viper.GetString("email.templates"))
	if err != nil {
		panic(err.Error())
	}
	defer apufferi.Close(emailDefinition)

	var mapping map[string]*emailDeclaration
	err = json.NewDecoder(emailDefinition).Decode(&mapping)
	if err != nil {
		panic(err.Error())
	}

	for templateName, data := range mapping {
		subjectTemplate, err := template.New(templateName).Parse(data.Subject)
		if err != nil {
			panic(errors.New(fmt.Sprintf("Error processing email template subject %s: %s", templateName, err.Error())))
		}

		renderedTemplate, err := template.New(templateName).ParseFiles(filepath.Join(parentDir, data.Body))
		if err != nil {
			panic(errors.New(fmt.Sprintf("Error processing email template body %s: %s", templateName, err.Error())))
		}

		globalEmailService.templates[templateName] = &emailTemplate{
			Subject: subjectTemplate,
			Body:    renderedTemplate,
		}
	}

	for k := range globalEmailService.templates {
		logging.Debug("Email template registered: %s", k)
	}
}

func GetEmailService() EmailService {
	return globalEmailService
}

func (es *emailService) SendEmail(to, template string, data interface{}, async bool) (err error) {
	tmpl := es.templates[template]

	if tmpl == nil {
		return pufferpanel.ErrNoTemplate(template)
	}

	subjectBuilder := &strings.Builder{}
	err = tmpl.Subject.Execute(subjectBuilder, data)
	if err != nil {
		return err
	}

	bodyBuilder := &strings.Builder{}
	err = tmpl.Body.Execute(bodyBuilder, data)
	if err != nil {
		return err
	}

	provider := viper.GetString("email.provider")
	if provider == "" {
		return pufferpanel.ErrEmailNotConfigured
	}

	switch provider {
	case "mailgun":
		return impl.SendEmailViaMailgun(to, subjectBuilder.String(), bodyBuilder.String(), async)
	default:
		return pufferpanel.ErrServiceInvalidProvider("email", provider)
	}
}
