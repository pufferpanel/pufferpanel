package logging

import "io"

type Rotator struct {
	io.WriteCloser
	backer io.WriteCloser
}

func (r *Rotator) Write(p []byte) (n int, err error) {
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
	oldBacker := r.backer
	r.backer = newBackend
	_ = oldBacker.Close()
}
