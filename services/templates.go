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
	"archive/zip"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const ReleaseUrl = "https://github.com/PufferPanel/templates/archive/v2.zip"

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

func (t *Template) ImportFromRepo() error {
	dir, err := ioutil.TempDir("", "pufferpaneltemplates")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	targetFile, err := ioutil.TempFile("", "pufferpaneltemplates*.zip")
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		pufferpanel.Close(f)
		os.Remove(f.Name())
	}(targetFile)

	logging.Info().Printf("Downloading %s\n", ReleaseUrl)
	response, err := http.Get(ReleaseUrl)
	if err != nil {
		return err
	}
	defer pufferpanel.Close(response.Body)

	_, err = io.Copy(targetFile, response.Body)
	if err != nil {
		return err
	}
	_ = response.Body.Close()

	err = unzip(targetFile.Name(), dir)
	if err != nil {
		return err
	}

	return filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}

		readmePath := filepath.Join(filepath.Dir(p), "README.md")

		logging.Info().Printf("Importing %s\n", p)
		name := strings.TrimSuffix(filepath.Base(info.Name()), filepath.Ext(info.Name()))
		err = t.ImportTemplate(name, p, readmePath)
		return err
	})
}

func (t *Template) ImportTemplate(name, templatePath, readmePath string) error {
	template, err := openTemplate(templatePath)

	if err != nil {
		return err
	}

	if name == "" {
		name = strings.TrimSuffix(filepath.Base(templatePath), filepath.Ext(templatePath))
	}

	model := &models.Template{
		Template: template,
		Name:     name,
		Readme:   "",
	}

	if readmePath != "" {
		data, err := openReadme(readmePath)
		if err == nil {
			model.Readme = data
		}
	}

	return t.Save(model)
}

func openTemplate(path string) (t pufferpanel.Template, err error) {
	file, err := os.Open(path)
	defer pufferpanel.Close(file)
	if err != nil {
		return
	}

	err = json.NewDecoder(file).Decode(&t)
	return
}

func openReadme(path string) (string, error) {
	file, err := os.Open(path)
	defer pufferpanel.Close(file)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), err
}

func unzip(sourceZip, targetDir string) error {
	zipFile, err := zip.OpenReader(sourceZip)
	defer pufferpanel.Close(zipFile)
	if err != nil {
		return err
	}

	for _, f := range zipFile.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.HasPrefix(filepath.Base(f.Name), ".") {
			continue
		}

		logging.Info().Printf("Extracting %s\n", f.Name)
		exportPath := filepath.Join(targetDir, f.Name)
		err := os.MkdirAll(filepath.Dir(exportPath), 0644)
		if err != nil {
			return err
		}
		err = writeFile(f, exportPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(source *zip.File, target string) error {
	s, err := source.Open()
	defer pufferpanel.Close(s)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	defer pufferpanel.Close(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, s)
	return err
}