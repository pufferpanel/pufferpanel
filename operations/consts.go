package operations

import (
	"github.com/google/cel-go/cel"
	"runtime"
)

var celGlobalConstants = []cel.EnvOption{
	cel.Variable(cel_os, cel.StringType),
	cel.Variable(cel_arch, cel.StringType),
	cel.Variable(cel_success, cel.BoolType),
	cel.Variable(cel_env, cel.StringType),
}

var celGlobalConstantValues = map[string]interface{}{
	cel_os:   runtime.GOOS,
	cel_arch: runtime.GOARCH,
}

var (
	cel_os      = "os"
	cel_arch    = "arch"
	cel_success = "success"
	cel_env     = "env"
)
