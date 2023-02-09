//go:build !windows

package logging

import (
	"os"
	"os/signal"
	"path"
	"syscall"
)

func (r *Rotator) StartRotation(dir string) {
	go func(directory string) {
		sig := make(chan os.Signal, 1)
		for {
			signal.Notify(sig, syscall.SIGUSR1)

			<-sig

			newFile, err := os.OpenFile(path.Join(directory, "pufferpanel.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}

			r.Rotate(newFile)
		}
	}(dir)
}
