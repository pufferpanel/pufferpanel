package servers

import (
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/files"
)

func CreateFunctions(server *Server, env *pufferpanel.Environment) []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function("file_exists",
			cel.Overload("file_exists_string_bool",
				[]*cel.Type{cel.StringType},
				cel.BoolType,
				cel.UnaryBinding(cel_file_exists(server.GetFileServer(), env)),
			)),
		cel.Function("in_path",
			cel.Overload("in_path_string_bool",
				[]*cel.Type{cel.StringType},
				cel.BoolType,
				cel.UnaryBinding(cel_in_path(env)),
			)),
		cel.Function("is_server_running", cel.Overload("is_server_running_bool",
			[]*cel.Type{},
			cel.BoolType,
			cel.FunctionBinding(cel_is_server_running(env)),
		)),
	}
}

func cel_file_exists(serverFS files.FileServer, env *pufferpanel.Environment) functions.UnaryOp {
	return func(fileName ref.Val) ref.Val {
		fullPath := filepath.Join(fileName.Value().(string))
		_, err := serverFS.Stat(fullPath)
		return types.Bool(err == nil)
	}
}

func cel_in_path(env *pufferpanel.Environment) functions.UnaryOp {
	return func(fileName ref.Val) ref.Val {
		_, err := exec.LookPath(fileName.Value().(string))
		return types.Bool(err == nil || errors.Is(err, exec.ErrDot))
	}
}

func cel_is_server_running(env *pufferpanel.Environment) functions.FunctionOp {
	return func(values ...ref.Val) ref.Val {
		r, err := env.IsRunning()
		return types.Bool(err == nil && r)
	}
}
