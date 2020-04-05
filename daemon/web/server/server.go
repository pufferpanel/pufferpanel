/*
 Copyright 2016 Padduck, LLC

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

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/itsjamie/gin-cors"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/middleware"
	"github.com/pufferpanel/pufferpanel/v2/programs"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/satori/go.uuid"
	"github.com/spf13/cast"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RegisterRoutes(e *gin.RouterGroup) {
	l := e.Group("/server")
	{
		l.PUT("/:id", middleware.OAuth2Handler(pufferpanel.ScopeServersCreate, false), CreateServer)
		l.DELETE("/:id", middleware.OAuth2Handler(pufferpanel.ScopeServersDelete, true), DeleteServer)
		l.GET("/:id", middleware.OAuth2Handler(pufferpanel.ScopeServersEditAdmin, true), GetServerAdmin)
		l.POST("/:id", middleware.OAuth2Handler(pufferpanel.ScopeServersEditAdmin, true), EditServerAdmin)
		l.OPTIONS("/:id", response.CreateOptions("PUT", "DELETE", "GET", "POST"))

		l.GET("/:id/data", middleware.OAuth2Handler(pufferpanel.ScopeServersEdit, true), GetServer)
		l.POST("/:id/data", middleware.OAuth2Handler(pufferpanel.ScopeServersEdit, true), EditServer)
		l.OPTIONS("/:id/data", response.CreateOptions("GET", "POST"))

		l.POST("/:id/reload", middleware.OAuth2Handler(pufferpanel.ScopeServersEditAdmin, true), ReloadServer)
		l.OPTIONS("/:id/reload", response.CreateOptions("POST"))

		l.POST("/:id/start", middleware.OAuth2Handler(pufferpanel.ScopeServersStart, true), StartServer)
		l.OPTIONS("/:id/start", response.CreateOptions("POST"))

		l.POST("/:id/stop", middleware.OAuth2Handler(pufferpanel.ScopeServersStop, true), StopServer)
		l.OPTIONS("/:id/stop", response.CreateOptions("POST"))

		l.POST("/:id/kill", middleware.OAuth2Handler(pufferpanel.ScopeServersStop, true), KillServer)
		l.OPTIONS("/:id/kill", response.CreateOptions("POST"))

		l.POST("/:id/install", middleware.OAuth2Handler(pufferpanel.ScopeServersInstall, true), InstallServer)
		l.OPTIONS("/:id/install", response.CreateOptions("POST"))

		l.GET("/:id/file/*filename", middleware.OAuth2Handler(pufferpanel.ScopeServersFilesGet, true), GetFile)
		l.PUT("/:id/file/*filename", middleware.OAuth2Handler(pufferpanel.ScopeServersFilesPut, true), PutFile)
		l.DELETE("/:id/file/*filename", middleware.OAuth2Handler(pufferpanel.ScopeServersFilesPut, true), DeleteFile)
		l.POST("/:id/file/*filename", middleware.OAuth2Handler(pufferpanel.ScopeServersFilesPut, true), response.NotImplemented)
		l.OPTIONS("/:id/file/*filename", response.CreateOptions("GET", "PUT", "DELETE", "POST"))

		l.GET("/:id/console", middleware.OAuth2Handler(pufferpanel.ScopeServersConsole, true), GetLogs)
		l.POST("/:id/console", middleware.OAuth2Handler(pufferpanel.ScopeServersConsoleSend, true), PostConsole)
		l.OPTIONS("/:id/console", response.CreateOptions("GET", "POST"))

		l.GET("/:id/stats", middleware.OAuth2Handler(pufferpanel.ScopeServersStat, true), GetStats)
		l.OPTIONS("/:id/stats", response.CreateOptions("GET"))

		l.GET("/:id/status", middleware.OAuth2Handler(pufferpanel.ScopeServersView, true), GetStatus)
		l.OPTIONS("/:id/status", response.CreateOptions("GET"))
	}

	p := e.Group("/socket")
	{
		p.GET("/:id", middleware.OAuth2Handler(pufferpanel.ScopeServersConsole, true), cors.Middleware(cors.Config{
			Origins:     "*",
			Credentials: true,
		}), OpenSocket)
		p.Handle("CONNECT", "/:id", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "false")
		})
		p.OPTIONS("/:id", response.CreateOptions("GET"))
	}

	l.POST("", middleware.OAuth2Handler(pufferpanel.ScopeServersCreate, false), CreateServer)
	l.OPTIONS("", response.CreateOptions("POST"))
}

// @Summary Starts server
// @Description Starts the given server
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "Server started"
// @Success 202 {object} response.Empty "Start has been queued"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param wait query bool false "Wait for the operation to complete"
// @Router /server/{id}/start [post]
func StartServer(c *gin.Context) {
	item, _ := c.Get("server")
	server := item.(*programs.Program)

	_, wait := c.GetQuery("wait")

	if wait {
		err := server.Start()
		if response.HandleError(c, err, http.StatusInternalServerError) {
		} else {
			c.Status(http.StatusNoContent)
		}
	} else {
		go func() {
			err := server.Start()
			if err != nil {
				logging.Error().Printf("Error starting server %s: %s", server.Id(), err)
			}
		}()
		c.Status(http.StatusAccepted)
	}
}

// @Summary Stop server
// @Description Stops the given server
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "Server stopped"
// @Success 202 {object} response.Empty "Stop has been queued"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param wait query bool false "Wait for the operation to complete"
// @Router /server/{id}/stop [post]
func StopServer(c *gin.Context) {
	item, _ := c.Get("server")
	server := item.(*programs.Program)

	_, wait := c.GetQuery("wait")

	err := server.Stop()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if wait {
		err = server.GetEnvironment().WaitForMainProcess()
		if response.HandleError(c, err, http.StatusInternalServerError) {
		} else {
			c.Status(http.StatusNoContent)
		}
	} else {
		c.Status(204)
	}
}

// @Summary Kill server
// @Description Stops the given server forcefully
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "Server killed"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Router /server/{id}/stop [post]
func KillServer(c *gin.Context) {
	item, _ := c.Get("server")
	server := item.(*programs.Program)

	err := server.Kill()
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Create server
// @Description Creates the server
// @Accept json
// @Produce json
// @Success 200 {object} daemon.ServerIdResponse "Server created"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param server body shared.Server true "Server to create"
// @Router /server/{id} [put]
func CreateServer(c *gin.Context) {
	serverId := c.Param("id")
	if serverId == "" {
		id := uuid.NewV4()
		serverId = id.String()
	}
	prg, _ := programs.Get(serverId)

	if prg != nil {
		response.HandleError(c, pufferpanel.ErrServerAlreadyExists, http.StatusConflict)
		return
	}

	prg = &programs.Program{}
	err := json.NewDecoder(c.Request.Body).Decode(prg)

	if err != nil {
		logging.Error().Printf("Error decoding JSON body: %s", err)
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}

	prg.Identifier = serverId

	if err := programs.Create(prg); err != nil {
		response.HandleError(c, err, http.StatusInternalServerError)
	} else {
		c.JSON(200, &pufferpanel.ServerIdResponse{Id: serverId})
	}
}

// @Summary Deletes server
// @Description Deletes the given server
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "Server deleted"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Router /server/{id} [delete]
func DeleteServer(c *gin.Context) {
	item, _ := c.Get("server")
	prg := item.(*programs.Program)
	err := programs.Delete(prg.Id())
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Installs server
// @Description installs the given server
// @Accept json
// @Produce json
// @Success 202 {object} response.Empty "Install has been queued"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param wait query bool false "Wait for the operation to complete"
// @Router /server/{id}/install [post]
func InstallServer(c *gin.Context) {
	item, _ := c.Get("server")
	prg := item.(*programs.Program)

	_, wait := c.GetQuery("wait")

	if wait {
		err := prg.Install()
		if response.HandleError(c, err, http.StatusInternalServerError) {
		} else {
			c.Status(http.StatusNoContent)
		}
	} else {
		go func(p *programs.Program) {
			_ = p.Install()
		}(prg)

		c.Status(http.StatusAccepted)
	}
}

// @Summary Edit server data
// @Description Edits the given server data
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "Server edited"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param data body daemon.ServerData true "Server data"
// @Router /server/{id}/data [post]
func EditServer(c *gin.Context) {
	item, _ := c.Get("server")
	prg := item.(*programs.Program)

	data := &pufferpanel.ServerData{}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = prg.Edit(data.Variables, false)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Edit server data as admin
// @Description Edits the given server
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "Server edited"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param data body daemon.ServerData true "Server data"
// @Router /server/{id} [post]
func EditServerAdmin(c *gin.Context) {
	item, _ := c.Get("server")
	prg := item.(*programs.Program)

	data := &pufferpanel.ServerData{}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = prg.Edit(data.Variables, true)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Reload server
// @Description Reloads the server from disk
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "Reloaded server"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Router /server/{id}/reload [post]
func ReloadServer(c *gin.Context) {
	item, _ := c.Get("server")
	prg := item.(*programs.Program)

	err := programs.Reload(prg.Id())
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Gets server data
// @Description Gets the given server data
// @Accept json
// @Produce json
// @Success 200 {object} daemon.ServerData "Data for this server"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Router /server/{id}/data [get]
func GetServer(c *gin.Context) {
	item, _ := c.Get("server")
	server := item.(*programs.Program)

	data := server.GetData()

	c.JSON(200, &pufferpanel.ServerData{Variables: data})
}

// @Summary Gets server data as admin
// @Description Gets the given server data from an admin's view
// @Accept json
// @Produce json
// @Success 200 {object} daemon.ServerDataAdmin "Data for this server"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Router /server/{id} [get]
func GetServerAdmin(c *gin.Context) {
	item, _ := c.MustGet("server").(*pufferpanel.Server)

	c.JSON(200, &pufferpanel.ServerDataAdmin{Server: item})
}

// @Summary Get file/list
// @Description Gets a file or a file list from the server
// @Accept json
// @Produce json
// @Produce octet-stream
// @Success 200 {object} string "File"
// @Success 200 {object} messages.FileDesc "File List"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param filename path string true "File name"
// @Router /server/{id}/{filename} [get]
func GetFile(c *gin.Context) {
	item, _ := c.Get("server")
	server := item.(*programs.Program)

	targetPath := c.Param("filename")

	data, err := server.GetItem(targetPath)
	defer func() {
		if data != nil {
			pufferpanel.Close(data.Contents)
		}
	}()

	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(404)
		} else if err == pufferpanel.ErrIllegalFileAccess {
			response.HandleError(c, err, http.StatusBadRequest)
		} else {
			response.HandleError(c, err, http.StatusInternalServerError)
		}
		return
	}

	if data.FileList != nil {
		c.JSON(200, data.FileList)
	} else if data.Contents != nil {
		fileName := filepath.Base(data.Name)

		extraHeaders := map[string]string{
			"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, fileName),
		}

		//discard the built-in response, we cannot use this one at all
		c.DataFromReader(http.StatusOK, data.ContentLength, "application/octet-stream", data.Contents, extraHeaders)
	} else {
		//uhhhhhhhhhhhhh
		response.HandleError(c, errors.New("no file content or file list"), http.StatusInternalServerError)
	}
}

// @Summary Put file/folder
// @Description Puts a file or folder on the server
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "If file/folder was created"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param filename path string true "File name"
// @Param folder path bool true "If this is a folder"
// @Param file formData file false "File to place"
// @Router /server/{id}/{filename} [put]
func PutFile(c *gin.Context) {
	item, _ := c.Get("server")
	server := item.(*programs.Program)

	targetPath := c.Param("filename")

	if targetPath == "" {
		c.Status(404)
		return
	}

	var err error

	_, mkFolder := c.GetQuery("folder")
	if mkFolder {
		err = server.CreateFolder(targetPath)
		response.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	var sourceFile io.ReadCloser

	v := c.Request.Header.Get("Content-Type")
	if t, _, _ := mime.ParseMediaType(v); t == "multipart/form-data" {
		sourceFile, _, err = c.Request.FormFile("file")
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
	} else {
		sourceFile = c.Request.Body
	}

	file, err := server.OpenFile(targetPath)
	defer pufferpanel.Close(file)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		_, err = io.Copy(file, sourceFile)
		if response.HandleError(c, err, http.StatusInternalServerError) {
		} else {
			c.Status(http.StatusNoContent)
		}
	}
}

// @Summary Delete file
// @Description Deletes a file from the server
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "If file was deleted"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param filename path string true "File name"
// @Router /server/{id}/{filename} [delete]
func DeleteFile(c *gin.Context) {
	item, _ := c.Get("server")
	server := item.(*programs.Program)

	targetPath := c.Param("filename")

	err := server.DeleteItem(targetPath)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Run command
// @Description Runs a command in the server
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "If command was ran"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Param commands body string true "Command to run"
// @Router /server/{id}/console [post]
func PostConsole(c *gin.Context) {
	item, _ := c.Get("server")
	prg := item.(*programs.Program)

	d, _ := ioutil.ReadAll(c.Request.Body)
	cmd := string(d)
	err := prg.Execute(cmd)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Gets server stats
// @Description Gets the given server stats
// @Accept json
// @Produce json
// @Success 200 {object} daemon.ServerStats "Stats for this server"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Router /server/{id}/stats [get]
func GetStats(c *gin.Context) {
	item, _ := c.Get("server")
	svr := item.(*programs.Program)

	results, err := svr.GetEnvironment().GetStats()
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.JSON(200, results)
	}
}

// @Summary Gets server logs
// @Description Gets the given server logs since a certain time period
// @Accept json
// @Produce json
// @Success 200 {object} daemon.ServerLogs "Logs for this server"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Router /server/{id}/logs [get]
func GetLogs(c *gin.Context) {
	item, _ := c.Get("server")
	program := item.(*programs.Program)

	time := c.DefaultQuery("time", "0")

	castedTime, ok := cast.ToInt64E(time)
	if ok != nil {
		response.HandleError(c, pufferpanel.ErrInvalidUnixTime, http.StatusBadRequest)
		return
	}

	console, epoch := program.GetEnvironment().GetConsoleFrom(castedTime)
	msg := ""
	for _, k := range console {
		msg += k
	}

	c.JSON(200, &pufferpanel.ServerLogs{
		Epoch: epoch,
		Logs:  msg,
	})
}

// @Summary Gets server status
// @Description Gets the given server status
// @Accept json
// @Produce json
// @Success 200 {object} daemon.ServerRunning
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Empty
// @Failure 404 {object} response.Empty
// @Failure 500 {object} response.Error
// @Param id path string true "Server Identifier"
// @Router /server/{id}/status [get]
func GetStatus(c *gin.Context) {
	item, _ := c.Get("server")
	program := item.(*programs.Program)

	running, err := program.IsRunning()

	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.JSON(200, &pufferpanel.ServerRunning{Running: running})
	}
}

func OpenSocket(c *gin.Context) {
	item, _ := c.Get("server")
	program := item.(*programs.Program)

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	internalMap, _ := c.Get("scopes")
	scopes := internalMap.([]pufferpanel.Scope)

	socket := pufferpanel.Create(conn)

	go listenOnSocket(socket, program, scopes)

	program.GetEnvironment().AddListener(socket)
}
