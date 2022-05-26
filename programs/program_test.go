package programs

import (
	"testing"
)

func TestProgram_valid(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "invalid json",
			data: "{\n  \"type\": \"\",\n  \"id\": \"902cb39e\",\n  \"run\": {\n    \"autostart\": true\n  }\n}",
			want: false,
		},
		{
			name: "valid json",
			data: "{\n  \"type\": \"minecraft-java\",\n  \"display\": \"Vanilla - Minecraft\",\n  \"id\": \"test\",\n  \"install\": [\n    {\n      \"type\": \"mojangdl\",\n      \"version\": \"latest\",\n      \"target\": \"server.jar\"\n    },\n    {\n      \"type\": \"writefile\",\n      \"text\": \"server-ip=${ip}\\nserver-port=${port}\",\n      \"target\": \"server.properties\"\n    }\n  ],\n  \"run\": {\n    \"stop\": \"stop\",\n    \"pre\": [],\n    \"post\": [],\n    \"command\": \"java -Xmx2048M -Dlog4j2.formatMsgNoLookups=true -jar server.jar nogui\"\n  },\n  \"environment\": {\n    \"type\": \"standard\"\n  },\n  \"data\": {\n    \"ip\": {\n      \"value\": \"0.0.0.0\",\n      \"required\": true,\n      \"desc\": \"What IP to bind the server to\",\n      \"display\": \"IP\",\n      \"internal\": false\n    },\n    \"port\": {\n      \"value\": \"25565\",\n      \"required\": true,\n      \"desc\": \"What port to bind the server to\",\n      \"display\": \"Port\",\n      \"internal\": false,\n      \"type\": \"integer\"\n    }\n  }\n}\n",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := LoadFromData("test", []byte(tt.data))
			if err != nil {
				t.Errorf("Failed loading test json: %s", err.Error())
			}

			if got := p.valid(); got != tt.want {
				t.Errorf("valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
