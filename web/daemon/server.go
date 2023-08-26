/*
 Copyright 2022 (c) PufferPanel

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

package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/servers"
	"github.com/spf13/cast"
	"io"
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

func RegisterServerRoutes(e *gin.RouterGroup) {
	l := e.Group("/server")
	{
		l.PUT("/:serverId", middleware.RequiresPermission(pufferpanel.ScopeServerCreate, false), createServer)
		l.DELETE("/:serverId", middleware.RequiresPermission(pufferpanel.ScopeServerDelete, true), deleteServer)
		l.OPTIONS("/:serverId", response.CreateOptions("PUT", "DELETE", "GET"))

		l.GET("/:serverId/definition", middleware.RequiresPermission(pufferpanel.ScopeServerViewAdmin, true), getServerAdmin)
		l.PUT("/:serverId/definition", middleware.RequiresPermission(pufferpanel.ScopeServerEditAdmin, true), editServerAdmin)
		l.OPTIONS("/:serverId/definition", response.CreateOptions("PUT", "DELETE", "GET"))

		l.GET("/:serverId/data", middleware.RequiresPermission(pufferpanel.ScopeServerViewData, true), getServerData)
		l.POST("/:serverId/data", middleware.RequiresPermission(pufferpanel.ScopeServerEditData, true), editServerData)
		l.OPTIONS("/:serverId/data", response.CreateOptions("GET", "POST"))

		l.GET("/:serverId/tasks", middleware.RequiresPermission(pufferpanel.ScopeServerTaskView, true), getServerTasks)
		l.OPTIONS("/:serverId/tasks", response.CreateOptions("GET"))

		l.GET("/:serverId/tasks/:taskId", middleware.RequiresPermission(pufferpanel.ScopeServerTaskView, true), getServerTask)
		l.PUT("/:serverId/tasks/:taskId", middleware.RequiresPermission(pufferpanel.ScopeServerTaskCreate, true), editServerTask)
		l.DELETE("/:serverId/tasks/:taskId", middleware.RequiresPermission(pufferpanel.ScopeServerTaskDelete, true), deleteServerTask)
		l.OPTIONS("/:serverId/tasks/:taskId", response.CreateOptions("GET", "PUT", "DELETE"))

		l.POST("/:serverId/tasks/:taskId/run", middleware.RequiresPermission(pufferpanel.ScopeServerTaskRun, true), runServerTask)
		l.OPTIONS("/:serverId/tasks/:taskId/run", response.CreateOptions("POST"))

		l.POST("/:serverId/reload", middleware.RequiresPermission(pufferpanel.ScopeServerReload, true), reloadServer)
		l.OPTIONS("/:serverId/reload", response.CreateOptions("POST"))

		l.POST("/:serverId/start", middleware.RequiresPermission(pufferpanel.ScopeServerStart, true), startServer)
		l.OPTIONS("/:serverId/start", response.CreateOptions("POST"))

		l.POST("/:serverId/stop", middleware.RequiresPermission(pufferpanel.ScopeServerStop, true), stopServer)
		l.OPTIONS("/:serverId/stop", response.CreateOptions("POST"))

		l.POST("/:serverId/kill", middleware.RequiresPermission(pufferpanel.ScopeServerKill, true), killServer)
		l.OPTIONS("/:serverId/kill", response.CreateOptions("POST"))

		l.POST("/:serverId/install", middleware.RequiresPermission(pufferpanel.ScopeServerInstall, true), installServer)
		l.OPTIONS("/:serverId/install", response.CreateOptions("POST"))

		l.GET("/:serverId/file/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileGet, true), getFile)
		l.PUT("/:serverId/file/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit, true), putFile)
		l.DELETE("/:serverId/file/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit, true), deleteFile)
		l.POST("/:serverId/file/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit, true), response.NotImplemented)
		l.OPTIONS("/:serverId/file/*filename", response.CreateOptions("GET", "PUT", "DELETE", "POST"))

		l.GET("/:serverId/console", middleware.RequiresPermission(pufferpanel.ScopeServerLogs, true), getLogs)
		l.POST("/:serverId/console", middleware.RequiresPermission(pufferpanel.ScopeServerSendCommand, true), postConsole)
		l.OPTIONS("/:serverId/console", response.CreateOptions("GET", "POST"))

		l.GET("/:serverId/stats", middleware.RequiresPermission(pufferpanel.ScopeServerStat, true), getStats)
		l.OPTIONS("/:serverId/stats", response.CreateOptions("GET"))

		l.GET("/:serverId/status", middleware.RequiresPermission(pufferpanel.ScopeServerStatus, true), getStatus)
		l.OPTIONS("/:serverId/status", response.CreateOptions("GET"))

		l.POST("/:serverId/archive/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit, true), archive)
		l.GET("/:serverId/extract/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit, true), extract)

		l.GET("/:serverId/socket", middleware.RequiresPermission(pufferpanel.ScopeServerList, true), cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowCredentials: true,
		}), openSocket)

		l.Handle("CONNECT", "/:serverId/socket", middleware.RequiresPermission(pufferpanel.ScopeServerList, true), func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "false")
		})
		l.OPTIONS("/:serverId/socket", response.CreateOptions("GET", "CONNECT"))

	}

	l.POST("", middleware.RequiresPermission(pufferpanel.ScopeServerCreate, false), createServer)
	l.OPTIONS("", response.CreateOptions("POST"))
}

func getServerFromGin(c *gin.Context) *servers.Server {
	return c.MustGet("program").(*servers.Server)
}

// @Summary Start server
// @Description Start server
// @Success 202 {object} nil
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/start [post]
func startServer(c *gin.Context) {
	server := getServerFromGin(c)
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
				logging.Error.Printf("Error starting server %s: %s", server.Id(), err)
			}
		}()
		c.Status(http.StatusAccepted)
	}
}

// @Summary Stop server
// @Description Stop server
// @Success 202 {object} nil
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/stop [post]
func stopServer(c *gin.Context) {
	server := getServerFromGin(c)

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
		c.Status(http.StatusAccepted)
	}
}

// @Summary Kill server
// @Description Kill server
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/start [post]
func killServer(c *gin.Context) {
	server := getServerFromGin(c)

	err := server.Kill()
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Create server
// @Description Create server
// @Success 200 {object} pufferpanel.ServerIdResponse
// @Param id path string false "Server ID"
// @Param server body pufferpanel.Server true "Server definition"
// @Router /api/servers/{id} [put]
// @Router /api/servers [post]
func createServer(c *gin.Context) {
	serverId := c.Param("serverId")
	if serverId == "" {
		id, err := uuid.NewV4()
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
		serverId = id.String()
	}
	prg, _ := servers.Get(serverId)

	if prg != nil {
		response.HandleError(c, pufferpanel.ErrServerAlreadyExists, http.StatusConflict)
		return
	}

	prg = servers.CreateProgram()
	err := json.NewDecoder(c.Request.Body).Decode(prg)
	if err != nil {
		logging.Error.Printf("Error decoding JSON body: %s", err)
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}
	prg.Identifier = serverId

	err = prg.Requirements.Test(prg.Server)
	if err != nil {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}

	if err := servers.Create(prg); err != nil {
		response.HandleError(c, err, http.StatusInternalServerError)
		_ = servers.Delete(prg.Id())
		return
	}

	c.JSON(http.StatusOK, &pufferpanel.ServerIdResponse{Id: serverId})
}

// @Summary Delete server
// @Description Delete server
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Router /api/servers/{id} [delete]
func deleteServer(c *gin.Context) {
	server := getServerFromGin(c)

	err := servers.Delete(server.Id())
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Install server
// @Description Install server
// @Success 202 {object} nil
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/install [post]
func installServer(c *gin.Context) {
	server := getServerFromGin(c)

	_, wait := c.GetQuery("wait")

	if wait {
		err := server.Install()
		if response.HandleError(c, err, http.StatusInternalServerError) {
		} else {
			c.Status(http.StatusNoContent)
		}
	} else {
		go func(p *servers.Server) {
			_ = p.Install()
		}(server)

		c.Status(http.StatusAccepted)
	}
}

// Not documented in swagger as overridden on frontend
func editServerData(c *gin.Context) {
	server := getServerFromGin(c)

	data := &pufferpanel.ServerData{}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = server.EditData(data.Variables)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Get server tasks
// @Description Get server tasks
// @Success 200 {object} pufferpanel.ServerTasks
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/tasks [get]
func getServerTasks(c *gin.Context) {
	server := getServerFromGin(c)

	result := pufferpanel.ServerTasks{
		Tasks: make(map[string]pufferpanel.ServerTask),
	}

	for k, v := range server.Scheduler.Tasks {
		result.Tasks[k] = pufferpanel.ServerTask{
			Task: pufferpanel.Task{
				Name:         v.Name,
				CronSchedule: v.CronSchedule,
				Description:  v.Description,
			},
			IsRunning: server.Scheduler.IsTaskRunning(k),
		}
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Get server task
// @Description Get server task by id
// @Success 200 {object} pufferpanel.ServerTask
// @Param id path string true "Server ID"
// @Param taskId path string true "Task ID"
// @Router /api/servers/{id}/tasks/{taskId} [get]
func getServerTask(c *gin.Context) {
	server := getServerFromGin(c)

	var result *pufferpanel.ServerTask

	for k, v := range server.Scheduler.Tasks {
		result = &pufferpanel.ServerTask{
			Task:      v,
			IsRunning: server.Scheduler.IsTaskRunning(k),
		}
	}

	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// @Summary Run server task
// @Description Run a specific task
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param taskId path string true "Task ID"
// @Router /api/servers/{id}/tasks/{taskId}/run [post]
func runServerTask(c *gin.Context) {
	server := getServerFromGin(c)

	taskId := c.Param("taskId")

	err := server.Scheduler.RunTask(taskId)
	if errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		c.Status(http.StatusNotFound)
		return
	}
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Edit server task
// @Description Edit server task by id
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param taskId path string true "Task ID"
// @Param task body pufferpanel.Task true "Task definition"
// @Router /api/servers/{id}/tasks/{taskId} [put]
func editServerTask(c *gin.Context) {
	server := getServerFromGin(c)

	taskId := c.Param("taskId")

	var task pufferpanel.Task
	err := c.ShouldBindJSON(&task)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = server.Scheduler.RemoveTask(taskId)
	if errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		err = nil
	}
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = server.Scheduler.AddTask(taskId, task)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Delete server task
// @Description Delete server task by id
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param taskId path string true "Task ID"
// @Router /api/servers/{id}/tasks/{taskId} [delete]
func deleteServerTask(c *gin.Context) {
	server := getServerFromGin(c)

	taskId := c.Param("taskId")

	var task pufferpanel.Task
	err := c.ShouldBindJSON(&task)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = server.Scheduler.RemoveTask(taskId)
	if errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		c.Status(http.StatusNotFound)
		return
	}
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Reload server
// @Description Reload server
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/reload [post]
func reloadServer(c *gin.Context) {
	server := getServerFromGin(c)

	err := servers.Reload(server.Id())
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Get server data
// @Description Get server variables
// @Success 200 {object} pufferpanel.ServerData
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/data [get]
func getServerData(c *gin.Context) {
	server := getServerFromGin(c)

	data := server.GetData()

	var replacement = make(map[string]pufferpanel.Variable)
	for k, v := range data {
		if v.UserEditable {
			replacement[k] = v
		}
	}
	data = replacement

	c.JSON(http.StatusOK, &pufferpanel.ServerData{Variables: data})
}

// @Summary Get server definition
// @Description Get server definition
// @Success 200 {object} pufferpanel.Server
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/definition [get]
func getServerAdmin(c *gin.Context) {
	server := getServerFromGin(c)

	c.JSON(http.StatusOK, &server.Server)
}

// @Summary Edit server definition
// @Description Updates the server with a new definition
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param server body pufferpanel.Server true "New definition"
// @Router /api/servers/{id}/definition [post]
func editServerAdmin(c *gin.Context) {
	prg := getServerFromGin(c)
	server := &prg.Server

	replacement := &pufferpanel.Server{}
	err := c.BindJSON(replacement)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	//backup, just in case we break
	backup := &pufferpanel.Server{}
	backup.CopyFrom(server)

	//copy from request
	server.CopyFrom(replacement)

	err = servers.Save(prg.Id())
	if response.HandleError(c, err, http.StatusInternalServerError) {
		//REVERT!!!!!!!
		server.CopyFrom(backup)
		return
	}

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func getFile(c *gin.Context) {
	server := getServerFromGin(c)

	targetPath := c.Param("filename")

	data, err := server.GetItem(targetPath)
	defer func() {
		if data != nil {
			pufferpanel.Close(data.Contents)
		}
	}()

	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
		} else if err == pufferpanel.ErrIllegalFileAccess {
			response.HandleError(c, err, http.StatusBadRequest)
		} else {
			response.HandleError(c, err, http.StatusInternalServerError)
		}
		return
	}

	if data.FileList != nil {
		c.JSON(http.StatusOK, data.FileList)
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

func putFile(c *gin.Context) {
	server := getServerFromGin(c)

	targetPath := c.Param("filename")

	if targetPath == "" {
		c.AbortWithStatus(http.StatusNotFound)
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

func deleteFile(c *gin.Context) {
	server := getServerFromGin(c)

	targetPath := c.Param("filename")

	err := server.DeleteItem(targetPath)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

func postConsole(c *gin.Context) {
	server := getServerFromGin(c)

	d, _ := io.ReadAll(c.Request.Body)
	cmd := string(d)
	err := server.Execute(cmd)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

func getStats(c *gin.Context) {
	server := getServerFromGin(c)

	results, err := server.GetEnvironment().GetStats()
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.JSON(http.StatusOK, results)
	}
}

func getLogs(c *gin.Context) {
	server := getServerFromGin(c)

	time := c.DefaultQuery("time", "0")

	castedTime, ok := cast.ToInt64E(time)
	if ok != nil {
		response.HandleError(c, pufferpanel.ErrInvalidUnixTime, http.StatusBadRequest)
		return
	}

	console, epoch := server.GetEnvironment().GetConsoleFrom(castedTime)
	msg := ""
	for _, k := range console {
		msg += k
	}

	c.JSON(http.StatusOK, &pufferpanel.ServerLogs{
		Epoch: epoch,
		Logs:  msg,
	})
}

func getStatus(c *gin.Context) {
	server := getServerFromGin(c)

	running, err := server.IsRunning()

	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.JSON(http.StatusOK, &pufferpanel.ServerRunning{Running: running})
	}
}

func archive(c *gin.Context) {
	server := getServerFromGin(c)
	var files []string

	if err := c.BindJSON(&files); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}
	if len(files) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	destination := c.Param("filename")

	err := server.ArchiveItems(files, destination)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

func extract(c *gin.Context) {
	server := getServerFromGin(c)

	targetPath := c.Param("filename")
	destination := c.Query("destination")

	err := server.Extract(targetPath, destination)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

func openSocket(c *gin.Context) {
	server := getServerFromGin(c)

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	internalMap, _ := c.Get("scopes")
	scopes := internalMap.([]pufferpanel.Scope)

	socket := pufferpanel.Create(conn)

	go listenOnSocket(socket, server, scopes)

	if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServerLogs) {
		server.GetEnvironment().AddListener(socket)
	}
}
