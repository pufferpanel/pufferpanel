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
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var existingPaths = make(map[string]repoCache, 0)
var pathLock sync.Mutex

type Template struct {
	DB *gorm.DB
}

type repoCache struct {
	LastRefreshed time.Time
	Path          string
}

func (t *Template) GetRepos() ([]*models.TemplateRepo, error) {
	var repos []*models.TemplateRepo
	err := t.DB.Find(&repos).Error

	//return list from the db, and add local
	return append(repos, &models.TemplateRepo{Name: "local"}), err
}

func (t *Template) GetAllFromRepo(repo string) ([]*models.Template, error) {
	var templates []*models.Template
	var err error

	if repo == "local" {
		err = t.DB.Find(&templates).Error
		if err != nil {
			return nil, err
		}

		//because we don't want to return a ton of data, we'll only return a few select fields
		replacement := make([]*models.Template, len(templates))

		for k, v := range templates {
			replacement[k] = &models.Template{
				Name: v.Name,
				Server: pufferpanel.Server{
					Display: v.Server.Display,
					Type:    v.Server.Type,
				},
			}
		}

		templates = replacement
	} else {
		repoDb := &models.TemplateRepo{
			Name: repo,
		}

		err = t.DB.First(repoDb).Error
		if err != nil {
			return nil, err
		}

		path, err := validateRepoOnDisk(repoDb)
		if err != nil {
			return nil, err
		}

		folders, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, folder := range folders {
			if !folder.IsDir() || strings.HasPrefix(folder.Name(), ".") {
				continue
			}

			jsonFolder, err := os.ReadDir(filepath.Join(path, folder.Name()))
			if err != nil {
				return nil, err
			}

			for _, v := range jsonFolder {
				if !strings.HasSuffix(v.Name(), ".json") {
					continue
				}

				name := strings.TrimSuffix(v.Name(), ".json")
				templatePath := filepath.Join(path, folder.Name(), v.Name())
				template, err := readTemplateFromDisk(name, templatePath)
				if err != nil {
					logging.Error.Printf("Error reading template from %s: %s", templatePath, err.Error())
					continue
				}

				templates = append(templates, &models.Template{
					Name: name,
					Server: pufferpanel.Server{
						Display:     template.Server.Display,
						Type:        template.Server.Type,
						Environment: template.Server.Environment,
					},
				})
			}
		}
	}

	return templates, err
}

func (t *Template) Get(repo, name string) (*models.Template, error) {
	template := &models.Template{
		Name: name,
	}
	if repo == "local" {
		err := t.DB.First(template).Error
		if err != nil {
			return nil, err
		}
	} else {
		repoDb := &models.TemplateRepo{
			Name: repo,
		}

		err := t.DB.First(repoDb).Error
		if err != nil {
			return nil, err
		}

		path, err := validateRepoOnDisk(repoDb)
		if err != nil {
			return nil, err
		}

		var folderName string
		var exists bool
		folders, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, folder := range folders {
			if !folder.IsDir() || strings.HasPrefix(folder.Name(), ".") {
				continue
			}

			jsonFolder, err := os.ReadDir(filepath.Join(path, folder.Name()))
			if err != nil {
				return nil, err
			}

			for _, v := range jsonFolder {
				if v.Name() == name+".json" {
					folderName = folder.Name()
					exists = true
					break
				}
			}
		}

		if !exists {
			return nil, pufferpanel.ErrNoTemplate(name)
		}

		templatePath := filepath.Join(path, folderName, name+".json")
		template, err = readTemplateFromDisk(name, templatePath)
		if err != nil {
			return nil, err
		}

		readmePath := filepath.Join(path, folderName, "README.md")
		readme, err := os.ReadFile(readmePath)
		if err != nil {
			logging.Error.Printf("Error reading readme %s: %s", readmePath, err.Error())
		} else {
			template.Readme = string(readme)
		}
	}

	return template, nil
}

func (t *Template) Save(template *models.Template) error {
	return t.DB.Save(template).Error
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

func readTemplateFromDisk(name, path string) (*models.Template, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	template := &models.Template{
		Name: name,
	}
	err = json.NewDecoder(file).Decode(template)
	return template, err
}

func validateRepoOnDisk(repo *models.TemplateRepo) (string, error) {
	//temp locations!!!
	pathLock.Lock()
	defer pathLock.Unlock()

	if cache, exists := existingPaths[repo.Name]; exists {
		if cache.LastRefreshed.Add(time.Minute * 5).After(time.Now()) {
			return cache.Path, nil
		}

		logging.Debug.Printf("Updating local git repo for %s: %s", repo, cache)

		r, err := git.PlainOpen(cache.Path)
		if err != nil {
			return "", err
		}

		w, err := r.Worktree()
		if err != nil {
			return "", err
		}

		err = w.Pull(&git.PullOptions{
			SingleBranch:  true,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + repo.Branch),
			RemoteName:    "origin"},
		)
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return "", err
		}

		existingPaths[repo.Name] = repoCache{
			LastRefreshed: time.Now(),
			Path:          cache.Path,
		}
	} else {
		path := filepath.Join(config.CacheFolder.Value(), "template-repos", repo.Name)

		//if the directory already exists, we may need to nuke it
		fi, err := os.Stat(path)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		} else if fi != nil {
			_ = os.RemoveAll(path)
		}

		err = os.MkdirAll(path, 0755)
		if err != nil && !os.IsExist(err) {
			return "", err
		}

		logging.Debug.Printf("Checking out repo %s: %s", repo, path)
		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL:           repo.Url,
			SingleBranch:  true,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + repo.Branch),
		})
		if err != nil {
			return "", err
		}

		existingPaths[repo.Name] = repoCache{
			LastRefreshed: time.Now(),
			Path:          path,
		}
	}

	return existingPaths[repo.Name].Path, nil
}
