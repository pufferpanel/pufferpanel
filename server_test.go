package pufferpanel

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseRequirementRow(t *testing.T) {
	tests := []struct {
		str  string
		want []string
	}{
		{str: "", want: []string{}},
		{str: "linux", want: []string{"linux"}},
		{str: "  linux   ", want: []string{"linux"}},
		{str: "linux||windows", want: []string{"linux", "windows"}},
		{str: " linux    || windows  ", want: []string{"linux", "windows"}},
	}

	for _, tt := range tests {
		got := parseRequirementRow(tt.str)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseRequirementRow(%#v) = %v, want %v", tt.str, got, tt.want)
		}
	}
}

func TestVariable_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		args         []byte
		expectedType string
	}{
		{
			name:         "DefaultString",
			args:         []byte(`{ "type": "", "value": "0.0.0.0", "display": "IP", "desc": "What IP to bind the server to", "required": true }`),
			expectedType: "string",
		},
		{
			name:         "TypeInt",
			args:         []byte(`{ "type": "integer", "value": 12345, "display": "Port", "desc": "Port", "required": true }`),
			expectedType: "integer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Variable
			err := json.Unmarshal(tt.args, &v)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedType, v.Type.Type)
		})
	}
}
