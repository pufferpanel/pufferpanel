package pufferpanel

import "io"

type Console interface {
	io.WriteCloser
	Start()
}

type NoStartConsole struct {
	Console
	Base io.WriteCloser
}

func (nsc *NoStartConsole) Start() {}

func (nsc *NoStartConsole) Close() error {
	return nsc.Base.Close()
}

func (nsc *NoStartConsole) Write(p []byte) (n int, err error) {
	return nsc.Base.Write(p)
}
