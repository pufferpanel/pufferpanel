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
	"github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/web/handlers"
)

func registerTemplates(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2(pufferpanel.ScopeViewTemplates, false), getAllTemplates)
	g.Handle("OPTIONS", "", pufferpanel.CreateOptions("GET"))
}

func getAllTemplates(c *gin.Context) {
	res := response.From(c)

	db, err := database.GetConnection()
	if pufferpanel.HandleError(res, err) {
		return
	}
	ts := &services.Template{DB: db}
	templates, err := ts.GetAll()
	if pufferpanel.HandleError(res, err) {
		return
	}

	res.Data(templates)
}
