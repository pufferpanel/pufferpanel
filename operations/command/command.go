package command

import (
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

type Command struct {
	Commands []string
	Env      map[string]string
	StdIn    pufferpanel.ConsoleConfiguration
}

func (c Command) Run(env pufferpanel.Environment) error {
	for _, cmd := range c.Commands {
		logging.Info.Printf("Executing command: %s", cmd)
		env.DisplayToConsole(true, fmt.Sprintf("Executing: %s\n", cmd))
		cmdToExec, args := pufferpanel.SplitArguments(cmd)
		ch := make(chan error, 1)
		err := env.Execute(pufferpanel.ExecutionData{
			Command:     cmdToExec,
			Arguments:   args,
			Environment: c.Env,
			Callback: func(exitCode int) {
				if exitCode != 0 {
					ch <- errors.New("failed to run command")
				}
				ch <- nil
			},
			StdInConfig: c.StdIn,
		})
		if err != nil {
			return err
		}
		err = <-ch
		if err != nil {
			return err
		}
	}

	return nil
}
