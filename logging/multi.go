// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//This is an altered variant of the MultiWriter provided by Golang.
//This adds a null-check to adding writings, along with holding off
//on returning when doing a write so that things don't look stupid....
//because windows

package logging

import "io"

type multiWriter struct {
	writers []io.Writer
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

var _ io.StringWriter = (*multiWriter)(nil)

func (t *multiWriter) WriteString(s string) (n int, err error) {
	var p []byte // lazily initialized if/when needed
	var firstErr error
	for _, w := range t.writers {
		if sw, ok := w.(io.StringWriter); ok {
			n, err = sw.WriteString(s)
		} else {
			if p == nil {
				p = []byte(s)
			}
			n, err = w.Write(p)
		}
		if err == nil && n != len(s) {
			err = io.ErrShortWrite
		}
		if firstErr == nil && err != nil {
			firstErr = err
		}
	}
	return len(s), firstErr
}

// MultiWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command.
//
// Each write is written to each listed writer, one at a time.
// If a listed writer returns an error, that overall write operation
// stops and returns the error; it does not continue down the list.
func MultiWriter(writers ...io.Writer) io.Writer {
	allWriters := make([]io.Writer, 0)
	for _, w := range writers {
		if w == nil {
			continue
		}
		if mw, ok := w.(*multiWriter); ok {
			allWriters = append(allWriters, mw.writers...)
		} else {
			allWriters = append(allWriters, w)
		}
	}
	return &multiWriter{allWriters}
}
