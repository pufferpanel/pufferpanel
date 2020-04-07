/*
 Copyright 2020 Padduck, LLC
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

package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/middleware"
	"github.com/pufferpanel/pufferpanel/v2/middleware/handlers"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/spf13/cast"
	"net/http"
	"strings"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	proxy := rg.Group("/daemon", handlers.HasOAuth2Token, middleware.NeedsDatabase)
	{
		g := proxy.Group("/server")
		{
			g.Any("/:id", proxyServerRequest)
			g.Any("/:id/*path", proxyServerRequest)
		}

		g = proxy.Group("/socket")
		{
			g.Any("/:id", proxyServerRequest)
		}
	}
}

func proxyServerRequest(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ss := &services.Server{DB: db}
	ns := &services.Node{DB: db}
	ps := &services.Permission{DB: db}

	serverId := c.Param("id")
	if serverId == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	path := strings.TrimPrefix(c.Request.URL.Path, "/proxy")

	s, err := ss.Get(serverId)
	if err != nil && !gorm.IsRecordNotFoundError(err) && response.HandleError(c, err, http.StatusInternalServerError) {
		return
	} else if s == nil || s.Identifier == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	token := c.MustGet("token").(*pufferpanel.Token)

	//if a session-token, we need to convert it to an oauth2 token instead
	if token.Claims.Audience == "session" {
		newToken, err := ps.GenerateOAuthForUser(cast.ToUint(token.Claims.Subject), &s.Identifier)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		//set new header
		c.Request.Header.Set("Authorization", "Bearer "+newToken)
	}

	if s.Node.IsLocal() {
		c.Request.URL.Path = path
		pufferpanel.Engine.HandleContext(c)
	} else {
		if c.IsWebsocket() {
			proxySocketRequest(c, path, ns, &s.Node)
		} else {
			proxyHttpRequest(c, path, ns, &s.Node)
		}
	}
	c.Abort()
}

func proxyHttpRequest(c *gin.Context, path string, ns *services.Node, node *models.Node) {
	callResponse, err := ns.CallNode(node, c.Request.Method, path, c.Request.Body, c.Request.Header)

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	//Even though apache isn't going to be in place, we can't set certain headers
	newHeaders := make(map[string]string, 0)
	for k, v := range callResponse.Header {
		switch k {
		case "Transfer-Encoding":
		case "Content-Type":
		case "Content-Length":
			continue
		default:
			newHeaders[k] = strings.Join(v, ", ")
		}
	}

	c.DataFromReader(callResponse.StatusCode, callResponse.ContentLength, callResponse.Header.Get("Content-Type"), callResponse.Body, newHeaders)
	c.Abort()
}

func proxySocketRequest(c *gin.Context, path string, ns *services.Node, node *models.Node) {
	if node.IsLocal() {
		//have gin handle the request again, but send it to daemon instead
		pufferpanel.Engine.HandleContext(c)
	} else {
		err := ns.OpenSocket(node, path, c.Writer, c.Request)
		response.HandleError(c, err, http.StatusInternalServerError)
	}
	c.Abort()
}
