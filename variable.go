package pufferpanel

import (
	"encoding/json"
	"errors"

	"github.com/spf13/cast"
)

type Variable struct {
	Type
	Value        any              `json:"value"`
	Display      string           `json:"display,omitempty"`
	Description  string           `json:"desc,omitempty"`
	Required     bool             `json:"required"`
	Internal     bool             `json:"internal,omitempty"`
	UserEditable bool             `json:"userEdit"`
	Options      []VariableOption `json:"options,omitempty"`
} //@name Variable
type variableAlias Variable

type VariableOption struct {
	Value   any    `json:"value"`
	Display string `json:"display"`
} //@name VariableOption

func (v *Variable) SetValue(val any) (err error) {
	var tmp any
	//convert variable to correct typing
	switch v.Type.Type {
	case "string":
		{
			tmp, err = cast.ToStringE(val)
			if err == nil {
				v.Value = tmp
			}
		}
	case "integer":
		{
			if tmp, err = cast.ToStringE(val); err == nil && tmp == "" {
				v.Value = 0
			} else {
				tmp, err = cast.ToIntE(val)
				if err == nil {
					v.Value = tmp
				}
			}
		}
	case "boolean":
		{
			tmp, err = cast.ToBoolE(val)
			if err == nil {
				v.Value = tmp
			}
		}
	case "option":
		{
			updated := false
			for _, k := range v.Options {
				if k.Value == val {
					v.Value = k.Value
					updated = true
					break
				}
			}
			if !updated {
				err = errors.New("unknown value")
			}
		}
	default:
		err = errors.New("unknown variable type")
	}

	return err
}

func (v *Variable) UnmarshalJSON(data []byte) (err error) {
	aux := variableAlias{}
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}
	if aux.Type.Type == "" {
		aux.Type = Type{Type: "string"}
	}

	//default any null value to empty string
	if aux.Value == nil {
		aux.Value = ""
	}

	*v = Variable(aux)
	err = v.SetValue(aux.Value)
	return
}
