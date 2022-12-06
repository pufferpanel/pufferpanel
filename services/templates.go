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
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const TemplateGithubUrl = "https://api.github.com/repos/pufferpanel/templates/git/trees/v2.5?recursive=true"

var ignoredPaths = []string{".git", ".github"}
var templateCache *pufferpanel.GithubFileList
var templateCacheExpireAt time.Time
var templateCacheLocker sync.Mutex

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
			Server: pufferpanel.Server{
				Display:               v.Server.Display,
				Type:                  v.Server.Type,
				SupportedEnvironments: v.Server.SupportedEnvironments,
			},
		}
	}
	return &replacement, err
}

func (t *Template) Get(name string) (*models.Template, error) {
	template := &models.Template{
		Name: name,
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
	files, err := t.callGithub()
	if err != nil {
		return err
	}

	var templateUrl string
	var readmeUrl string

	var root string
	for _, file := range files.Tree {
		if strings.HasSuffix(file.Path, templateName+".json") {
			templateUrl = file.Url
			root = strings.TrimSuffix(file.Path, "/"+templateName+".json")
		}
		if strings.HasSuffix(file.Path, templateName+".md") {
			readmeUrl = file.Url
		}
	}
	//re-search, look for readme if not one for the specific template
	if readmeUrl == "" {
		for _, file := range files.Tree {
			if file.Path == root+"/README.md" {
				readmeUrl = file.Url
			}
		}
	}

	response, err := templateClient.Get(templateUrl)
	if err != nil {
		return err
	}
	defer pufferpanel.Close(response.Body)

	var readme io.Reader
	if readmeUrl != "" {
		readmeResponse, err := templateClient.Get(readmeUrl)
		if err == nil && readmeResponse.StatusCode == 200 {
			defer pufferpanel.Close(readmeResponse.Body)
			readme = readmeResponse.Body
		}
	}

	return t.ImportTemplate(templateName, response.Body, readme)
}

func (t *Template) ImportTemplate(name string, template, readme io.Reader) error {
	var templateData pufferpanel.Server
	err := json.NewDecoder(template).Decode(&templateData)
	if err != nil {
		return err
	}

	if len(templateData.SupportedEnvironments) == 0 {
		templateData.SupportedEnvironments = []interface{}{templateData.Environment}
	}

	model := &models.Template{
		Server: templateData,
		Name:   strings.ToLower(name),
	}

	if readme != nil {
		data, err := io.ReadAll(readme)
		if err != nil {
			return err
		}
		model.Readme = string(data)
	}

	return t.Save(model)
}

func (t *Template) Delete(name string) error {
	model := &models.Template{
		Name: name,
	}

	res := t.DB.Delete(model)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (t *Template) GetImportableTemplates() ([]string, error) {
	d, err := t.callGithub()
	if err != nil {
		return nil, err
	}

	results := make([]string, 0)
	for _, v := range d.Tree {
		if v.Type == "blob" && strings.HasSuffix(v.Path, ".json") {
			ignore := false
			for _, i := range ignoredPaths {
				if strings.HasPrefix(v.Path, i) {
					ignore = true
				}
			}
			if ignore {
				continue
			}
			results = append(results, v.Path)
		}
	}
	return results, nil
}

func (t *Template) callGithub() (pufferpanel.GithubFileList, error) {
	templateCacheLocker.Lock()
	defer templateCacheLocker.Unlock()

	if templateCache != nil && templateCacheExpireAt.After(time.Now()) {
		return *templateCache, nil
	}

	var d pufferpanel.GithubFileList
	u, err := url.Parse(TemplateGithubUrl)
	if err != nil {
		return d, err
	}

	request := &http.Request{
		Method: "GET",
		URL:    u,
	}
	if request.Header == nil {
		request.Header = http.Header{}
	}

	request.Header.Add("Accept", "application/vnd.github.v3+json")

	res, err := pufferpanel.Http().Do(request)
	if err != nil {
		return d, err
	}
	defer pufferpanel.Close(res.Body)

	err = json.NewDecoder(res.Body).Decode(&d)

	templateCache = &d
	templateCacheExpireAt = time.Now().Add(time.Minute * 5)
	return d, err
}
