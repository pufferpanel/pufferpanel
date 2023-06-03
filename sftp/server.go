/*
 Copyright 2016 Padduck, LLC

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

package sftp

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/sftp"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/oauth2"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path/filepath"
)

var sftpServer net.Listener

var auth pufferpanel.SFTPAuthorization

func Run() {
	err := runServer()
	if err != nil {
		logging.Error.Printf("Error starting SFTP server: %s", err)
	}
}

func SetAuthorization(service pufferpanel.SFTPAuthorization) {
	auth = service
}

func Stop() {
	if sftpServer != nil {
		_ = sftpServer.Close()
	}
}

func runServer() error {
	if auth == nil {
		auth = &oauth2.WebSSHAuthorization{}
	}

	serverConfig := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return auth.Validate(c.User(), string(pass))
		},
	}

	serverKeyFile := config.SftpKey.Value()

	_, e := os.Stat(serverKeyFile)

	if e != nil && os.IsNotExist(e) {
		logging.Info.Printf("Generating new key")
		var key ed25519.PrivateKey
		_, key, e = ed25519.GenerateKey(nil)
		if e != nil {
			return e
		}

		data, e := x509.MarshalPKCS8PrivateKey(key)
		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: data,
		}
		if e != nil {
			return e
		}

		e = os.WriteFile(serverKeyFile, pem.EncodeToMemory(block), 0700)
		if e != nil {
			return e
		}
	} else if e != nil {
		return e
	}

	logging.Info.Printf("Loading existing key")
	var data []byte
	data, e = os.ReadFile(serverKeyFile)
	if e != nil {
		return e
	}

	hkey, e := ssh.ParsePrivateKey(data)

	if e != nil {
		return e
	}

	serverConfig.AddHostKey(hkey)

	bind := config.SftpHost.Value()

	sftpServer, e = net.Listen("tcp", bind)
	if e != nil {
		return e
	}
	logging.Info.Printf("Started SFTP Server on %s", bind)

	go func() {
		for {
			conn, _ := sftpServer.Accept()
			if conn != nil {
				go HandleConn(conn, serverConfig)
			}
		}
	}()

	return nil
}

func HandleConn(conn net.Conn, serverConfig *ssh.ServerConfig) {
	defer pufferpanel.Close(conn)
	defer pufferpanel.Recover()
	logging.Info.Printf("SFTP connection from %s", conn.RemoteAddr().String())
	e := handleConn(conn, serverConfig)
	if e != nil {
		if e.Error() != "EOF" {
			logging.Error.Printf("sftpd connection error: %s", e)
		}
	}
}
func handleConn(conn net.Conn, serverConfig *ssh.ServerConfig) error {
	sc, chans, reqs, e := ssh.NewServerConn(conn, serverConfig)
	defer pufferpanel.Close(sc)
	if e != nil {
		return e
	}

	// The incoming Request channel must be serviced.
	go PrintDiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of an SFTP session, this is "subsystem"
		// with a payload string of "<length=4>sftp"
		if newChannel.ChannelType() != "session" {
			err := newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			if err != nil {
				return err
			}
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			return err
		}

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "subsystem" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				switch req.Type {
				case "subsystem":
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				_ = req.Reply(ok, nil)
			}
		}(requests)

		fs := CreateRequestPrefix(filepath.Join(config.ServersFolder.Value(), sc.Permissions.Extensions["server_id"]))

		server := sftp.NewRequestServer(channel, fs)

		if err := server.Serve(); err != nil {
			return err
		}
	}
	return nil
}

func PrintDiscardRequests(in <-chan *ssh.Request) {
	for req := range in {
		if req.WantReply {
			_ = req.Reply(false, nil)
		}
	}
}
