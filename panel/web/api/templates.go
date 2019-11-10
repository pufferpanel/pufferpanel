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
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/handlers"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/scope"
	"net/http"
)

func registerTemplates(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(scope.TemplatesView, false), getAllTemplates)
	g.Handle("OPTIONS", "", response.CreateOptions("GET"))
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
