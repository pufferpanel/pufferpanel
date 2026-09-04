package pufferpanel

import (
	"github.com/pufferpanel/pufferpanel/v3/files"
)

type Server struct {
	Type
	Identifier            string                    `json:"id,omitempty"`
	Display               string                    `json:"display,omitempty"`
	Icon                  string                    `json:"icon,omitempty"`
	Variables             map[string]Variable       `json:"data"`
	Groups                []Group                   `json:"groups,omitempty"`
	Installation          []ConditionalMetadataType `json:"install"`
	Uninstallation        []ConditionalMetadataType `json:"uninstall"`
	Execution             Execution                 `json:"run"`
	Environment           MetadataType              `json:"environment"`
	SupportedEnvironments []MetadataType            `json:"supportedEnvironments,omitempty"`
	Requirements          Requirements              `json:"requirements,omitempty"`
	Stats                 MetadataType              `json:"stats,omitempty"`
	Query                 MetadataType              `json:"query,omitempty"`
	KeepAlive             KeepAlive                 `json:"keepAlive,omitempty"`
} //@name ServerDefinition

type Execution struct {
	Command                 interface{}               `json:"command"`
	StopCommand             string                    `json:"stop,omitempty"`
	StopCode                int                       `json:"stopCode,omitempty"`
	PreExecution            []ConditionalMetadataType `json:"pre,omitempty"`
	PostExecution           []ConditionalMetadataType `json:"post,omitempty"`
	EnvironmentVariables    map[string]string         `json:"environmentVars,omitempty"`
	WorkingDirectory        string                    `json:"workingDirectory,omitempty"`
	Stdin                   StdinConsoleConfiguration `json:"stdin,omitempty"`
	AutoStart               bool                      `json:"autostart"`
	AutoRestartFromCrash    bool                      `json:"autorecover"`
	AutoRestartFromGraceful bool                      `json:"autorestart"`
	ExpectedExitCode        int                       `json:"expectedExitCode,omitempty"`
} //@name Execution

type Name struct {
	Name string `json:"name"`
} //@name Name

type Command struct {
	If      string                    `json:"if,omitempty"`
	Command string                    `json:"command"`
	StdIn   StdinConsoleConfiguration `json:"stdin"`
} //@name Command

type Type struct {
	Type string `json:"type"`
} //@name Type

type Group struct {
	If          string   `json:"if,omitempty"`
	Display     string   `json:"display"`
	Description string   `json:"description"`
	Variables   []string `json:"variables"`
	Order       int      `json:"order"`
} //@name Group

type KeepAlive struct {
	Frequency string `json:"frequency"`
	Command   string `json:"command"`
} //@name KeepAlive

func (s *Server) CopyFrom(replacement *Server) {
	s.Variables = replacement.Variables
	s.Type = replacement.Type
	s.Execution = replacement.Execution
	s.Display = replacement.Display
	s.Installation = replacement.Installation
	s.Uninstallation = replacement.Uninstallation
	s.Environment = replacement.Environment
	s.Requirements = replacement.Requirements
	s.SupportedEnvironments = replacement.SupportedEnvironments
	s.Groups = replacement.Groups
	s.Stats = replacement.Stats
}

func (s *Server) DataToMap() map[string]interface{} {
	var result = make(map[string]interface{})

	for k, v := range s.Variables {
		result[k] = v.Value
	}

	return result
}

type DaemonServer interface {
	Get() Server

	GetFileServer() files.FileServer

	Extract(source, destination string) error

	ArchiveItems(files []string, destination string) error

	DataToMap() map[string]interface{}

	Id() string
}
