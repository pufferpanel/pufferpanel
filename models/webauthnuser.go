package models

import (
	"encoding/binary"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthnUser struct {
	User *User
	Credentials []webauthn.Credential
}

func (wau *WebAuthnUser) WebAuthnID() []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(wau.User.ID))
	return bytes
}

func (wau *WebAuthnUser) WebAuthnName() string {
	return wau.User.Email
}

func (wau *WebAuthnUser) WebAuthnDisplayName() string {
	return wau.User.Username
}

func (wau *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return wau.Credentials
}
