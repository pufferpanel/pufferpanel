package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"net/http"
)

func DecryptPayload(c *gin.Context) {
	//original body is a JWT-encoded string
	//convert it to the original JSON message
	var err error

	reader := c.Request.Body
	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

}
