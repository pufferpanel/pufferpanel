package curseforge

import (
	"errors"
	test "github.com/pufferpanel/pufferpanel/v3/testing"
	"os"
	"testing"
)

func TestCurseForge_Run(t *testing.T) {
	tests := []struct {
		name    string
		fields  CurseForge
		wantErr bool
	}{
		/*{
			name: "All The Mods",
			fields:  CurseForge{ProjectId: 520914, FileId: 0, JavaBinary: "java", Key: os.Getenv("CURSEFORGE_KEY")},
			wantErr: false,
		},
		{
			name: "RLCraft",
			fields:  CurseForge{ProjectId: 285109, FileId: 0, JavaBinary: "java", Key: os.Getenv("CURSEFORGE_KEY")},
			wantErr: false,
		},
		{
			name:    "Pixelmon",
			fields:  CurseForge{ProjectId: 389615, FileId: 4352459, JavaBinary: "java"},
			wantErr: false,
		},
		{
			name: "Better MC Fabric",
			fields:  CurseForge{ProjectId: 452013, FileId: 0, JavaBinary: "java", Key: os.Getenv("CURSEFORGE_KEY")},
			wantErr: false,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			env := test.CreateEnvironment(tt.name)
			err := os.RemoveAll(env.GetRootDirectory())
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Error(err)
				return
			}
			err = os.MkdirAll(env.GetRootDirectory(), 0644)
			if err != nil {
				t.Error(err)
				return
			}
			if err := c.Run(env); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
