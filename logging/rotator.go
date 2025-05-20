package logging

import (
	"io"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

type Rotator struct {
	sync.RWMutex
	io.WriteCloser
	backer io.WriteCloser
}

func (r *Rotator) Write(p []byte) (n int, err error) {
	r.RLock()
	defer r.RUnlock()
	if r.backer == nil {
		return len(p), nil
	}
	return r.backer.Write(p)
}

func (r *Rotator) Close() error {
	if r.backer == nil {
		return nil
	}
	return r.backer.Close()
}

func (r *Rotator) Rotate(newBackend io.WriteCloser) {
	r.Lock()
	defer r.Unlock()
	oldBacker := r.backer
	r.backer = newBackend
	_ = oldBacker.Close()
}

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
