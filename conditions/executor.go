package conditions

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v3/thirdparty/cel-go/cel"
	"runtime"
)

var GlobalConstantValues = map[string]interface{}{
	VariableOs:   runtime.GOOS,
	VariableArch: runtime.GOARCH,
}

func ResolveIf(condition interface{}, data map[string]interface{}, extraCels []cel.EnvOption) (bool, error) {
	var cond string
	var ok bool
	if cond, ok = condition.(string); !ok {
		return false, errors.New("unknown type for condition")
	}

	celVars := extraCels
	if celVars == nil {
		celVars = make([]cel.EnvOption, 0)
	}

	inputData := map[string]interface{}{}

	if data != nil {
		for k, v := range data {
			celVars = append(celVars, cel.Variable(k, cel.AnyType))
			inputData[k] = v
		}
	}

	for k, v := range GlobalConstantValues {
		inputData[k] = v
	}

	celEnv, err := cel.NewEnv(celVars...)
	if err != nil {
		return false, err
	}

	ast, issues := celEnv.Compile(cond)
	if issues != nil && issues.Err() != nil {
		return false, issues.Err()
	}

	prg, err := celEnv.Program(ast)
	if err != nil {
		return false, err
	}

	out, _, err := prg.Eval(inputData)
	if err != nil {
		return false, err
	}
	return out.Value().(bool), nil
}
