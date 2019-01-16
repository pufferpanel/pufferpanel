package services

import (
	"github.com/mailgun/mailgun-go"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/config"
	"github.com/pufferpanel/pufferpanel/errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type EmailService interface {
	SendEmail(to string, subject string, template string, data interface{}, async bool) error
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
			logging.Error("Error processing email template "+tmpl, err)
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

func (es *emailService) SendEmail(to, subject, template string, data interface{}, async bool) (err error) {
	tmpl := es.templates[template]

	if tmpl == nil {
		return errors.New("no template with name " + template)
	}

	builder := &strings.Builder{}

	err = tmpl.Execute(builder, data)
	if err != nil {
		return err
	}

	cfg := config.Get()

	domain, exist := cfg.GetString("email.domain")
	if !exist {
		return errors.NewEmailNotConfigured("no domain defined")
	}

	from, exist := cfg.GetString("email.from")
	if !exist {
		return errors.NewEmailNotConfigured("no from email defined")
	}

	key, exist := cfg.GetString("email.key")
	if !exist {
		return errors.NewEmailNotConfigured("no api key defined")
	}

	mg := mailgun.NewMailgun(domain, key)
	message := mg.NewMessage(from, subject, builder.String(), to)

	if async {
		go func(mgI *mailgun.MailgunImpl, messageI *mailgun.Message) {
			mgI.Send(messageI)
		}(mg, message)
	} else {
		_, _, err = mg.Send(message)
	}

	return
}