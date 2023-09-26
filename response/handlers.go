package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func NotImplemented(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotImplemented, pufferpanel.ErrNotImplemented)
}

func CreateOptions(options ...string) gin.HandlerFunc {
	replacement := make([]string, len(options)+1)

	copy(replacement, options)

	replacement[len(options)] = http.MethodOptions
	res := strings.Join(replacement, ",")

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", res)
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", res)
		c.Header("Content-Type", binding.MIMEJSON)
		c.AbortWithStatus(http.StatusOK)
	}
}

func HandleError(c *gin.Context, err error, statusCode int) bool {
	if err != nil {
		logging.Error.Printf("%s", err.Error())

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatusJSON(statusCode, &pufferpanel.ErrorResponse{Error: pufferpanel.FromError(err)})
		}

		return true
	}

	return false
}
