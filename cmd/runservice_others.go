//go:build !windows

package main

import (
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/spf13/cobra"
	"net"
	"os"
)

const NotifyReady = "READY=1"
const NotifyShutdown = "STOPPING=1"

func executeRunService(cmd *cobra.Command, args []string) {
	term, success := internalRun()

	if !success || term == nil {
		return
	}

	notifySystemd(NotifyReady)

	<-term
	notifySystemd(NotifyShutdown)
	closePanel()
}

func notifySystemd(msg string) {
	//tell systemd we started
	socketAddr := &net.UnixAddr{
		Name: os.Getenv("NOTIFY_SOCKET"),
		Net:  "unixgram",
	}

	// NOTIFY_SOCKET not set
	if socketAddr.Name == "" {
		logging.Debug.Printf("NOTIFY_SOCKET not set, will not communicate to systemd")
		return
	}

	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	// Error connecting to NOTIFY_SOCKET
	if err != nil {
		logging.Error.Printf("Failed opening NOTIFY_SOCKET: %s", err.Error())
		return
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(msg)); err != nil {
		logging.Error.Printf("Failed sending message to NOTIFY_SOCKET: %s", err.Error())
		return
	}
}
