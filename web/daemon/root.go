package daemon

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/servers"
	"net/http"
	"runtime"
	"time"
)

func RegisterDaemonRoutes(e *gin.RouterGroup) {
	e.GET("", getStatusGET)
	e.HEAD("", getStatusHEAD)
	e.Handle("OPTIONS", "", response.CreateOptions("GET", "HEAD"))

	e.GET("features", getFeatures)
	e.Handle("OPTIONS", "features", response.CreateOptions("GET"))

	RegisterServerRoutes(e)
}

// @Summary Check daemon status
// @Description Check to see if the daemon is online or not
// @Success 200 {object} pufferpanel.DaemonRunning
// @Router /daemon [get]
// @Security OAuth2Application[none]
func getStatusGET(c *gin.Context) {
	c.JSON(http.StatusOK, &pufferpanel.DaemonRunning{Message: "daemon is running"})
}

// @Summary Check daemon status
// @Description Check to see if the daemon is online or not
// @Success 204 {object} nil
// @Router /daemon [head]
// @Security OAuth2Application[none]
func getStatusHEAD(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// @Summary Get features of the node
// @Description Gets the features that the node supports, like it's OS and environments
// @Success 200 {object} Features
// @Router /daemon/features [get]
// @Security OAuth2Application[none]
func getFeatures(c *gin.Context) {
	features := make([]string, 0)

	envs := servers.GetSupportedEnvironments()

	if testDocker() {
		features = append(features, "docker")
	}

	c.JSON(http.StatusOK, Features{Features: features, Environments: envs, OS: runtime.GOOS, Arch: runtime.GOARCH})
}

func testDocker() bool {
	d, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = d.Ping(ctx)
	return err == nil
}

type Features struct {
	Features     []string `json:"features"`
	Environments []string `json:"environments"`
	OS           string   `json:"os"`
	Arch         string   `json:"arch"`
} //@name Features
