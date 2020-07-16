/*
 Copyright 2019 Padduck, LLC
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
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const TemplateJson = "https://raw.githubusercontent.com/PufferPanel/templates/v2/{name}/{name}.json"
const TemplateReadme = "https://raw.githubusercontent.com/PufferPanel/templates/v2/{name}/README.md"

var templateClient = http.Client{}

type Template struct {
	DB *gorm.DB
}

func (t *Template) GetAll() (*models.Templates, error) {
	templates := &models.Templates{}
	err := t.DB.Find(&templates).Error
	if err != nil {
		return nil, err
	}
	//because we don't want to return a ton of data, we'll only return a few select fields
	replacement := make(models.Templates, len(*templates))

	for k, v := range *templates {
		replacement[k] = &models.Template{
			Name: v.Name,
			Template: pufferpanel.Template{
				Server: pufferpanel.Server{
					Display: v.Template.Server.Display,
					Type:    v.Template.Server.Type,
				},
			},
		}
	}
	return &replacement, err
}

func (t *Template) Get(name string) (*models.Template, error) {
	template := &models.Template{
		Name:  name,
	}
	err := t.DB.Find(&template).Error
	if err != nil {
		return nil, err
	}
	return template, err
}


func (t *Template) Save(template *models.Template) error {
	return t.DB.Save(template).Error
}

func (t *Template) ImportFromRepo(templateName string) error {
	url := strings.Replace(TemplateJson, "{name}", templateName, -1)

	response, err := templateClient.Get(url)
	if err != nil {
		return err
	}
	defer pufferpanel.Close(response.Body)

	var readme io.Reader
	readmeUrl := strings.Replace(TemplateReadme, "{name}", templateName, -1)
	readmeResponse, err := templateClient.Get(readmeUrl)
	if err == nil {
		defer pufferpanel.Close(readmeResponse.Body)
		readme = readmeResponse.Body
	}

	return t.ImportTemplate(templateName, response.Body, readme)
}

func (t *Template) ImportTemplate(name string, template, readme io.Reader) error {
	var templateData pufferpanel.Template
	err := json.NewDecoder(template).Decode(&templateData)
	if err != nil {
		return err
	}

	model := &models.Template{
		Template: templateData,
		Name:     name,
	}

	if readme != nil {
		data, err := ioutil.ReadAll(readme)
		if err != nil {
			return err
		}
		model.Readme = string(data)
	}

	return t.Save(model)
}
