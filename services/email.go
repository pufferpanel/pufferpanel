package services

import (
	"github.com/pufferpanel/apufferi/logging"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type EmailService interface {
	SendEmail(to string, subject string, template string)
}

var globalEmailService *emailService

type emailService struct {
	templates map[string]*template.Template
}

func LoadEmailService() {
	globalEmailService = &emailService{templates: make(map[string]*template.Template)}

	//validate all emails in the email folder are valid and register templates
	prefix := "assets" + string(os.PathSeparator) + "email" + string(os.PathSeparator)
	templates, err := filepath.Glob(prefix + "*.html")
	if err != nil {
		panic(err.Error())
	}

	for _, tmpl := range templates {
		templateName := strings.TrimSuffix(strings.TrimPrefix(tmpl, prefix), ".html")
		renderedTemplate, err := template.New(templateName).ParseFiles(tmpl)
		if err != nil {
			logging.Error("Error processing email template " + tmpl, err)
			continue
		}
		globalEmailService.templates[templateName] = renderedTemplate
	}

	for k := range globalEmailService.templates {
		logging.Debugf("Email template registered: %s", k)
	}
}

func GetEmailService() EmailService {
	return globalEmailService
}

func (es *emailService) SendEmail(to, subject, template string) {

}
