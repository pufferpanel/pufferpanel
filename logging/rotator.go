package logging

import (
	"io"
	"sync"
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
