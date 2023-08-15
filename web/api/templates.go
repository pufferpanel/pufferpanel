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

package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
	"net/http"
)

func registerTemplates(g *gin.RouterGroup) {
	g.Handle("GET", "/", middleware.RequiresPermission(pufferpanel.ScopeTemplatesView, false), getRepos)
	g.Handle("OPTIONS", "/", response.CreateOptions("GET"))

	g.Handle("GET", "/:repo", middleware.RequiresPermission(pufferpanel.ScopeTemplatesView, false), getsTemplatesForRepo)
	g.Handle("PUT", "/:repo", middleware.RequiresPermission(pufferpanel.ScopeTemplatesEdit, false), addRepo)
	g.Handle("DELETE", "/:repo", middleware.RequiresPermission(pufferpanel.ScopeTemplatesEdit, false), deleteRepo)
	g.Handle("OPTIONS", "/:repo", response.CreateOptions("GET", "PUT", "DELETE"))

	g.Handle("GET", "/:repo/:name", middleware.RequiresPermission(pufferpanel.ScopeTemplatesView, false), getTemplateFromRepo)
	g.Handle("DELETE", "/local/:name", middleware.RequiresPermission(pufferpanel.ScopeTemplatesEdit, false), deleteTemplate)
	g.Handle("PUT", "/local/:name", middleware.RequiresPermission(pufferpanel.ScopeTemplatesEdit, false), putTemplate)
	g.Handle("OPTIONS", "/:repo/:name", response.CreateOptions("GET"))
	g.Handle("OPTIONS", "/local/:name", response.CreateOptions("GET", "DELETE", "PUT"))
}

// @Summary Get all repos
// @Description Gets all repos that are available to pull template from
// @Success 200 {object} []models.TemplateRepo
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/templates [get]
func getRepos(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	repos, err := ts.GetRepos()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, repos)
}

// @Summary Get all templates from repo
// @Description Gets all templates from a repository
// @Param repo path string true "Repo name"
// @Success 200 {object} []models.Template
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/templates/{repo} [get]
// @Security OAuth2Application[none]
func getsTemplatesForRepo(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	templates, err := ts.GetAllFromRepo(c.Param("repo"))
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, templates)
}

func addRepo(c *gin.Context) {
	//TODO: Implement
	response.NotImplemented(c)
}

func deleteRepo(c *gin.Context) {
	//TODO: Implement
	response.NotImplemented(c)
}

// @Summary Get template
// @Description Gets a template from the repo
// @Param repo path string true "Repo name"
// @Param template path string true "Template name"
// @Success 200 {object} models.Template
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/templates/{repo}/{template} [get]
// @Security OAuth2Application[none]
func getTemplateFromRepo(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	template, err := ts.Get(c.Param("repo"), c.Param("name"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(404)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, template)
}

// @Summary Adds or updates a template
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param template body pufferpanel.Server true "Template"
// @Param name path string true "Template name"
// @Router /api/templates/local/{name} [put]
// @Security OAuth2Application[none]
func putTemplate(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	templateName := c.Param("name")
	templateRequest := pufferpanel.Server{}
	err := c.MustBindWith(&templateRequest, binding.JSON)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	template, err := ts.Get("local", templateName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		template = &models.Template{
			Name: templateName,
		}
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	template.Server = templateRequest
	err = ts.Save(template)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Deletes template
// @Description Deletes template
// @Success 204 {object} nil
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param name path string true "Template name"
// @Router /api/templates/local/{name} [delete]
// @Security OAuth2Application[none]
func deleteTemplate(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	err := ts.Delete(c.Param("name"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}
