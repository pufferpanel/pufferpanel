package services

import "testing"

func TestLoadEmailService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Does it work",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadEmailService()
		})
	}
}
