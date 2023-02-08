/*
 Copyright 2020 Padduck, LLC
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

package logging

import (
	"github.com/pufferpanel/pufferpanel/v3/config"
	"io"
	"log"
	"os"
	"path"
	"time"
)

var logFile io.WriteCloser
var flags = log.LstdFlags

var Error = log.New(os.Stderr, "[ERROR] ", flags)
var Debug = log.New(os.Stdout, "[DEBUG] ", flags)
var Info = log.New(os.Stdout, "[INFO] ", flags)

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

		logFile, err = os.OpenFile(path.Join(directory, time.Now().Format("2006-01-02T15-04-05")+".log"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
	}

	//just create them ourselves.....

	//first, create STDERR

	stderr := MultiWriter(logFile, os.Stderr, CreateServiceLogger("error"))
	Error = log.New(stderr, "[ERROR] ", flags)

	//now, STDOUT
	stdout := MultiWriter(logFile, os.Stdout, CreateServiceLogger("info"))
	Info = log.New(stdout, "[INFO] ", flags)

	//and now, a DEBUG
	stddebug := MultiWriter(logFile, os.Stdout)
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
	if logFile != nil {
		_ = logFile.Close()
	}
}
