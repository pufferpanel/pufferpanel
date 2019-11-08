package services

import "golang.org/x/crypto/ssh"

type DatabaseSFTPAuthorization struct {
}

func (s *DatabaseSFTPAuthorization) Validate(username, password string) (perms *ssh.Permissions, err error) {
	return
}
