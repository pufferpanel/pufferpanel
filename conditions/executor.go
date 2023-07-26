package conditions

import (
	"errors"
	"github.com/google/cel-go/cel"
	"github.com/spf13/cast"
	"runtime"
)

var GlobalConstantValues = map[string]interface{}{
	VariableOs:   runtime.GOOS,
	VariableArch: runtime.GOARCH,
}

func ResolveIf(condition interface{}, data map[string]interface{}, extraCels []cel.EnvOption) (bool, error) {
	if str, ok := condition.(string); condition == nil || (ok && str == "") {
		if success, exists := data[VariableSuccess]; exists {
			return cast.ToBoolE(success)
		}
		return true, nil
	}

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
