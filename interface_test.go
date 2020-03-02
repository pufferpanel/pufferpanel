/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

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
