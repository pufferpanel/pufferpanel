// +build !nodocs

package swagger

import (
	"github.com/gin-gonic/gin"
	_ "github.com/pufferpanel/pufferpanel/v2/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Load(e *gin.Engine) {
	e.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
