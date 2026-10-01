package utils_test

import (
	"testing"

	"github.com/pufferpanel/pufferpanel/v3/utils"
)

func TestParsePortMap(t *testing.T) {
	tests := []struct {
		str               string
		wantIp            string
		wantHostPort      string
		wantContainerPort string
		wantProtocol      string
		wantErr           bool
	}{
		{
			str:               "80",
			wantIp:            "0.0.0.0",
			wantHostPort:      "80",
			wantContainerPort: "80",
			wantProtocol:      "tcp",
			wantErr:           false,
		},
		{
			str:               "80/tcp",
			wantIp:            "0.0.0.0",
			wantHostPort:      "80",
			wantContainerPort: "80",
			wantProtocol:      "tcp",
			wantErr:           false,
		},
		{
			str:               "80/udp",
			wantIp:            "0.0.0.0",
			wantHostPort:      "80",
			wantContainerPort: "80",
			wantProtocol:      "udp",
			wantErr:           false,
		},
		{
			str:               "1234:80",
			wantIp:            "0.0.0.0",
			wantHostPort:      "1234",
			wantContainerPort: "80",
			wantProtocol:      "tcp",
			wantErr:           false,
		},
		{
			str:               "1234:80/udp",
			wantIp:            "0.0.0.0",
			wantHostPort:      "1234",
			wantContainerPort: "80",
			wantProtocol:      "udp",
			wantErr:           false,
		},
		{
			str:               "192.168.1.1:1234:80",
			wantIp:            "192.168.1.1",
			wantHostPort:      "1234",
			wantContainerPort: "80",
			wantProtocol:      "tcp",
			wantErr:           false,
		},
		{
			str:               "192.168.1.1:1234:80/udp",
			wantIp:            "192.168.1.1",
			wantHostPort:      "1234",
			wantContainerPort: "80",
			wantProtocol:      "udp",
			wantErr:           false,
		},
		{
			str:               "[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1234:80/udp",
			wantIp:            "2001:db8:85a3:8d3:1319:8a2e:370:7348",
			wantHostPort:      "1234",
			wantContainerPort: "80",
			wantProtocol:      "udp",
			wantErr:           false,
		},
		{
			str:               "[::1]:1234:80/udp",
			wantIp:            "::1",
			wantHostPort:      "1234",
			wantContainerPort: "80",
			wantProtocol:      "udp",
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			got, gotErr := utils.ParsePortMap(tt.str)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParsePortMap() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParsePortMap() succeeded unexpectedly")
			}

			for k, v := range got {
				if k.Port() != tt.wantContainerPort {
					t.Errorf("ParsePortMap() Container Port = %v, want %v", k.Port(), tt.wantContainerPort)
				}
				if len(v) != 1 {
					t.Errorf("ParsePortMap() = want 1 binding, got %v", len(v))
				}
				if v[0].HostPort != tt.wantHostPort {
					t.Errorf("ParsePortMap() Host Port = %v, want %v", v[0].HostPort, tt.wantHostPort)
				}
				if v[0].HostIP.String() != tt.wantIp {
					t.Errorf("ParsePortMap() Host IP = %v, want %v", v[0].HostIP.String(), tt.wantIp)
				}
			}
		})
	}
}
