package daemon

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"sync"
)

type SFTPAuthorization interface {
	Validate(username, password string) (perms *ssh.Permissions, err error)
}

var publicKey *ecdsa.PublicKey

var atLocker = &sync.RWMutex{}

func SetPublicKey(key *ecdsa.PublicKey) {
	atLocker.Lock()
	defer atLocker.Unlock()
	publicKey = key
}

func GetPublicKey() *ecdsa.PublicKey {
	atLocker.RLock()
	defer atLocker.RUnlock()
	return publicKey
}

func LoadPublicKey() (*ecdsa.PublicKey, error) {
	publicKey := GetPublicKey()
	if publicKey != nil {
		return publicKey, nil
	}

	f, err := os.OpenFile(viper.GetString("daemon.auth.publicKey"), os.O_RDONLY, 660)
	defer pufferpanel.Close(f)

	var buf bytes.Buffer

	_, _ = io.Copy(&buf, f)

	block, _ := pem.Decode(buf.Bytes())
	if block == nil {
		return nil, pufferpanel.ErrKeyNotPEM
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, pufferpanel.ErrKeyNotECDSA
	}

	SetPublicKey(publicKey)
	return publicKey, nil
}
