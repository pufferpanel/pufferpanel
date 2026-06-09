package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fileServer_Symlink(t *testing.T) {
	tests := []struct {
		name    string
		oldpath string
		newpath string
		wantErr error
	}{
		{
			name:    "symlink out of root",
			oldpath: "testfile",
			newpath: "../escape",
			wantErr: ErrInvalidAccess,
		},
		{
			name:    "symlink within allowed folders",
			oldpath: "internal/testfile",
			newpath: "../escape",
			wantErr: nil,
		},
		{
			name:    "symlink to a root file resolves to localized",
			oldpath: "testfile",
			newpath: "/escape",
			wantErr: nil,
		},
	}

	p, err := os.MkdirTemp("", "*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(p)

	sfp, err := NewFileServer(p, 0, 0)
	if !assert.NoError(t, err) {
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := filepath.Dir(tt.oldpath)
			if parent != "" {
				os.MkdirAll(filepath.Join(p, parent), 0755)
			}
			defer func() {
				f, _ := os.ReadDir(p)
				for _, v := range f {
					z := filepath.Join(p, v.Name())
					os.RemoveAll(z)
				}
			}()

			gotErr := sfp.Symlink(tt.oldpath, tt.newpath)
			if gotErr != tt.wantErr {
				t.Errorf("Symlink() failed: got: %v, wanted: %v", gotErr, tt.wantErr)
				return
			}

			if gotErr == nil {
				fmt.Printf("Confirming this did not bypass and write to %s\n", tt.newpath)
				if !assert.NoFileExists(t, tt.newpath) {
					return
				}

				writtenPath := filepath.Join(p, filepath.Dir("/"+tt.newpath), filepath.Base(tt.newpath))
				fmt.Printf("Opening %s\n", writtenPath)
				result, err := os.Readlink(writtenPath)
				if !assert.NoError(t, err) {
					return
				}
				if filepath.IsAbs(result) {
					if assert.Truef(t, strings.HasPrefix(result, p), "Expected prefix of %s but is missing: %s", p, result) {
						return
					}
				}
			}
		})
	}
}
