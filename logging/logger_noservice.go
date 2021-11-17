// +build !windows

package logging

import (
	"io"
)

func CreateServiceLogger(string) io.Writer {
	return nil
}