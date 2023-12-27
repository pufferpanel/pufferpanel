package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func Test_tokenService_GenerateRequest(t *testing.T) {
	type args struct {
		request interface{}
	}

	desiredObj := map[string]interface{}{
		"test": "12345",
		"time": 45201,
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestGenerate",
			args:    args{request: desiredObj},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, err := NewTokenService()
			if !assert.NoError(t, err) {
				return
			}
			got, err := ts.GenerateRequest(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			data, err := io.ReadAll(got)
			if !assert.NoError(t, err) {
				return
			}
			fmt.Println(string(data))

			//now check if we can decrypt it
			response, err := ts.DecryptRequest(bytes.NewReader(data))
			if !assert.NoError(t, err) {
				return
			}

			var res map[string]interface{}
			err = json.NewDecoder(response).Decode(&res)
			if !assert.NoError(t, err) {
				return
			}
		})
	}
}
