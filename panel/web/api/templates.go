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
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/panel/models"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/handlers"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"net/http"
)

func registerTemplates(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(pufferpanel.ScopeTemplatesView, false), getAllTemplates)

	g.Handle("OPTIONS", "", response.CreateOptions("GET"))

	g.Handle("GET", "/:name", handlers.OAuth2Handler(pufferpanel.ScopeTemplatesView, false), getTemplate)
	g.Handle("PUT", "/:name", handlers.OAuth2Handler(pufferpanel.ScopeTemplatesEdit, false), putTemplate)
	g.Handle("OPTIONS", "/:name", response.CreateOptions("GET", "POST"))

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
// @Router /templates [get]
func getAllTemplates(c *gin.Context) {
	db := handlers.GetDatabase(c)
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
// @Router /templates [get]
func getTemplate(c *gin.Context) {
	db := handlers.GetDatabase(c)
	ts := &services.Template{DB: db}

	template, err := ts.Get(c.Param("name"))
	if err != nil && gorm.IsRecordNotFoundError(err) {
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
// @Success 204 {object} models.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /templates [get]
func putTemplate(c *gin.Context) {
	db := handlers.GetDatabase(c)
	ts := &services.Template{DB: db}

	templateName := c.Param("name")
	templateRequest := pufferpanel.Template{}
	err := c.MustBindWith(&templateRequest, binding.JSON)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	template, err := ts.Get(templateName)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		template = &models.Template{
			Name:     templateName,
		}
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	template.Template = templateRequest
	err = ts.Save(template)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}
