package spongedl

import (
	"testing"
)

func TestSpongeDl_Run(t *testing.T) {
	type fields struct {
		Recommended      bool
		SpongeType       string
		SpongeVersion    string
		MinecraftVersion string
	}

	testEnv := &TestEnvironment{}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "SpongeVanilla_1.16.5",
			fields: fields{
				MinecraftVersion: "1.16.5",
				SpongeType: "spongevanilla",
			},
			wantErr: false,
		},
		{
			name: "SpongeVanilla_1.17.1",
			fields: fields{
				MinecraftVersion: "1.17.1",
				SpongeType: "spongevanilla",
			},
			wantErr: false,
		},
		{
			name: "SpongeVanilla_Bad",
			fields: fields{
				MinecraftVersion: "1.0.0",
				SpongeType: "spongevanilla",
			},
			wantErr: true,
		},
		{
			name: "SpongeForge_1.16.5",
			fields: fields{
				MinecraftVersion: "1.16.5",
				SpongeType: "spongeforge",
			},
			wantErr: false,
		},
		{
			name: "SpongeForge_Bad",
			fields: fields{
				MinecraftVersion: "1.0.0",
				SpongeType: "spongeforge",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := SpongeDl{
				Recommended:      tt.fields.Recommended,
				SpongeType:       tt.fields.SpongeType,
				SpongeVersion:    tt.fields.SpongeVersion,
				MinecraftVersion: tt.fields.MinecraftVersion,
			}
			if err := op.Run(testEnv); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
