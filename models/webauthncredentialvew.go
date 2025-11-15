package models

import (
	webauthnProto "github.com/go-webauthn/webauthn/protocol"
)

type WebauthnCredentialView struct {
	ID webauthnProto.URLEncodedBase64 `json:"id" swaggertype:"string"`
	Name string `json:"name"`
	Descriptor webauthnProto.CredentialDescriptor `json:"descriptor"`
}