package main

import (
	"flag"
	"strings"
)

var CmdFlags = &Flags{}

func init() {
	var skip string
	var required string
	var files string

	flag.StringVar(&CmdFlags.GitRef, "gitRef", "refs/heads/v3", "")
	flag.StringVar(&skip, "skip", "", "")
	flag.StringVar(&required, "require", "*", "")
	flag.StringVar(&CmdFlags.TemplateFolder, "path", "", "")
	flag.BoolVar(&CmdFlags.DeleteTemp, "delete", true, "")
	flag.StringVar(&CmdFlags.WorkingDir, "workDir", "", "")
	flag.BoolVar(&CmdFlags.Reuse, "reuse", false, "")
	flag.StringVar(&files, "files", "", "")
	flag.StringVar(&CmdFlags.PufferpanelBinary, "binary", "pufferpanel", "")
	flag.StringVar(&CmdFlags.Host, "host", "http://127.0.0.1:8080", "")
	flag.Parse()

	if skip != "" {
		CmdFlags.Skip = strings.Split(skip, ",")
	} else {
		CmdFlags.Skip = make([]string, 0)
	}

	if skip != "" {
		CmdFlags.Required = strings.Split(required, ",")
	} else {
		CmdFlags.Required = make([]string, 0)
	}

	if skip != "" {
		CmdFlags.Files = strings.Split(files, ",")
	} else {
		CmdFlags.Files = make([]string, 0)
	}
}

type Flags struct {
	GitRef            string
	Skip              []string
	Required          []string
	TemplateFolder    string
	DeleteTemp        bool
	WorkingDir        string
	Reuse             bool
	Files             []string
	PufferpanelBinary string
	Host              string
}
