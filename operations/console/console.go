package console

import "github.com/pufferpanel/pufferpanel/v3"

type Console struct {
	Text string
}

func (d Console) Run(env pufferpanel.Environment) error {
	env.DisplayToConsole(true, "Message: %s \n", d.Text)
	return nil
}
