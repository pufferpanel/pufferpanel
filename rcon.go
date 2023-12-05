package pufferpanel

import (
	"errors"
	"fmt"
	"github.com/gorcon/rcon"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"io"
	"time"
)

type RCONConnection struct {
	io.WriteCloser
	IP         string
	Port       string
	Password   string
	Reconnect  bool
	connection *rcon.Conn
	ready      bool
	closer     chan error
}

func (tc *RCONConnection) Write(p []byte) (n int, err error) {
	if !tc.ready {
		time.Sleep(1 * time.Second)
		if !tc.ready {
			return 0, errors.New("rcon not available")
		}
	}
	if tc.connection != nil {
		var res string
		res, err = tc.connection.Execute(string(p))
		if err != nil {
			tc.closer <- err
		}
		logging.Debug.Printf("RCON Response: %s\n", res)
		return len(p), err
	}
	return 0, errors.New("rcon not available")
}

func (tc *RCONConnection) Start() {
	tc.Reconnect = true
	if tc.IP == "" {
		tc.IP = "127.0.0.1"
	}

	go tc.reconnector()
}

func (tc *RCONConnection) Close() error {
	tc.Reconnect = false
	if tc.connection == nil {
		return nil
	}
	return tc.connection.Close()
}

func (tc *RCONConnection) reconnector() {
	init := true
	for tc.Reconnect {
		tc.recon(init)
		if init {
			init = false
		}
	}
}

func (tc *RCONConnection) recon(init bool) {
	tc.ready = false
	if !init {
		time.Sleep(5 * time.Second)
	}

	var err error
	tc.connection, err = rcon.Dial(fmt.Sprintf("%s:%s", tc.IP, tc.Port), tc.Password)
	if err != nil {
		logging.Debug.Printf("Error sending password for TCP TELNET socket: %s", err.Error())
		return
	}
	defer tc.connection.Close()
	<-tc.closer
}
