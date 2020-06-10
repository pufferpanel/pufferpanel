package auth

import (
	"bytes"
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/services"
)

func GetToken(c *gin.Context) {
	key := services.GetPublicKey()

	pubKeyEncoded, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		logging.Error().Printf("Internal error on token service: %s", err)
		return
	}

	var buffer bytes.Buffer

	err = pem.Encode(&buffer, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyEncoded})
	if err != nil {
		logging.Error().Printf("Internal error on token service: %s", err)
		return
	}

	hasher := md5.New()
	hasher.Write(buffer.Bytes())
	hash := hex.EncodeToString(hasher.Sum(nil))
	c.Header("Cache-Control", "public")
	c.Header("ETag", hash)

	if c.GetHeader("ETag") == hash {
		c.Status(304)
		return
	}

	c.Data(200, "application/x-pem-file", buffer.Bytes())
}
