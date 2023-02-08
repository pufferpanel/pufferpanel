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
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/middleware/handlers"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
	"net/http"
)

var client = http.Client{}

func registerTemplates(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(pufferpanel.ScopeTemplatesView, false), getAllTemplates)
	g.Handle("OPTIONS", "", response.CreateOptions("GET"))

	g.Handle("POST", "/import", handlers.OAuth2Handler(pufferpanel.ScopeTemplatesEdit, false), getImportableTemplates)
	g.Handle("POST", "/import/:name", handlers.OAuth2Handler(pufferpanel.ScopeTemplatesEdit, false), importTemplate)

	g.Handle("GET", "/:name", handlers.OAuth2Handler(pufferpanel.ScopeTemplatesView, false), getTemplate)
	g.Handle("DELETE", "/:name", handlers.OAuth2Handler(pufferpanel.ScopeTemplatesView, false), deleteTemplate)
	g.Handle("PUT", "/:name", handlers.OAuth2Handler(pufferpanel.ScopeTemplatesEdit, false), putTemplate)
	g.Handle("OPTIONS", "/:name", response.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

// @Summary Get templates
// @Description Gets all templates registered
// @Accept json
// @Produce json
// @Success 200 {object} models.Templates
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/templates [get]
func getAllTemplates(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	templates, err := ts.GetAll()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, templates)
}

// @Summary Get single template
// @Description Gets a template if registered
// @Accept json
// @Produce json
// @Success 200 {object} models.Template
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/templates [get]
func getTemplate(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	template, err := ts.Get(c.Param("name"))
	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatus(404)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, template)
}

// @Summary Adds or updates a template
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param template body pufferpanel.Server true "Template"
// @Param name path string true "Template name"
// @Router /api/templates/{name} [put]
func putTemplate(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	templateName := c.Param("name")
	templateRequest := pufferpanel.Server{}
	err := c.MustBindWith(&templateRequest, binding.JSON)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	template, err := ts.Get(templateName)
	if err != nil && err == gorm.ErrRecordNotFound {
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

// @Summary Import template from repo
// @Description Imports the given template from our main repo
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param name path string true "Template"
// @Router /api/templates/import/{name} [post]
func importTemplate(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	err := ts.ImportFromRepo(c.Param("name"))
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Deletes template
// @Description Deletes template
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param name path string true "Template"
// @Router /api/templates/{name} [delete]
func deleteTemplate(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	err := ts.Delete(c.Param("name"))
	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatus(404)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Gets importable templates
// @Description Gets all templates which can be imported from https://github.com/PufferPanel/templates
// @Accept json
// @Produce json
// @Success 200 {object} []string
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param template body pufferpanel.Server true "Template"
// @Router /api/templates/import [post]
func getImportableTemplates(c *gin.Context) {
	ts := &services.Template{}
	results, err := ts.GetImportableTemplates()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(200, results)
}
