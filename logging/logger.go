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

func init() {
	_ = Info()
	_ = Debug()
	_ = Error()
}

func Error() *log.Logger {
	return get("ERROR")
}

func Debug() *log.Logger {
	return get("DEBUG")
}

func Info() *log.Logger {
	return get("INFO")
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

func get(name string) *log.Logger {
	name = strings.ToUpper(name)
	l := mapper[name]
	if l == nil {
		lock.Lock()
		defer lock.Unlock()

		l = create(name)
	}

	return l
}
