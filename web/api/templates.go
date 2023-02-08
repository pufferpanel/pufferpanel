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
	"github.com/pufferpanel/pufferpanel/v3/middleware/panelmiddleware"
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

func getRepos(c *gin.Context) {
	db := panelmiddleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	repos, err := ts.GetRepos()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, repos)
}

func getsTemplatesForRepo(c *gin.Context) {
	db := panelmiddleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	templates, err := ts.GetAllFromRepo(c.Param("repo"))
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, templates)
}

func addRepo(c *gin.Context) {
}

func deleteRepo(c *gin.Context) {
}

// @Summary Value single template
// @Description Gets a template if registered
// @Accept json
// @Produce json
// @Success 200 {object} models.Template
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/templates [get]
func getTemplateFromRepo(c *gin.Context) {
	db := panelmiddleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	template, err := ts.Get(c.Param("repo"), c.Param("name"))
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
	db := panelmiddleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	templateName := c.Param("name")
	templateRequest := pufferpanel.Server{}
	err := c.MustBindWith(&templateRequest, binding.JSON)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	template, err := ts.Get("local", templateName)
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
	db := panelmiddleware.GetDatabase(c)
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
