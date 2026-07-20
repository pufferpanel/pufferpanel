package pufferpanel

import (
	"encoding/json"

	"github.com/pufferpanel/pufferpanel/v3/utils"
)

type StdinConsoleConfiguration struct {
	Type     string `json:"type,omitempty"`
	IP       string `json:"ip,omitempty"`
	Port     string `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
} //@name StdinConsoleConfiguration

type stdinConfigAlias StdinConsoleConfiguration

func (v *StdinConsoleConfiguration) Replace(variables map[string]interface{}) StdinConsoleConfiguration {
	return StdinConsoleConfiguration{
		Type:     v.Type,
		IP:       utils.ReplaceTokens(v.IP, variables, utils.PlainReplace),
		Port:     utils.ReplaceTokens(v.Port, variables, utils.PlainReplace),
		Password: utils.ReplaceTokens(v.Password, variables, utils.PlainReplace),
	}
}

func (v *StdinConsoleConfiguration) UnmarshalJSON(data []byte) error {
	aux := stdinConfigAlias{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Type == "" {
		aux.Type = "stdin"
	}
	*v = StdinConsoleConfiguration(aux)
	return nil
}
