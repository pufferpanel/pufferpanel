package models

import (
	"github.com/gobuffalo/uuid"
	"time"
)

type Scope struct {
	ID        int       `json:"-" db:"id"`
	ClientId  uuid.UUID `json:"-" db:"client_id" rw:"r"`
	ServerId  uuid.UUID `json:"server_id" db:"server_id" rw:"r"`
	value     string    `json:"value" db:"value" rw:"r"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Scopes []Scope

func (s *Scope) GetValue() string {
	return s.ServerId.String() + ":" + s.value
}

func (s *Scope) GetRawValue() string {
	return s.value
}
