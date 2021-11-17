// +build windows

package main

import (
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"golang.org/x/sys/windows/svc"
)

func serviceCheck(c chan error) {
	inService, _ := svc.IsWindowsService()
	if inService {
		logging.Debug.Printf("Detecting this is running as a service\n")
		c <- svc.Run("PufferPanel", &service{})
	}
}

type service struct{}

func (m *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				logging.Info.Printf("Received stop command\n")
				break loop
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}
