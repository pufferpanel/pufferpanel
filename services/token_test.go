package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_tokenService_GenerateRequest(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestGenerate",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, err := NewTokenService()
			if !assert.NoError(t, err) {
				return
			}
			got, err := ts.GenerateRequest()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			//now check if we can decrypt it
			err = ts.ValidateRequest(got)
			if !assert.NoError(t, err) {
				return
			}
		})
	}
}
