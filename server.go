/*
 Copyright 2019 Padduck, LLC
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
	"context"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

type Server struct {
	Type
	Variables             map[string]Variable `json:"data,omitempty"`
	Display               string              `json:"display,omitempty"`
	Environment           interface{}         `json:"environment,omitempty"`
	SupportedEnvironments interface{}         `json:"supportedEnvironments,omitempty"`
	Installation          []interface{}       `json:"install,omitempty"`
	Uninstallation        []interface{}       `json:"uninstall,omitempty"`
	Identifier            string              `json:"id,omitempty"`
	Execution             Execution           `json:"run,omitempty"`
	Requirements          Requirements        `json:"requirements,omitempty"`
	Groups                []Group             `json:"groups,omitempty"`
} //@name ServerDefinition

type Task struct {
	Name         string        `json:"name"`
	CronSchedule string        `json:"cronSchedule,omitempty"`
	Operations   []interface{} `json:"operations" binding:"required"`
	Description  string        `json:"description,omitempty"`
} //@name Task

type Variable struct {
	Type
	Description  string           `json:"desc,omitempty"`
	Display      string           `json:"display,omitempty"`
	Internal     bool             `json:"internal,omitempty"`
	Required     bool             `json:"required,omitempty"`
	Value        interface{}      `json:"value,omitempty"`
	UserEditable bool             `json:"userEdit,omitempty"`
	Options      []VariableOption `json:"options,omitempty"`
} //@name Variable

type VariableOption struct {
	Value   interface{} `json:"value"`
	Display string      `json:"display"`
} //@name VariableOption

type Execution struct {
	Command                 interface{}       `json:"command,omitempty"`
	StopCommand             string            `json:"stop,omitempty"`
	Disabled                bool              `json:"disabled,omitempty"`
	AutoStart               bool              `json:"autostart,omitempty"`
	AutoRestartFromCrash    bool              `json:"autorecover,omitempty"`
	AutoRestartFromGraceful bool              `json:"autorestart,omitempty"`
	PreExecution            []interface{}     `json:"pre,omitempty"`
	PostExecution           []interface{}     `json:"post,omitempty"`
	StopCode                int               `json:"stopCode,omitempty"`
	EnvironmentVariables    map[string]string `json:"environmentVars,omitempty"`
	LegacyRun               string            `json:"program,omitempty"`
	LegacyArguments         []string          `json:"arguments,omitempty"`
	WorkingDirectory        string            `json:"workingDirectory,omitempty"`
} //@name Execution

type Name struct {
	Name string `json:"name"`
} //@name Name

type Command struct {
	Command string `json:"command"`
	If      string `json:"if,omitempty"`
} //@name Command

type Type struct {
	Type string `json:"type"`
} //@name Type

type Requirements struct {
	Binaries []string `json:"binaries,omitempty"`
	OS       string   `json:"os,omitempty"`
	Arch     string   `json:"arch,omitempty"`
} //@name Requirements

type Group struct {
	Variables   []string `json:"variables"`
	Display     string   `json:"string"`
	Description string   `json:"description"`
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

func (s Server) DataToMap() map[string]interface{} {
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
