// +build !nodocs

package swagger

import (
	"github.com/gin-gonic/gin"
	_ "github.com/pufferpanel/pufferpanel/v2/panel/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Load(e *gin.RouterGroup) {
	e.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
