package conditions

import (
	"github.com/google/cel-go/cel"
	"runtime"
	"testing"
)

func TestResolveIf(t *testing.T) {
	type args struct {
		condition string
		data      RunData
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
				data:      RunData{Variables: map[string]interface{}{"success": true}},
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "empty condition with success false flag",
			args: args{
				condition: "",
				data:      RunData{Variables: map[string]interface{}{"success": false}},
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "success condition with success false flag",
			args: args{
				condition: "vars.success",
				data:      RunData{Variables: map[string]interface{}{"success": false}},
				extraCels: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success condition with success false flag",
			args: args{
				condition: "vars.success",
				data:      RunData{Variables: map[string]interface{}{"success": false}},
				extraCels: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "string condition with variable true",
			args: args{
				condition: "vars.loader == \"vanilla\"",
				data:      RunData{Variables: map[string]interface{}{"loader": "vanilla"}},
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "string condition with variable true using 's",
			args: args{
				condition: "vars.loader == 'vanilla'",
				data:      RunData{Variables: map[string]interface{}{"loader": "vanilla"}},
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "string condition with variable false",
			args: args{
				condition: "vars.loader == \"vanilla\"",
				data:      RunData{Variables: map[string]interface{}{"loader": "notvanilla"}},
				extraCels: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "string condition with variable type",
			args: args{
				condition: "sys.type == \"" + runtime.GOOS + "\"",
				data:      RunData{},
				extraCels: nil,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run[bool](tt.args.condition, tt.args.data, tt.args.extraCels)
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
				str: "Hello {{ vars.world }}",
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
				str: "{{ vars.hello }} {{ vars.world }}",
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
				str: "{{ vars.ishello ? 'Hello' : 'Goodbye'}} world!",
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
				str: "Hello {{ vars.asdf }}!",
				data: map[string]interface{}{
					"world": "and bye",
				},
				extras: nil,
			},
			want:    "Hello ERROR: <input>:1:2: undeclared reference to 'asdf' (in container '')\n |  asdf \n | .^!",
			wantErr: true,
		},
		{
			name: "Int variable return",
			args: args{
				str: "{{vars.item}}M",
				data: map[string]interface{}{
					"item": 4096,
				},
				extras: nil,
			},
			want: "4096M",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceInString(tt.args.str, RunData{
				Variables: tt.args.data,
			}, tt.args.extras)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceInString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("ReplaceInString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
