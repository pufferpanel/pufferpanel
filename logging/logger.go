package logging

import (
	"github.com/pufferpanel/pufferpanel/v3/config"
	"io"
	"log"
	"os"
	"path"
)

var rotation *Rotator
var flags = log.LstdFlags

var Error = log.New(os.Stderr, "[ERROR] ", flags)
var Debug = log.New(os.Stdout, "[DEBUG] ", flags)
var Info = log.New(os.Stdout, "[INFO] ", flags)
var Server = log.New(os.Stdout, "[SERVER] ", flags)
var OriginalStdOut = os.Stdout

func Initialize(useFiles bool) {
	if useFiles {
		directory := config.LogsFolder.Value()
		if directory == "" {
			directory = "logs"
		}

		err := os.MkdirAll(directory, 0755)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}

		logFile, err := os.OpenFile(path.Join(directory, "pufferpanel.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		rotation = &Rotator{backer: logFile}
		rotation.StartRotation(directory)
	}

	//just create them ourselves.....

	//first, create STDERR

	stderr := MultiWriter(rotation, os.Stderr, CreateServiceLogger("error"))
	Error = log.New(stderr, "[ERROR] ", flags)

	//now, STDOUT
	stdout := MultiWriter(rotation, os.Stdout, CreateServiceLogger("info"))
	Info = log.New(stdout, "[INFO] ", flags)

	//and now, a DEBUG
	stddebug := MultiWriter(rotation, os.Stdout)
	Debug = log.New(stddebug, "[DEBUG] ", flags)

	log.SetOutput(Info.Writer())

	//i hate go's idea of how stdout should work
	outR, outW, _ := os.Pipe()
	errR, errW, _ := os.Pipe()

	go func() {
		_, _ = io.Copy(Info.Writer(), outR)
	}()

	go func() {
		_, _ = io.Copy(Error.Writer(), errR)
	}()

	os.Stdout = outW
	os.Stderr = errW
}

func Close() {
	if rotation != nil {
		_ = rotation.Close()
	}
}
