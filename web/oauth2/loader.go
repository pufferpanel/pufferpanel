package oauth2

import "github.com/gin-gonic/gin"

func Register(rg *gin.RouterGroup) {
	registerTokens(rg)
}