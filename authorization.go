package pufferpanel

import (
	"golang.org/x/crypto/ssh"
)

type SFTPAuthorization interface {
	Validate(username, password string) (perms *ssh.Permissions, err error)
}
