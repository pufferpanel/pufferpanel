package conditions

import (
	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	type tests struct {
		name       string
		expression string
		result     bool
		variables  map[string]interface{}
	}
	ts := []tests{
		{
			name:       "defined item",
			expression: `has(vars.item) && vars.item == "exists"`,
			result:     true,
			variables:  map[string]interface{}{"item": "exists"},
		},
		{
			name:       "undefined item",
			expression: `has(vars.items) && vars.items == "exists"`,
			result:     false,
			variables:  map[string]interface{}{"item": "exists"},
		},
		{
			name:       "defined item as map",
			expression: `has(vars.item) && vars["item"] == "exists"`,
			result:     true,
			variables:  map[string]interface{}{"item": "exists"},
		},
		{
			name:       "undefined item as map",
			expression: `has(vars.items) && vars["items"] == "exists"`,
			result:     false,
			variables:  map[string]interface{}{"item": "exists"},
		},
	}

	for _, tt := range ts {
		t.Run(tt.name, func(t *testing.T) {
			celVars := make([]cel.EnvOption, 0)

			inputData := map[string]interface{}{}

			celVars = append(celVars, cel.Variable("vars", cel.MapType(cel.StringType, cel.DynType)))

			for k, v := range tt.variables {
				inputData[k] = v
			}

			celEnv, err := cel.NewEnv(celVars...)
			if !assert.NoError(t, err) {
				return
			}

			ast, issues := celEnv.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				assert.NoError(t, issues.Err())
				return
			}

			prg, err := celEnv.Program(ast)
			if !assert.NoError(t, err) {
				return
			}

			out, _, err := prg.Eval(map[string]any{
				"vars": inputData,
			})
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.result, out.Value())
		})
	}
}
