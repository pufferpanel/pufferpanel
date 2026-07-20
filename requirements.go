package pufferpanel

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

type Requirements struct {
	OS       string   `json:"os,omitempty"`
	Arch     string   `json:"arch,omitempty"`
	Binaries []string `json:"binaries,omitempty"`
} //@name Requirements

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
	err := utils.UnmarshalTo(server.Environment, &envType)
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
				parsed := utils.ReplaceTokens(binary, server.DataToMap(), utils.PlainReplace)
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
