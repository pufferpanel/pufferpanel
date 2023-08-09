package models

type Setting struct {
	Value interface{} `json:"value"`
} //@name Setting

type ChangeMultipleSettings map[string]interface{} //@name ChangeMultipleSettings
