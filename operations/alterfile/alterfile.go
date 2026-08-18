package alterfile

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

type AlterFile struct {
	TargetFile string
	Search     string
	Replace    string
	Regex      bool
}

func (c AlterFile) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment

	fs := args.Server.GetFileServer()

	logging.Info.Printf("Changing data in file: %s", c.TargetFile)
	env.DisplayToConsole(true, "Changing some data in file: %s\n ", c.TargetFile)

	var regex *regexp.Regexp
	var err error
	if c.Regex {
		regex, err = regexp.Compile("(?m)" + c.Search)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	var file *os.File
	file, err = fs.OpenFile(c.TargetFile, os.O_RDWR, 0755)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	defer utils.Close(file)

	err = replaceLinesInFile(file, regex, []byte(c.Search), []byte(c.Replace))

	return pufferpanel.OperationResult{Error: err}
}

func replaceLinesInFile(file io.ReadWriter, regex *regexp.Regexp, needle []byte, replacement []byte) (err error) {
	s := bufio.NewScanner(file)

	for s.Scan() {
		line := []byte(s.Text())
		if regex != nil {
			line = regex.ReplaceAllLiteral(line, replacement)
		} else {
			line = bytes.ReplaceAll(line, needle, replacement)
		}
		_, err = file.Write(line)
		if err != nil {
			return
		}
		log.Println()
		_, err = file.Write([]byte{'\n'})
		if err != nil {
			return
		}
	}
	err = s.Err()
	return
}
