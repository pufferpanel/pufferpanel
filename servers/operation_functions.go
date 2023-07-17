package servers

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/pufferpanel/pufferpanel/v3"
	"os"
	"os/exec"
	"path/filepath"
)

func CreateFunctions(env pufferpanel.Environment) []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function("file_exists",
			cel.Overload("file_exists_string_bool",
				[]*cel.Type{cel.StringType},
				cel.BoolType,
				cel.UnaryBinding(cel_file_exists(env)),
			)),
		cel.Function("in_path",
			cel.Overload("in_path_string_bool",
				[]*cel.Type{cel.StringType},
				cel.BoolType,
				cel.UnaryBinding(cel_in_path(env)),
			)),
	}
}

func cel_file_exists(env pufferpanel.Environment) functions.UnaryOp {
	return func(fileName ref.Val) ref.Val {
		fullPath := filepath.Join(env.GetRootDirectory(), fileName.Value().(string))
		_, err := os.Stat(fullPath)
		return types.Bool(err == nil)
	}
}

func cel_in_path(environment pufferpanel.Environment) functions.UnaryOp {
	return func(fileName ref.Val) ref.Val {
		_, err := exec.LookPath(fileName.Value().(string))
		return types.Bool(err == nil)
	}
}
