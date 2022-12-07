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
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var ignoredPaths = []string{".git", ".github"}

var templateRepo *git.Repository
var templateFiles billy.Filesystem
var templateRepoExpiredAt time.Time
var templateCacheLocker sync.Mutex

type Template struct {
	DB *gorm.DB
}

type templateFile struct {
	Name       string
	Path       string
	ReadmePath string
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
	files, err := t.getTemplateFiles()
	if err != nil {
		return err
	}

	templateCacheLocker.Lock()
	defer templateCacheLocker.Unlock()

	var template io.ReadCloser
	var readme io.ReadCloser

	defer pufferpanel.Close(template)
	defer pufferpanel.Close(readme)

	for _, file := range files {
		if file.Name == templateName {
			template, err = templateFiles.Open(file.Path)
			if err != nil {
				return err
			}
			if file.ReadmePath != "" {
				readme, err = templateFiles.Open(file.ReadmePath)
				if err != nil {
					return err
				}
			}
			break
		}
	}

	if template == nil {
		return pufferpanel.ErrNoTemplate(templateName)
	}

	return t.ImportTemplate(templateName, template, readme)
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
	files, err := t.getTemplateFiles()
	if err != nil {
		return nil, err
	}

	results := make([]string, 0)
	for _, f := range files {
		results = append(results, f.Name)
	}

	return results, nil
}

func (t *Template) refreshGithub() (err error) {
	templateCacheLocker.Lock()
	defer templateCacheLocker.Unlock()

	if templateRepo == nil {
		client.InstallProtocol("https", githttp.NewClient(pufferpanel.Http()))

		templateFiles = memfs.New()
		templateRepo, err = git.Clone(memory.NewStorage(), templateFiles, &git.CloneOptions{
			URL:           "https://github.com/PufferPanel/templates",
			ReferenceName: "refs/heads/v2.5",
		})

		if err != nil {
			return
		}
	}

	if templateRepoExpiredAt.Before(time.Now()) {
		err = templateRepo.Fetch(&git.FetchOptions{})
		if err == git.NoErrAlreadyUpToDate {
			err = nil
		}

		if err != nil {
			return
		}
	}

	templateRepoExpiredAt = time.Now().Add(time.Minute * 5)
	return
}

func (t *Template) getTemplateFiles() ([]templateFile, error) {
	err := t.refreshGithub()
	if err != nil {
		return nil, err
	}

	templateCacheLocker.Lock()
	defer templateCacheLocker.Unlock()

	directories, err := templateFiles.ReadDir("/")
	if err != nil {
		return nil, err
	}

	results := make([]templateFile, 0)

	for _, directory := range directories {
		if !directory.IsDir() || pufferpanel.ContainsString(ignoredPaths, directory.Name()) {
			continue
		}

		var files []os.FileInfo
		files, err = templateFiles.ReadDir(directory.Name())
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			var readme string

			folder := filepath.Join("/", directory.Name())

			testPath := filepath.Join(folder, strings.TrimSuffix(file.Name(), ".json")+".md")
			if fi, err := templateFiles.Stat(testPath); err == nil && !fi.IsDir() {
				readme = testPath
			} else {
				testPath = filepath.Join(folder, "README.md")
				if fi, err = templateFiles.Stat(testPath); err == nil && !fi.IsDir() {
					readme = testPath
				}
			}
			results = append(results, templateFile{
				Name:       strings.TrimSuffix(file.Name(), ".json"),
				Path:       filepath.Join("/", directory.Name(), file.Name()),
				ReadmePath: readme,
			})
		}
	}

	return results, nil
}
