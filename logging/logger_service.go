//go:build windows
// +build windows

package logging

import (
	"golang.org/x/sys/windows/svc/eventlog"
	"io"
)

var elog *eventlog.Log

func CreateServiceLogger(t string) io.Writer {
	if elog == nil {
		var err error
		elog, err = eventlog.Open("PufferPanel")
		if err != nil {
			panic(err)
		}
	}

	return &eventLogWriter{Type: t}
}

type eventLogWriter struct {
	Type string
}

func (e *eventLogWriter) Write(p []byte) (n int, err error) {
	switch e.Type {
	case "error":
		err = elog.Error(1, string(p))
		break
	case "info":
		err = elog.Info(1, string(p))
	}

	return len(p), err
}
