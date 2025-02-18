package stdin

import "github.com/pufferpanel/pufferpanel/v3"

type Stdin struct {
	Command string
}

func (d Stdin) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment

	running, err := env.IsRunning()
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	} else if !running {
		return pufferpanel.OperationResult{Error: pufferpanel.ErrServerOffline}
	}

	err = env.ExecuteInMainProcess(d.Command)
	return pufferpanel.OperationResult{Error: err}
}
