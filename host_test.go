package pufferpanel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHostname(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		expected string
	}{
		{
			name:     "IPV6 Localhost",
			request:  "[::1]",
			expected: "[::1]",
		},
		{
			name:     "IPV6 Localhost with port",
			request:  "[::1]:8080",
			expected: "[::1]",
		},
		{
			name:     "Machine IPv6",
			request:  "[fe80::b1d1:8c48:ca2d:8262%9]",
			expected: "[fe80::b1d1:8c48:ca2d:8262%9]",
		},
		{
			name:     "Machine IPv6 with port",
			request:  "[fe80::b1d1:8c48:ca2d:8262%9]:8080",
			expected: "[fe80::b1d1:8c48:ca2d:8262%9]",
		},
		{
			name:     "Domain",
			request:  "example.com",
			expected: "example.com",
		},
		{
			name:     "Domain with port",
			request:  "example.com:8080",
			expected: "example.com",
		},
		{
			name:     "IPV4",
			request:  "127.0.0.1",
			expected: "127.0.0.1",
		},
		{
			name:     "IPV4 with port",
			request:  "127.0.0.1:8080",
			expected: "127.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.expected, GetHostname(tt.request), "GetHostname(%v)", tt.request)
		})
	}
}
