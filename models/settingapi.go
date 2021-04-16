package models

type ChangeSetting struct {
	Value interface{} `json:"value"`
}

type SettingResponse struct {
	Value string `json:"value"`
}
