package pufferpanel

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/spf13/cast"
	"io"
	"net"
	"time"
)

type TelnetConnection struct {
	io.WriteCloser
	IP         string
	Port       string
	Password   string
	Reconnect  bool
	connection *net.TCPConn
	ready      bool
}

func (tc *TelnetConnection) Write(p []byte) (n int, err error) {
	if !tc.ready {
		time.Sleep(1 * time.Second)
		if !tc.ready {
			return 0, errors.New("telnet not available")
		}
	}
	if tc.connection != nil {
		return tc.connection.Write(p)
	}
	return 0, errors.New("telnet not available")
}

func (tc *TelnetConnection) Start() {
	tc.Reconnect = true
	if tc.IP == "" {
		tc.IP = "127.0.0.1"
	}

	go tc.reconnector()
}

func (tc *TelnetConnection) Close() error {
	tc.Reconnect = false
	if tc.connection == nil {
		return nil
	}
	return tc.connection.Close()
}

func (tc *TelnetConnection) reconnector() {
	init := true
	for tc.Reconnect {
		tc.ready = false
		if !init {
			time.Sleep(5 * time.Second)
		} else {
			init = false
		}

		ipAddr := &net.TCPAddr{
			IP:   net.ParseIP(tc.IP),
			Port: cast.ToInt(tc.Port),
		}
		conn, err := net.DialTCP("tcp", nil, ipAddr)
		if err != nil {
			logging.Debug.Printf("Error waiting for TCP TELNET socket: %s", err.Error())
			continue
		}
		_ = conn.SetKeepAlive(true)

		//wait a second for the prompt for passwords/other delays
		time.Sleep(1 * time.Second)

		if tc.Password != "" {
			_, err = conn.Write([]byte(tc.Password + "\n"))
			if err != nil {
				logging.Debug.Printf("Error sending password for TCP TELNET socket: %s", err.Error())
				continue
			}
		}
		tc.connection = conn
		tc.ready = true
		listening := true
		for listening {
			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
			if err != nil {
				listening = false
			}
		}
	}
}
