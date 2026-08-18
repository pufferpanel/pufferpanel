package alterfile

import (
	"bytes"
	"io"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_replaceLinesInFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		file           *MemoryFile
		regex          *regexp.Regexp
		needle         string
		replacement    string
		wantErr        bool
		expectedResult []byte
	}{
		{
			name: "BasicReplace",
			file: &MemoryFile{
				start: []byte("test {replacement} in this"),
			},
			regex:          nil,
			needle:         "{replacement}",
			replacement:    "world",
			expectedResult: []byte("test world in this\n"),
		},
		{
			name: "BasicReplaceMultiline",
			file: &MemoryFile{
				start: []byte("line 1 should not change\ntest {replacement} in this\nand this one is as-is"),
			},
			regex:          nil,
			needle:         "{replacement}",
			replacement:    "world",
			expectedResult: []byte("line 1 should not change\ntest world in this\nand this one is as-is\n"),
		},
		{
			name: "BasicRegex",
			file: &MemoryFile{
				start: []byte("test {replacement} in this"),
			},
			regex:          regexp.MustCompile(`\{replacement\}`),
			needle:         "",
			replacement:    "world",
			expectedResult: []byte("test world in this\n"),
		},
		{
			name: "BasicRegexMultiline",
			file: &MemoryFile{
				start: []byte("line 1 should not change\ntest {replacement} in this\nand this one is as-is"),
			},
			regex:          regexp.MustCompile(`\{replacement\}`),
			needle:         "",
			replacement:    "world",
			expectedResult: []byte("line 1 should not change\ntest world in this\nand this one is as-is\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := replaceLinesInFile(tt.file, tt.regex, []byte(tt.needle), []byte(tt.replacement))
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("replaceLinesInFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("replaceLinesInFile() succeeded unexpectedly")
			}
			if !assert.Equal(t, tt.expectedResult, tt.file.end) {
				return
			}
		})
	}
}

type MemoryFile struct {
	io.ReadWriter

	start       []byte
	index       int
	end         []byte
	readWrapper *bytes.Reader
}

func (m *MemoryFile) Read(p []byte) (n int, err error) {
	if m.readWrapper == nil {
		m.readWrapper = bytes.NewReader(m.start)
	}
	return m.readWrapper.Read(p)
}

func (m *MemoryFile) Write(p []byte) (n int, err error) {
	m.end = append(m.end, p...)
	return len(p), nil
}
