package pufferpanel

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestReplaceTokens(t *testing.T) {
	mapping := createSourceMap()

	resultTest := ReplaceTokens("TEST ${val1}", mapping)
	assert.Equal(t, "TEST RESULT1", resultTest)

	resultTest = ReplaceTokens("TEST val1", mapping)
	assert.Equal(t, "TEST val1", resultTest)

	resultTest = ReplaceTokens("TEST val1", mapping)
	assert.Equal(t, "TEST val1", resultTest)
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
			wantArguments: []string{"-jar", "\"test.jar\""},
		},
		{
			args:          "java -jar \"test this.jar\"",
			wantCmd:       "java",
			wantArguments: []string{"-jar", "\"test this.jar\""},
		},
		{
			args:          "java -jar \"test this.jar\" noGui",
			wantCmd:       "java",
			wantArguments: []string{"-jar", "\"test this.jar\"", "noGui"},
		},
		{
			args:          "\"C:\\\\Program Files\\\\Java\\\\bin\\\\java.exe\" -jar \"test this.jar\" noGui",
			wantCmd:       "\"C:\\Program Files\\Java\\bin\\java.exe\"",
			wantArguments: []string{"-jar", "\"test this.jar\"", "noGui"},
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
