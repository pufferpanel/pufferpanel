package pufferpanel

import (
	"reflect"
	"testing"
)

func TestParseRequirementRow(t *testing.T) {
	tests := []struct {
		str  string
		want []string
	}{
		{str: "", want: []string{}},
		{str: "linux", want: []string{"linux"}},
		{str: "  linux   ", want: []string{"linux"}},
		{str: "linux||windows", want: []string{"linux", "windows"}},
		{str: " linux    || windows  ", want: []string{"linux", "windows"}},
	}

	for _, tt := range tests {
		got := parseRequirementRow(tt.str)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseRequirementRow(%#v) = %v, want %v", tt.str, got, tt.want)
		}
	}
}
