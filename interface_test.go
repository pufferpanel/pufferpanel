package pufferpanel

import (
	"errors"
	"reflect"
	"testing"
)

var interfaceStringArray = []string{
	"12345",
	"true",
	"test3",
}

var interfaceObjectArray = []interface{}{
	12345,
	true,
	"test3",
}

var interfaceInvalid = map[string]string{
	"invalid": "input",
}

func TestToStringArray(t *testing.T) {
	type args struct {
		element interface{}
	}
	type wants struct {
		result []string
		err    error
	}

	tests := []struct {
		name string
		args args
		want wants
	}{
		/*{
			name: "Test null input",
			args: args{
				element: nil,
			},
			want: []string{},
		},
		{
			name: "Test empty input",
			args: args{
				element: make([]string, 0),
			},
			want: make([]string, 0),
		},
		{
			name: "Test all valid input",
			args: args{
				element: interfaceStringArray,
			},
			want: interfaceStringArray,
		},
		{
			name: "Test mixed input",
			args: args{
				element: interfaceObjectArray,
			},
			want: interfaceStringArray,
		},
		{
			name: "Test single string input",
			args: args{
				element: "test",
			},
			want: []string{"test"},
		},*/
		{
			name: "Test invalid type",
			args: args{
				element: interfaceInvalid,
			},
			want: wants{
				nil,
				errors.New("unable to cast map[string]string{\"invalid\":\"input\"} of type map[string]string to []string"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Convert(tt.args.element, []string{}); (err != nil && tt.want.err.Error() != err.Error()) || !reflect.DeepEqual(got, tt.want.result) {
				if tt.want.err != nil {
					t.Errorf("ToStringArray() = %v, expected %v", err, tt.want.err)
				} else {
					t.Errorf("ToStringArray() = %v, want %v", got, tt.want.result)
				}
			}
		})
	}
}
