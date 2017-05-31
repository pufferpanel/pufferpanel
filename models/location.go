package models

import (
	"encoding/json"
	"time"
)

type Location struct {
	ID          int       `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
}

// String is not required by pop and may be deleted
func (l Location) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}
