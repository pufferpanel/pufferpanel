/*
 Copyright 2018 Padduck, LLC
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
	"github.com/pufferpanel/pufferpanel/shared"
)

func registerUsers(g *gin.RouterGroup) {
	g.Handle("GET", "", shared.NotImplemented)
	g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("PUT", "/:username", shared.NotImplemented)
	g.Handle("GET", "/:username", shared.NotImplemented)
	g.Handle("POST", "/:username", shared.NotImplemented)
	g.Handle("DELETE", "/:username", shared.NotImplemented)
	g.Handle("OPTIONS", "/:username", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))
}