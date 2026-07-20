package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceTokens(t *testing.T) {
	mapping := createSourceMap()

	resultTest := ReplaceTokens("TEST ${val1}", mapping, PlainReplace)
	assert.Equal(t, "TEST RESULT1", resultTest)

	resultTest = ReplaceTokens("TEST val1", mapping, PlainReplace)
	assert.Equal(t, "TEST val1", resultTest)

	resultTest = ReplaceTokens("TEST val1", mapping, PlainReplace)
	assert.Equal(t, "TEST val1", resultTest)
}

func TestReplaceTokensInMaps(t *testing.T) {
	mapping := createSourceMap()
	data := []string{
		"TEST ${val1}",
		"TEST2",
		"TEST ${1234567}",
		"${int}",
	}
	expected := []string{
		"TEST RESULT1",
		"TEST2",
		"TEST RESULT3",
		"436",
	}

	resultTest := ReplaceTokensInArr(data, mapping)
	assert.Equal(t, expected, resultTest)
}

func createSourceMap() map[string]interface{} {
	source := make(map[string]interface{})

	source["val1"] = "RESULT1"
	source["value2"] = "RESULT2"
	source["1234567"] = "RESULT3"
	source["val123"] = "RESULT4"
	source["int"] = 436

	return source
}

func TestSplitArguments(t *testing.T) {
	tests := []struct {
		name          string
		args          string
		wantCmd       string
		wantArguments []string
	}{
		{
			args:          "java -jar test.jar",
			wantCmd:       "java",
			wantArguments: []string{"-jar", "test.jar"},
		},
		{
			args:          "java -jar \"test.jar\"",
			wantCmd:       "java",
			wantArguments: []string{"-jar", "test.jar"},
		},
		{
			args:          "java -jar \"test this.jar\"",
			wantCmd:       "java",
			wantArguments: []string{"-jar", "test this.jar"},
		},
		{
			args:          "java -jar \"test this.jar\" noGui",
			wantCmd:       "java",
			wantArguments: []string{"-jar", "test this.jar", "noGui"},
		},
		{
			args:          "\"C:\\\\Program Files\\\\Java\\\\bin\\\\java.exe\" -jar \"test this.jar\" noGui",
			wantCmd:       "C:\\Program Files\\Java\\bin\\java.exe",
			wantArguments: []string{"-jar", "test this.jar", "noGui"},
		},
		{
			args:          "java",
			wantCmd:       "java",
			wantArguments: []string{},
		},
		{
			args:          "java -jar  server.jar",
			wantCmd:       "java",
			wantArguments: []string{"-jar", "server.jar"},
		},
		{
			args:          "java -jar      server.jar    ",
			wantCmd:       "java",
			wantArguments: []string{"-jar", "server.jar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotArguments := SplitArguments(tt.args)
			if gotCmd != tt.wantCmd {
				t.Errorf("SplitArguments() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
			if !reflect.DeepEqual(gotArguments, tt.wantArguments) {
				t.Errorf("SplitArguments() gotArguments = %v, want %v", gotArguments, tt.wantArguments)
			}
		})
	}
}

func TestMergeArguments(t *testing.T) {
	tests := []struct {
		args     []string
		expected string
	}{
		{
			expected: "java -jar test.jar",
			args:     []string{"java", "-jar", "test.jar"},
		},
		{
			expected: "java -jar \"test.jar\"",
			args:     []string{"java", "-jar", "\"test.jar\""},
		},
		{
			expected: "java -jar \"test this.jar\"",
			args:     []string{"java", "-jar", "test this.jar"},
		},
		{
			expected: "java -jar \"test this.jar\" noGui",
			args:     []string{"java", "-jar", "test this.jar", "noGui"},
		},
		{
			expected: "\"C:\\\\Program Files\\\\Java\\\\bin\\\\java.exe\" -jar \"test this.jar\" noGui",
			args:     []string{"\"C:\\\\Program Files\\\\Java\\\\bin\\\\java.exe\"", "-jar", "test this.jar", "noGui"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := MergeArguments(tt.args)
			if result != tt.expected {
				t.Errorf("MergeArguments() got = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShellReplace(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		data     map[string]any
		expected string
	}{
		{
			str:      "",
			data:     map[string]any{},
			expected: "",
		},
		{
			str:      "test",
			data:     map[string]any{},
			expected: "test",
		},
		{
			str:      "testing \" test",
			data:     map[string]any{},
			expected: "testing \" test",
		},
		{
			str: "testing \" test",
			data: map[string]any{
				"test": "FAKE",
			},
			expected: "testing \" test",
		},
		{
			str: "testing ${test}",
			data: map[string]any{
				"test": "FAKE",
			},
			expected: "testing FAKE",
		},
		{
			str: "${test}ing ${test}",
			data: map[string]any{
				"test": "FAKE",
			},
			expected: "FAKEing FAKE",
		},
		{
			str: "testing ${test}",
			data: map[string]any{
				"test": " > FAKE",
			},
			expected: "testing \" \\> FAKE\"",
		},
		{
			str: "testing ${test}",
			data: map[string]any{
				"test": "FAKE \\\"test data\"",
			},
			expected: "testing \"FAKE \\\\\\\"test data\\\"\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			result := ReplaceTokens(tt.str, tt.data, ShellReplace)
			if result != tt.expected {
				t.Errorf("ReplaceTokens() got = %v, want %v", result, tt.expected)
			}
		})
	}
}
