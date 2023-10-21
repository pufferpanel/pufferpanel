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
	l := e.Group("/server", middleware.IsPanelCaller)
	{
		l.PUT("/:serverId", createServer)
		l.DELETE("/:serverId", middleware.ResolveServerNode, deleteServer)
		l.OPTIONS("/:serverId", response.CreateOptions("PUT", "DELETE", "GET"))

		l.GET("/:serverId/definition", middleware.ResolveServerNode, getServerAdmin)
		l.PUT("/:serverId/definition", middleware.ResolveServerNode, editServerAdmin)
		l.OPTIONS("/:serverId/definition", response.CreateOptions("PUT", "DELETE", "GET"))

		l.GET("/:serverId/data", middleware.ResolveServerNode, getServerData)
		l.POST("/:serverId/data", middleware.ResolveServerNode, editServerData)
		l.OPTIONS("/:serverId/data", response.CreateOptions("GET", "POST"))

		l.GET("/:serverId/tasks", middleware.ResolveServerNode, getServerTasks)
		l.OPTIONS("/:serverId/tasks", response.CreateOptions("GET"))

		l.GET("/:serverId/tasks/:taskId", middleware.ResolveServerNode, getServerTask)
		l.PUT("/:serverId/tasks/:taskId", middleware.ResolveServerNode, editServerTask)
		l.DELETE("/:serverId/tasks/:taskId", middleware.ResolveServerNode, deleteServerTask)
		l.OPTIONS("/:serverId/tasks/:taskId", response.CreateOptions("GET", "PUT", "DELETE"))

		l.POST("/:serverId/tasks/:taskId/run", middleware.ResolveServerNode, runServerTask)
		l.OPTIONS("/:serverId/tasks/:taskId/run", response.CreateOptions("POST"))

		l.POST("/:serverId/reload", middleware.ResolveServerNode, reloadServer)
		l.OPTIONS("/:serverId/reload", response.CreateOptions("POST"))

		l.POST("/:serverId/start", middleware.ResolveServerNode, startServer)
		l.OPTIONS("/:serverId/start", response.CreateOptions("POST"))

		l.POST("/:serverId/stop", middleware.ResolveServerNode, stopServer)
		l.OPTIONS("/:serverId/stop", response.CreateOptions("POST"))

		l.POST("/:serverId/kill", middleware.ResolveServerNode, killServer)
		l.OPTIONS("/:serverId/kill", response.CreateOptions("POST"))

		l.POST("/:serverId/install", middleware.ResolveServerNode, installServer)
		l.OPTIONS("/:serverId/install", response.CreateOptions("POST"))

		l.GET("/:serverId/file/*filename", middleware.ResolveServerNode, getFile)
		l.PUT("/:serverId/file/*filename", middleware.ResolveServerNode, putFile)
		l.DELETE("/:serverId/file/*filename", middleware.ResolveServerNode, deleteFile)
		l.POST("/:serverId/file/*filename", middleware.ResolveServerNode, response.NotImplemented)
		l.OPTIONS("/:serverId/file/*filename", response.CreateOptions("GET", "PUT", "DELETE", "POST"))

		l.GET("/:serverId/console", middleware.ResolveServerNode, getLogs)
		l.POST("/:serverId/console", middleware.ResolveServerNode, postConsole)
		l.OPTIONS("/:serverId/console", response.CreateOptions("GET", "POST"))

		l.GET("/:serverId/flags", middleware.ResolveServerNode, getFlags)
		l.POST("/:serverId/flags", middleware.ResolveServerNode, setFlags)
		l.OPTIONS("/:serverId/flags", response.CreateOptions("GET", "POST"))

		l.GET("/:serverId/stats", middleware.ResolveServerNode, getStats)
		l.OPTIONS("/:serverId/stats", response.CreateOptions("GET"))

		l.GET("/:serverId/status", middleware.ResolveServerNode, getStatus)
		l.OPTIONS("/:serverId/status", response.CreateOptions("GET"))

		l.POST("/:serverId/archive/*filename", middleware.ResolveServerNode, archive)
		l.GET("/:serverId/extract/*filename", middleware.ResolveServerNode, extract)

		p := l.Group("/:serverId/socket")
		{
			p.GET("", middleware.ResolveServerNode, cors.New(cors.Config{
				AllowAllOrigins:  true,
				AllowCredentials: true,
			}), openSocket)
			p.Handle("CONNECT", "", func(c *gin.Context) {
				c.Header("Access-Control-Allow-Origin", "*")
				c.Header("Access-Control-Allow-Credentials", "false")
			})
			p.OPTIONS("", response.CreateOptions("GET", "CONNECT"))
		}
	}
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
// @Security OAuth2Application[server.start]
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
// @Security OAuth2Application[server.stop]
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
// @Router /api/servers/{id}/kill [post]
// @Security OAuth2Application[server.kill]
func killServer(c *gin.Context) {
	server := getServerFromGin(c)

	err := server.Kill()
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// Already declared in panel routing
func createServer(c *gin.Context) {
	serverId := c.Param("serverId")
	if serverId == "" {
		id, err := uuid.NewV4()
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
		serverId = id.String()
	}
	prg := servers.GetFromCache(serverId)

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

	if prg, err = servers.Create(prg); err != nil {
		response.HandleError(c, err, http.StatusInternalServerError)
		_ = servers.Delete(prg.Id())
		return
	}

	c.JSON(http.StatusOK, &pufferpanel.ServerIdResponse{Id: serverId})
}

// Already declared in panel routing
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
// @Security OAuth2Application[server.install]
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

	var data map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = server.EditData(data)
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
// @Security OAuth2Application[server.tasks.view]
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
// @Security OAuth2Application[server.tasks.view]
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
// @Security OAuth2Application[server.tasks.run]
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
// @Security OAuth2Application[server.tasks.edit]
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
// @Security OAuth2Application[server.tasks.delete]
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
// @Security OAuth2Application[server.reload]
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
// @Security OAuth2Application[server.data.view]
func getServerData(c *gin.Context) {
	server := getServerFromGin(c)

	data := server.GetData()

	var replacement = make(map[string]pufferpanel.Variable)
	for k, v := range data {
		if v.UserEditable {
			replacement[k] = v
		}
	}

	c.JSON(http.StatusOK, &pufferpanel.ServerData{Variables: replacement})
}

// @Summary Get server definition
// @Description Get server definition
// @Success 200 {object} pufferpanel.Server
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/definition [get]
// @Security OAuth2Application[server.definition.view]
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
// @Security OAuth2Application[server.definition.edit]
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

// @Summary Get file/folder
// @Description Gets a specific file or a list of files in a folder. This will either return
// @Description a) A raw file if the path points to a valid file
// @Description or b) An array of files for the folder contents
// @Success 200 {object} nil
// @Param id path string true "Server ID"
// @Param filepath path string true "File path"
// @Router /api/servers/{id}/file/{filepath} [get]
// @Security OAuth2Application[server.files.view]
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

// @Summary Edit file
// @Description Adds or edit a file, replacing the contents with the provided body
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param filepath path string true "File path"
// @Param file formData file true "File contents"
// @Accept multipart/form-data
// @Router /api/servers/{id}/file/{filepath} [put]
// @Security OAuth2Application[server.files.edit]
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
		return
	}

	_, err = io.Copy(file, sourceFile)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Delete file
// @Description Deletes a file or folder.
// @Description WARNING: This is a recursive operation, specifying a folder will delete all children
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param filepath path string true "File path"
// @Router /api/servers/{id}/file/{filepath} [delete]
// @Security OAuth2Application[server.files.edit]
func deleteFile(c *gin.Context) {
	server := getServerFromGin(c)

	targetPath := c.Param("filename")

	err := server.DeleteItem(targetPath)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Send command
// @Description Sends a command to the server
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param command body string true "Command"
// @Router /api/servers/{id}/console [post]
// @Security OAuth2Application[server.console.send]
func postConsole(c *gin.Context) {
	server := getServerFromGin(c)

	d, err := io.ReadAll(c.Request.Body)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	cmd, err := cast.ToStringE(d)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = server.Execute(cmd)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	c.Status(http.StatusNoContent)

}

// @Summary Get stats
// @Description Gets the CPU and memory usage of the server
// @Success 200 {object} pufferpanel.ServerStats
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/stats [get]
// @Security OAuth2Application[server.stats]
func getStats(c *gin.Context) {
	server := getServerFromGin(c)

	results, err := server.GetEnvironment().GetStats()
	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.JSON(http.StatusOK, results)
	}
}

// @Summary Get logs
// @Description Get the console logs for the server
// @Success 200 {object} pufferpanel.ServerLogs
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/console [get]
// @Security OAuth2Application[server.console]
func getLogs(c *gin.Context) {
	server := getServerFromGin(c)

	time := c.DefaultQuery("time", "0")

	castedTime, ok := cast.ToInt64E(time)
	if ok != nil {
		response.HandleError(c, pufferpanel.ErrInvalidUnixTime, http.StatusBadRequest)
		return
	}

	console, epoch := server.GetEnvironment().GetConsoleFrom(castedTime)

	c.JSON(http.StatusOK, &pufferpanel.ServerLogs{
		Epoch: epoch,
		Logs:  console,
	})
}

// @Summary Get status
// @Description Get the server's status (is it running)
// @Success 200 {object} pufferpanel.ServerRunning
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/status [get]
// @Security OAuth2Application[server.status]
func getStatus(c *gin.Context) {
	server := getServerFromGin(c)

	installing := server.GetEnvironment().IsInstalling()

	if installing {
		c.JSON(http.StatusOK, &pufferpanel.ServerRunning{Installing: installing})
		return
	}

	running, err := server.IsRunning()

	if response.HandleError(c, err, http.StatusInternalServerError) {
	} else {
		c.JSON(http.StatusOK, &pufferpanel.ServerRunning{Running: running})
	}
}

// @Summary Create archive
// @Description Creates an archive of files or folders
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param files body []string true "Files to archive"
// @Param filename path string true "Archive name"
// @Router /api/servers/{id}/archive/{filename} [post]
// @Security OAuth2Application[server.files.edit]
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

// @Summary Extract archive
// @Description Extracts an archive to the server
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param filename path string true "Target file to extract"
// @Param destination query string true "Path to place files"
// @Router /api/servers/{id}/extract/{filename} [post]
// @Security OAuth2Application[server.files.edit]
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

// @Summary Get flags
// @Description Get the management flags for a server
// @Success 200 {object} pufferpanel.ServerFlags
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/flags [get]
// @Security OAuth2Application[server.flags.view]
func getFlags(c *gin.Context) {
	server := getServerFromGin(c)

	c.JSON(http.StatusOK, &pufferpanel.ServerFlags{
		AutoStart:             &server.Execution.AutoStart,
		AutoRestartOnCrash:    &server.Execution.AutoRestartFromCrash,
		AutoRestartOnGraceful: &server.Execution.AutoRestartFromGraceful,
	})
}

// @Summary Set flags
// @Description Sets management flags for a server
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param flags body pufferpanel.ServerFlags true "Flags to change"
// @Router /api/servers/{id}/flags [post]
// @Security OAuth2Application[server.flags.edit]
func setFlags(c *gin.Context) {
	server := getServerFromGin(c)

	var req pufferpanel.ServerFlags
	err := c.BindJSON(&req)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}
	if req.AutoRestartOnCrash != nil {
		server.Execution.AutoRestartFromCrash = *req.AutoRestartOnCrash
	}
	if req.AutoRestartOnGraceful != nil {
		server.Execution.AutoRestartFromGraceful = *req.AutoRestartOnGraceful
	}
	if req.AutoStart != nil {
		server.Execution.AutoStart = *req.AutoStart
	}
	err = server.Save()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	c.Status(http.StatusNoContent)
}

func openSocket(c *gin.Context) {
	server := getServerFromGin(c)

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	socket := pufferpanel.Create(conn)

	if _, exists := c.GetQuery("console"); exists {
		server.GetEnvironment().AddConsoleListener(socket)
	}

	if _, exists := c.GetQuery("stats"); exists {
		server.GetEnvironment().AddStatsListener(socket)
	}

	if _, exists := c.GetQuery("status"); exists {
		server.GetEnvironment().AddStatusListener(socket)
	}
}
