package logging

import (
	"io"
	"sync"
)

type Rotator struct {
	io.WriteCloser
	backer io.WriteCloser
	lock   sync.RWMutex //exists so we can rotate without loss of data
}

func (r *Rotator) Write(p []byte) (n int, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if r == nil || r.backer == nil {
		return len(p), nil
	}
	return r.backer.Write(p)
}

func (r *Rotator) Close() error {
	if r == nil || r.backer == nil {
		return nil
	}
	return r.backer.Close()
}

func (r *Rotator) Rotate(newBackend io.WriteCloser) {
	r.lock.Lock()
	defer r.lock.Unlock()
	oldBacker := r.backer
	r.backer = newBackend
	_ = oldBacker.Close()
}
