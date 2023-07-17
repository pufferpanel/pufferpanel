package conditions

import (
	"github.com/google/cel-go/cel"
	"github.com/spf13/cast"
	"runtime"
)

var GlobalConstantValues = map[string]interface{}{
	Os:   runtime.GOOS,
	Arch: runtime.GOARCH,
}

func ResolveIf(condition interface{}, data map[string]interface{}, extraCels []cel.EnvOption) (bool, error) {
	cond, err := cast.ToStringE(condition)
	if err != nil {
		return false, err
	}

	celVars := extraCels

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
