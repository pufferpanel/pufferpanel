package neoforgedl

import "testing"

func Test_getLatestForMCVersion(t *testing.T) {
	type args struct {
		minecraftVersion string
	}
	tests := []struct {
		name    string
		args    args
		wantVer bool
		wantErr bool
	}{
		{
			name:    "ValidateParser",
			args:    args{minecraftVersion: "1.20.4"},
			wantVer: true,
			wantErr: false,
		},
		{
			name:    "NotSupportedVersion",
			args:    args{minecraftVersion: "1.12.2"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getLatestForMCVersion(tt.args.minecraftVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLatestForMCVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantVer {
				if got == "" {
					t.Errorf("getLatestForMCVersion() got nothing, but wanted something")
				}
			} else {
				if got != "" {
					t.Errorf("getLatestForMCVersion() got %s, but wanted nothing", got)
				}
			}
		})
	}
}
