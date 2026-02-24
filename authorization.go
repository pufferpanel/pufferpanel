package pufferpanel

import (
	"golang.org/x/crypto/ssh"
)

type SFTPAuthorization interface {
	ValidatePassword(username, password string) (perms *ssh.Permissions, err error)

	ValidatePublicKey(username string, key ssh.PublicKey) (perms *ssh.Permissions, err error)
}
