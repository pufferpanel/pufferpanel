/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package pufferpanel

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
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
	defer Close(f)

	var buf bytes.Buffer

	_, _ = io.Copy(&buf, f)

	block, _ := pem.Decode(buf.Bytes())
	if block == nil {
		return nil, ErrKeyNotPEM
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrKeyNotECDSA
	}

	SetPublicKey(publicKey)
	return publicKey, nil
}
