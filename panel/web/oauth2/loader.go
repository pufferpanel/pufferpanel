package oauth2

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	registerTokens(rg)
}