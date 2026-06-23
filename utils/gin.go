package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetQueryBool(c *gin.Context, key string) bool {
	var val string
	exists := false

	for k, v := range c.Request.URL.Query() {
		if strings.EqualFold(k, key) {
			exists = true
			if len(v) > 0 {
				val = v[0]
			}
		}
	}

	if !exists {
		return false
	}

	if val == "" {
		return true
	}

	return cast.ToBool(val)
}
