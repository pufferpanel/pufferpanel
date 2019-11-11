package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/panel/services/impl"
	"github.com/spf13/viper"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type EmailService interface {
	SendEmail(to string, template string, data map[string]interface{}, async bool) error
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

	jsonPath := viper.GetString("panel.email.templates")
	parentDir := filepath.Dir(jsonPath)
	emailDefinition, err := os.Open(jsonPath)
	if err != nil {
		panic(err.Error())
	}
	defer pufferpanel.Close(emailDefinition)

	var mapping map[string]*emailDeclaration
	err = json.NewDecoder(emailDefinition).Decode(&mapping)
	if err != nil {
		panic(err.Error())
	}

	for templateName, data := range mapping {
		subjectTemplate, err := template.New(templateName + "-subject").Parse(data.Subject)
		if err != nil {
			panic(errors.New(fmt.Sprintf("Error processing email template subject %s: %s", templateName, err.Error())))
		}

		body, err := ioutil.ReadFile(filepath.Join(parentDir, data.Body))
		if err != nil {
			panic(errors.New(fmt.Sprintf("Error processing email template subject %s: %s", templateName, err.Error())))
		}

		renderedTemplate, err := template.New(templateName + "-body").Parse(string(body))
		if err != nil {
			panic(errors.New(fmt.Sprintf("Error processing email template body %s: %s", templateName, err.Error())))
		}

		globalEmailService.templates[templateName] = &emailTemplate{
			Subject: subjectTemplate,
			Body:    renderedTemplate,
		}
	}

	var logger = logging.Debug()
	for k := range globalEmailService.templates {
		logger.Printf("Email template registered: %s", k)
	}
}

func GetEmailService() EmailService {
	return globalEmailService
}

func (es *emailService) SendEmail(to, template string, data map[string]interface{}, async bool) (err error) {
	provider := viper.GetString("panel.email.provider")
	if provider == "" {
		return pufferpanel.ErrEmailNotConfigured
	}

	tmpl := es.templates[template]

	if tmpl == nil {
		return pufferpanel.ErrNoTemplate(template)
	}

	if data == nil {
		data = make(map[string]interface{})
	}

	data["COMPANY_NAME"] = viper.GetString("panel.settings.companyName")
	data["MASTER_URL"] = viper.GetString("panel.settings.masterUrl")

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

	logging.Debug().Printf("Sending email to %s using %s", to, provider)

	switch provider {
	case "mailgun":
		return impl.SendEmailViaMailgun(to, subjectBuilder.String(), bodyBuilder.String(), async)
	case "debug":
		return nil
	default:
		return pufferpanel.ErrServiceInvalidProvider("email", provider)
	}
}
