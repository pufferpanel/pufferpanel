package tests

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"github.com/pufferpanel/pufferpanel/v3/web/daemon"
	"github.com/stretchr/testify/assert"
	"net/http"
	"runtime"
	"testing"
)

func Test_getFeatures(t *testing.T) {

	tests := []struct {
		name       string
		expected   daemon.Features
		dockerFlag bool
	}{
		{
			name: "docker not forced",
			expected: daemon.Features{
				Features:     []string{"docker"},
				Environments: []string{"docker", "host", "tty", "standard"},
				OS:           runtime.GOOS,
				Arch:         runtime.GOARCH,
				Version:      pufferpanel.Version,
			},
			dockerFlag: false,
		},
		{
			name: "docker forced",
			expected: daemon.Features{
				Features:     []string{"docker"},
				Environments: []string{"docker"},
				OS:           runtime.GOOS,
				Arch:         runtime.GOARCH,
				Version:      pufferpanel.Version,
			},
			dockerFlag: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = config.DockerDisallowHost.Set(tt.dockerFlag, false)

			if runtime.GOOS == "windows" {
				tt.expected.Environments = utils.Remove(tt.expected.Environments, "tty")
			}

			result := CallAPI("GET", "/daemon/features", nil, "")
			if !assert.Equal(t, http.StatusOK, result.Code) {
				return
			}

			var results daemon.Features
			err := json.NewDecoder(result.Body).Decode(&results)
			if !assert.NoError(t, err) {
				return
			}
			assert.ElementsMatch(t, tt.expected.Environments, results.Environments)
		})
	}
}
