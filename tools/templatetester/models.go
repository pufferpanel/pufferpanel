package main

type TestScenario struct {
	Name string
	Test *TestTemplate
}

type TestTemplate struct {
	Template           []byte
	Name               string
	Variables          map[string]interface{}
	Environment        map[string]interface{}
	IgnoreExitCode     bool
	RuntimeRequirement int
}

type TestData struct {
	Name               string                 `json:"name"`
	Variables          map[string]interface{} `json:"variables"`
	Environment        map[string]interface{} `json:"environment"`
	IgnoreExitCode     bool                   `json:"ignoreExitCode"`
	RuntimeRequirement int                    `json:"runtimeRequirement"`
}
