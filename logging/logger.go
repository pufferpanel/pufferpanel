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
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var mapper = make(map[string]*log.Logger)

var lock sync.Mutex

var logFile *os.File

func Initialize() {
	_ = Info()
	_ = Debug()
	_ = Error()
}

func Error() *log.Logger {
	return Get("ERROR")
}

func Debug() *log.Logger {
	return Get("DEBUG")
}

func Info() *log.Logger {
	return Get("INFO")
}

func Close() {
	if logFile != nil {
		_ = logFile.Close()
	}
}

func AsWriter() io.Writer {
	return nil
}

func create(prefix string) *log.Logger {
	if logFile == nil {
		directory := viper.GetString("logs")
		if directory == "" {
			directory = "logs"
		}

		err := os.MkdirAll(directory, 0755)
		if err != nil && !os.IsExist(err) {

		}

		logFile, err = os.OpenFile(path.Join(directory, time.Now().Format("2006-01-02T15-04-05.log")), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
		}
	}

	var writer io.Writer
	if prefix == "ERROR" || prefix == "STDERR" {
		writer = os.Stderr
	} else {
		writer = os.Stdout
	}

	if logFile != nil {
		writer = io.MultiWriter(writer, logFile)
	}

	if prefix == "INFO" {
		log.SetOutput(writer)
	}

	l := log.New(writer, fmt.Sprintf("[%s] ", prefix), 0)
	mapper[prefix] = l
	return l
}

func Get(name string) *log.Logger {
	name = strings.ToUpper(name)
	l := mapper[name]
	if l == nil {
		lock.Lock()
		defer lock.Unlock()

		l = create(name)
	}

	return l
}
