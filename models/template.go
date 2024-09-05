package models

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
	"gorm.io/gorm"
	"strings"
)

type Template struct {
	pufferpanel.Server `gorm:"-"`

	Name     string `gorm:"column:name;size:100;primaryKey" json:"name"`
	RawValue string `gorm:"column:raw_value;not null;size:4000" json:"-"`

	Readme string `gorm:"column:readme;size:4000" json:"readme,omitempty"`
} //@name Template

func (t *Template) AfterFind(*gorm.DB) error {
	err := json.NewDecoder(strings.NewReader(t.RawValue)).Decode(&t.Server)
	if err != nil {
		return err
	}
	t.RawValue = ""
	return nil
}

func (t *Template) BeforeSave(*gorm.DB) error {
	data, err := json.Marshal(&t.Server)
	if err != nil {
		return err
	}
	t.RawValue = string(data)
	return nil
}
