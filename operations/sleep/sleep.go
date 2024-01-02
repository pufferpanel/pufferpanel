package sleep

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"time"
)

type Sleep struct {
	Duration time.Duration
}

func (d Sleep) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	time.Sleep(d.Duration)
	return pufferpanel.OperationResult{Error: nil}
}
