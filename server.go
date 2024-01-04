package pufferpanel

import (
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

type Server struct {
	Type
	Identifier            string                    `json:"id,omitempty"`
	Display               string                    `json:"display,omitempty"`
	Icon                  string                    `json:"icon,omitempty"`
	Variables             map[string]Variable       `json:"data,omitempty"`
	Groups                []Group                   `json:"groups,omitempty"`
	Installation          []ConditionalMetadataType `json:"install,omitempty"`
	Uninstallation        []ConditionalMetadataType `json:"uninstall,omitempty"`
	Execution             Execution                 `json:"run"`
	Environment           MetadataType              `json:"environment"`
	SupportedEnvironments []MetadataType            `json:"supportedEnvironments,omitempty"`
	Requirements          Requirements              `json:"requirements,omitempty"`
} //@name ServerDefinition

type Task struct {
	Name         string                    `json:"name"`
	CronSchedule string                    `json:"cronSchedule"`
	Description  string                    `json:"description,omitempty"`
	Operations   []ConditionalMetadataType `json:"operations" binding:"required"`
} //@name Task

type Variable struct {
	Type
	Value        interface{}      `json:"value"`
	Display      string           `json:"display,omitempty"`
	Description  string           `json:"desc,omitempty"`
	Required     bool             `json:"required"`
	Internal     bool             `json:"internal,omitempty"`
	UserEditable bool             `json:"userEdit"`
	Options      []VariableOption `json:"options,omitempty"`
} //@name Variable
type variableAlias Variable

type VariableOption struct {
	Value   interface{} `json:"value"`
	Display string      `json:"display"`
} //@name VariableOption

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

type StdinConsoleConfiguration struct {
	Type     string `json:"type,omitempty"`
	IP       string `json:"ip,omitempty"`
	Port     string `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
} //@name StdinConsoleConfiguration

type stdinConfigAlias StdinConsoleConfiguration

type Type struct {
	Type string `json:"type"`
} //@name Type

type Requirements struct {
	OS       string   `json:"os,omitempty"`
	Arch     string   `json:"arch,omitempty"`
	Binaries []string `json:"binaries,omitempty"`
} //@name Requirements

type Group struct {
	If          string   `json:"if,omitempty"`
	Display     string   `json:"display"`
	Description string   `json:"description"`
	Variables   []string `json:"variables"`
	Order       int      `json:"order"`
} //@name Group

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
}

func (r Requirements) Test(server Server) error {
	osReq := parseRequirementRow(r.OS)
	if len(osReq) > 0 {
		passes := false
		for _, v := range osReq {
			if v == runtime.GOOS {
				passes = true
				break
			}
		}
		if !passes {
			return ErrUnsupportedOS(runtime.GOOS, strings.ReplaceAll(r.OS, "||", " OR "))
		}
	}

	archReq := parseRequirementRow(r.Arch)
	if len(archReq) > 0 {
		passes := false
		for _, v := range archReq {
			if v == runtime.GOARCH {
				passes = true
				break
			}
		}
		if !passes {
			return ErrUnsupportedArch(runtime.GOARCH, strings.ReplaceAll(r.Arch, "||", " OR "))
		}
	}

	//check to see if we support the environment
	//AKA.... if docker, do we support it
	var envType Type
	err := UnmarshalTo(server.Environment, &envType)
	if err != nil {
		return err
	}

	if envType.Type == "docker" {
		d, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return ErrDockerNotSupported
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = d.Ping(ctx)
		if err != nil {
			return ErrDockerNotSupported
		}
	} else {
		//we cannot check in docker if the binary requirements are good, so we'll skip it for docker
		//and check them now

		for _, v := range r.Binaries {
			binaries := parseRequirementRow(v)

			found := true
			for k, binary := range binaries {
				parsed := ReplaceTokens(binary, server.DataToMap())
				binaries[k] = parsed
				_, err := exec.LookPath(parsed)
				if err != nil {
					found = false
				}
			}
			if !found {
				return ErrMissingBinary(strings.Join(binaries, " OR "))
			}
		}
	}

	return nil
}

func (s *Server) DataToMap() map[string]interface{} {
	var result = make(map[string]interface{})

	for k, v := range s.Variables {
		result[k] = v.Value
	}
	result["serverId"] = s.Identifier

	return result
}

func parseRequirementRow(str string) []string {
	if str == "" {
		return []string{}
	}
	d := strings.Split(str, "||")
	for k, v := range d {
		d[k] = strings.TrimSpace(v)
	}
	return d
}

func (c *StdinConsoleConfiguration) Replace(variables map[string]interface{}) StdinConsoleConfiguration {
	return StdinConsoleConfiguration{
		Type:     c.Type,
		IP:       ReplaceTokens(c.IP, variables),
		Port:     ReplaceTokens(c.Port, variables),
		Password: ReplaceTokens(c.Password, variables),
	}
}

func (v *Variable) UnmarshalJSON(data []byte) (err error) {
	aux := variableAlias{}
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}
	if aux.Type.Type == "" {
		aux.Type = Type{Type: "string"}
	}

	//convert variable to correct typing
	switch aux.Type.Type {
	case "integer":
		{
			aux.Value, err = cast.ToIntE(aux.Value)
			if err != nil {
				var str string
				if str, err = cast.ToStringE(aux.Value); err == nil {
					if str == "" {
						aux.Value = 0
					}
				}
			}
		}
	case "boolean":
		{
			aux.Value, err = cast.ToBoolE(aux.Value)
		}
	}

	*v = Variable(aux)
	return
}

func (v *StdinConsoleConfiguration) UnmarshalJSON(data []byte) error {
	aux := stdinConfigAlias{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Type == "" {
		aux.Type = "stdin"
	}
	*v = StdinConsoleConfiguration(aux)
	return nil
}
