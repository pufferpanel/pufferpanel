package conditions

import (
	"fmt"
	"github.com/google/cel-go/cel"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var GlobalConstantValues = map[string]interface{}{
	VariableOs:   runtime.GOOS,
	VariableArch: runtime.GOARCH,
}

func ResolveIf(condition string, data map[string]interface{}, extraCels []cel.EnvOption) (bool, error) {
	if condition == "" {
		return true, nil
	}
	return Run[bool](condition, data, extraCels)
}

func Run[T interface{}](statement string, data map[string]interface{}, extras []cel.EnvOption) (T, error) {
	var res T

	celVars := extras
	if celVars == nil {
		celVars = make([]cel.EnvOption, 0)
	}

	inputData := map[string]interface{}{}

	for k, v := range data {
		celVars = append(celVars, cel.Variable(k, cel.AnyType))
		inputData[k] = v
	}

	for k, v := range GlobalConstantValues {
		celVars = append(celVars, cel.Variable(k, cel.AnyType))
		inputData[k] = v
	}

	celEnv, err := cel.NewEnv(celVars...)
	if err != nil {
		return res, err
	}

	ast, issues := celEnv.Compile(statement)
	if issues != nil && issues.Err() != nil {
		return res, issues.Err()
	}

	prg, err := celEnv.Program(ast)
	if err != nil {
		return res, err
	}

	out, _, err := prg.Eval(inputData)
	if err != nil {
		return res, err
	}
	if cast, ok := out.Value().(T); ok {
		return cast, nil
	} else {
		return res, fmt.Errorf("invalid return type, expected %s, got %s", reflect.TypeOf(res), reflect.TypeOf(cast))
	}
}

var conditionalStatementRegex = regexp.MustCompile("{{.*?}}")

func ReplaceInString(str string, data map[string]interface{}, extras []cel.EnvOption) (string, error) {
	var err error

	result := conditionalStatementRegex.ReplaceAllStringFunc(str, func(part string) string {
		part = strings.TrimSuffix(strings.TrimPrefix(part, "{{"), "}}")
		result, innErr := Run[string](part, data, extras)
		if innErr != nil {
			err = innErr
			return err.Error()
		}
		return result
	})

	return result, err
}
