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
	"github.com/pufferpanel/pufferpanel/v2/config"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

var logFile io.WriteCloser

var Error *log.Logger = log.New(os.Stderr, "[ERROR] ", log.Default().Flags())
var Debug *log.Logger = log.New(os.Stdout, "[DEBUG] ", log.Default().Flags())
var Info *log.Logger = log.New(os.Stdout, "[INFO] ", log.Default().Flags())

func Initialize(useFiles bool) {
	if useFiles {
		directory := config.GetString("logs")
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
	var serviceLogger io.Writer
	if runtime.GOOS == "windows" {
		serviceLogger = CreateServiceLogger("error")
	}

	flags := log.LstdFlags

	stderr := MultiWriter(logFile, os.Stderr, serviceLogger)
	Error = log.New(stderr, "[ERROR] ", flags)

	//now, STDOUT
	if runtime.GOOS == "windows" {
		serviceLogger = CreateServiceLogger("info")
	}
	stdout := MultiWriter(logFile, os.Stdout, serviceLogger)
	Info = log.New(stdout, "[INFO] ", flags)

	//and now, a DEBUG
	stddebug := MultiWriter(logFile, os.Stdout, serviceLogger)
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
