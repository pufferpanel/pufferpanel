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

package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/web/api"
	"github.com/pufferpanel/pufferpanel/web/auth"
	"github.com/pufferpanel/pufferpanel/web/oauth2"
)

func RegisterRoutes(e *gin.Engine) {
	e.LoadHTMLGlob("assets/web/**/*")

	api.Register(e.Group("/api"))
	e.Group("/assets").Static("", "assets/web")
	oauth2.Register(e.Group("/oauth2"))
	auth.Register(e.Group("/auth"))
}