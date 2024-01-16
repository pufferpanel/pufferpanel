package conditions

import (
	"github.com/google/cel-go/cel"
	"runtime"
	"testing"
)

func TestResolveIf(t *testing.T) {
	type args struct {
		condition string
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
			name: "string condition with variable true using 's",
			args: args{
				condition: "loader == 'vanilla'",
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
				condition: "os == \"" + runtime.GOOS + "\"",
				data:      nil,
				extraCels: nil,
			},
			want:    true,
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

func TestReplaceInString(t *testing.T) {
	type args struct {
		str    string
		data   map[string]interface{}
		extras []cel.EnvOption
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "SimpleReplace",
			args: args{
				str:    "Hello {{ 'world!' }}",
				data:   nil,
				extras: nil,
			},
			want: "Hello world!",
		},
		{
			name: "VariableReplace",
			args: args{
				str: "Hello {{ world }}",
				data: map[string]interface{}{
					"world": "world!",
				},
				extras: nil,
			},
			want: "Hello world!",
		},
		{
			name: "Multiple Variable Replace",
			args: args{
				str: "{{ hello }} {{ world }}",
				data: map[string]interface{}{
					"hello": "Hello",
					"world": "world!",
				},
				extras: nil,
			},
			want: "Hello world!",
		},
		{
			name: "Logic Replace",
			args: args{
				str: "{{ ishello ? 'Hello' : 'Goodbye'}} world!",
				data: map[string]interface{}{
					"ishello": true,
				},
				extras: nil,
			},
			want: "Hello world!",
		},
		{
			name: "No Replace",
			args: args{
				str: "Hello world!",
				data: map[string]interface{}{
					"world": "and bye",
				},
				extras: nil,
			},
			want: "Hello world!",
		},
		{
			name: "Invalid",
			args: args{
				str: "Hello {{ asdf }}!",
				data: map[string]interface{}{
					"world": "and bye",
				},
				extras: nil,
			},
			want:    "Hello ERROR: <input>:1:2: undeclared reference to 'asdf' (in container '')\n |  asdf \n | .^!",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceInString(tt.args.str, tt.args.data, tt.args.extras)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceInString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReplaceInString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
