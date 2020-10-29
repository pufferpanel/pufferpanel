package legacy

type ServerJson struct {
	ProgramData ProgramData `json:"pufferd"`
}

type ProgramData struct {
	Data            map[string]DataObject  `json:"data,omitempty"`
	Display         string                 `json:"display,omitempty"`
	EnvironmentData map[string]interface{} `json:"environment,omitempty"`
	InstallData     InstallSection         `json:"install,omitempty"`
	UninstallData   InstallSection         `json:"uninstall,omitempty"`
	Type            string                 `json:"type,omitempty"`
	Identifier      string                 `json:"id,omitempty"`
	RunData         RunObject              `json:"run,omitempty"`
	Template        string                 `json:"template,omitempty"`
}

type DataObject struct {
	Description  string      `json:"desc,omitempty"`
	Display      string      `json:"display,omitempty"`
	Internal     bool        `json:"internal,omitempty"`
	Required     bool        `json:"required,omitempty"`
	Value        interface{} `json:"value,omitempty"`
	UserEditable bool        `json:"userEdit,omitempty"`
}

type RunObject struct {
	Arguments               []string                 `json:"arguments,omitempty"`
	Program                 string                   `json:"program,omitempty"`
	Stop                    string                   `json:"stop,omitempty"`
	Enabled                 bool                     `json:"enabled,omitempty"`
	AutoStart               bool                     `json:"autostart,omitempty"`
	AutoRestartFromCrash    bool                     `json:"autorecover,omitempty"`
	AutoRestartFromGraceful bool                     `json:"autorestart,omitempty"`
	Pre                     []map[string]interface{} `json:"pre,omitempty"`
	Post                    []map[string]interface{} `json:"post,omitempty"`
	StopCode                int                      `json:"stopCode,omitempty"`
	EnvironmentVariables    map[string]string        `json:"environmentVars,omitempty"`
}

type InstallSection struct {
	Operations []map[string]interface{} `json:"commands,,omitempty"`
}
