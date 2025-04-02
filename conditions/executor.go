package conditions

import (
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/pufferpanel/pufferpanel/v3"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type RunData struct {
	Variables map[string]interface{}
	Server    map[string]interface{}
}

func CreateRunData(server pufferpanel.Server) RunData {
	return RunData{
		Variables: server.DataToMap(),
		Server: map[string]interface{}{
			"env": server.Environment.Type,
			"id":  server.Identifier,
		},
	}
}

func Run[T string | bool](statement string, data RunData, extras []cel.EnvOption) (T, error) {
	var res T

	//if we didn't define a statement, then set as success if the map has one
	if statement == "" {
		switch any(res).(type) {
		case bool:
			res = any(true).(T)
		case string:
			res = any("").(T)
		case int:
			res = any(0).(T)
		}
		return res, nil
	}

	celVars := extras
	if celVars == nil {
		celVars = make([]cel.EnvOption, 0)
	}
	celVars = append(celVars,
		cel.Variable("var", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("sys", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("server", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("os", cel.StringType),
		cel.Variable("arch", cel.StringType))

	inputData := map[string]interface{}{}
	for k, v := range data.Variables {
		inputData[k] = v
	}

	celEnv, err := cel.NewEnv(celVars...)
	if err != nil {
		return res, err
	}

	for k := range inputData {
		a, err := celEnv.Extend(cel.Variable(k, cel.DynType))
		if err != nil {
			continue
		}
		celEnv = a
	}

	inputData["os"] = runtime.GOOS
	inputData["arch"] = runtime.GOARCH
	inputData["sys"] = map[string]string{
		"type": runtime.GOOS,
		"arch": runtime.GOARCH,
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

func ReplaceInString(str string, data RunData, extras []cel.EnvOption) (string, error) {
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
