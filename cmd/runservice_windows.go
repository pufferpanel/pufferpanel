package main

import (
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/svc"
)

func executeRunService(cmd *cobra.Command, args []string) {
	_ = svc.Run("PufferPanel", &service{})
}

type service struct{}

func (m *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	term, success := internalRun()

	if !success || term == nil {
		return
	}

	go func() {
		<-term
		closePanel()
	}()

loop:
	for c := range r {
		switch c.Cmd {
		case svc.Stop, svc.Shutdown:
			logging.Info.Printf("Received stop command\n")
			break loop
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	term <- true
	return
}
