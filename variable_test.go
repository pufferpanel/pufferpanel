package pufferpanel_test

import (
	"testing"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestVariable_SetValue(t *testing.T) {
	tests := []struct {
		name        string
		valType     string
		options     []any
		val         any
		overrideVal any
		wantErr     bool
	}{
		{
			valType: "string",
			val:     "",
		},
		{
			valType: "string",
			val:     "test",
		},
		{
			valType:     "string",
			val:         123,
			overrideVal: "123",
		},
		{
			valType:     "string",
			val:         false,
			overrideVal: "false",
		},
		{
			valType: "string",
			val:     "12345",
		},
		{
			valType: "integer",
			val:     "bad data",
			wantErr: true,
		},
		{
			valType:     "integer",
			val:         "",
			overrideVal: 0,
		},
		{
			valType: "integer",
			val:     0,
		},
		{
			valType: "integer",
			val:     1,
		},
		{
			valType: "integer",
			val:     -1,
		},
		{
			valType: "integer",
			val:     60030201201,
		},
		{
			valType: "integer",
			val:     0,
		},
		{
			valType:     "boolean",
			val:         0,
			overrideVal: false,
		},
		{
			valType:     "boolean",
			val:         1,
			overrideVal: true,
		},
		{
			valType:     "boolean",
			val:         "0",
			overrideVal: false,
		},
		{
			valType:     "boolean",
			val:         "1",
			overrideVal: true,
		},
		{
			valType:     "boolean",
			val:         "false",
			overrideVal: false,
		},
		{
			valType:     "boolean",
			val:         "true",
			overrideVal: true,
		},
		{
			valType: "boolean",
			val:     false,
		},
		{
			valType: "boolean",
			val:     true,
		},
		{
			valType: "boolean",
			val:     "invalid",
			wantErr: true,
		},
		{
			valType: "badtype",
			wantErr: true,
		},
		{
			valType: "option",
			options: []any{
				"one", "two", "three",
			},
			val: "one",
		},
		{
			valType: "option",
			options: []any{
				"one", "two", "three",
			},
			val: "two",
		},
		{
			valType: "option",
			options: []any{
				1, 2, 3,
			},
			val: 2,
		},
		{
			valType: "option",
			options: []any{
				"one", "two", "three",
			},
			val:     "four",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.valType+"_"+cast.ToString(tt.val), func(t *testing.T) {
			v := &pufferpanel.Variable{
				Type: pufferpanel.Type{
					Type: tt.valType,
				},
				Value: "START_VALUE",
			}
			if len(tt.options) > 0 {
				for _, o := range tt.options {
					v.Options = append(v.Options, pufferpanel.VariableOption{Value: o})
				}
			}
			gotErr := v.SetValue(tt.val)
			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Equal(t, v.Value, "START_VALUE")
				return
			} else if !assert.NoError(t, gotErr) {
				return
			}

			if tt.overrideVal == nil {
				assert.Equal(t, tt.val, v.Value)
			} else {
				assert.Equal(t, tt.overrideVal, v.Value)
			}
		})
	}
}
