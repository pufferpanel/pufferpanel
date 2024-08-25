package services

import (
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	emailAssets "github.com/pufferpanel/pufferpanel/v3/assets/email"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/email"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"html/template"
	"io/fs"
	"os"
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

	var merged fs.ReadFileFS
	if config.EmailTemplateFolder.Value() != "" {
		merged = pufferpanel.NewMergedFS(os.DirFS(config.EmailTemplateFolder.Value()), emailAssets.FS)
	} else {
		merged = emailAssets.FS
	}

	emailDefinition, err := merged.Open("emails.json")
	if err != nil {
		panic(err)
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
			panic(fmt.Errorf("error processing email template subject %s: %s", templateName, err.Error()))
		}

		body, err := merged.ReadFile(data.Body)
		if err != nil {
			panic(fmt.Errorf("error processing email template subject %s: %s", templateName, err.Error()))
		}

		renderedTemplate, err := template.New(templateName + "-body").Parse(string(body))
		if err != nil {
			panic(fmt.Errorf("error processing email template body %s: %s", templateName, err.Error()))
		}

		globalEmailService.templates[templateName] = &emailTemplate{
			Subject: subjectTemplate,
			Body:    renderedTemplate,
		}
	}

	var logger = logging.Debug
	for k := range globalEmailService.templates {
		logger.Printf("Email template registered: %s", k)
	}
}

func GetEmailService() EmailService {
	return globalEmailService
}

func (es *emailService) SendEmail(to, template string, data map[string]interface{}, async bool) (err error) {
	provider := config.EmailProvider.Value()
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

	data["COMPANY_NAME"] = config.CompanyName.Value()
	data["MASTER_URL"] = config.MasterUrl.Value()

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

	logging.Debug.Printf("Sending email to %s using %s", to, provider)

	svc := email.GetProvider(provider)
	if svc == nil {
		return pufferpanel.ErrServiceInvalidProvider("email", provider)
	}

	if async {
		go func(emailService email.Provider, toEmail, subject, body string) {
			err := emailService.Send(toEmail, subject, body)
			if err != nil {
				logging.Error.Printf("Error sending email: %s", err)
			}
		}(svc, to, subjectBuilder.String(), bodyBuilder.String())
		return nil
	} else {
		return svc.Send(to, subjectBuilder.String(), bodyBuilder.String())
	}
}
