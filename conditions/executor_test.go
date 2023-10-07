package conditions

import (
	"github.com/google/cel-go/cel"
	"testing"
)

func TestResolveIf(t *testing.T) {
	type args struct {
		condition interface{}
		data      map[string]interface{}
		extraCels []cel.EnvOption
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "null condition with nil data",
			args: args{
				condition: nil,
				data:      nil,
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "null condition with empty data",
			args: args{
				condition: nil,
				data:      map[string]interface{}{},
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "null condition with success true flag",
			args: args{
				condition: nil,
				data:      map[string]interface{}{"success": true},
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "empty condition with success true flag",
			args: args{
				condition: "",
				data:      map[string]interface{}{"success": true},
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "null condition with success false flag",
			args: args{
				condition: nil,
				data:      map[string]interface{}{"success": false},
				extraCels: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty condition with success false flag",
			args: args{
				condition: "",
				data:      map[string]interface{}{"success": false},
				extraCels: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success condition with success false flag",
			args: args{
				condition: "success",
				data:      map[string]interface{}{"success": false},
				extraCels: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success condition with success false flag",
			args: args{
				condition: "success",
				data:      map[string]interface{}{"success": false},
				extraCels: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "string condition with variable true",
			args: args{
				condition: "loader == \"vanilla\"",
				data:      map[string]interface{}{"loader": "vanilla"},
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "string condition with variable false",
			args: args{
				condition: "loader == \"vanilla\"",
				data:      map[string]interface{}{"loader": "notvanilla"},
				extraCels: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "string condition with variable os",
			args: args{
				condition: "os == \"linux\"",
				data:      nil,
				extraCels: nil,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveIf(tt.args.condition, tt.args.data, tt.args.extraCels)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveIf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveIf() got = %v, want %v", got, tt.want)
			}
		})
	}
}
