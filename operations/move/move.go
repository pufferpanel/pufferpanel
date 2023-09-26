package move

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"os"
	"path/filepath"
)

type Move struct {
	SourceFile string
	TargetFile string
}

func (m Move) Run(env pufferpanel.Environment) error {
	source := filepath.Join(env.GetRootDirectory(), m.SourceFile)
	target := filepath.Join(env.GetRootDirectory(), m.TargetFile)
	result, valid := validateMove(source, target)
	if !valid {
		return nil
	}

	for k, v := range result {
		logging.Info.Printf("Moving file from %s to %s", source, target)
		env.DisplayToConsole(true, "Moving file from %s to %s\n", m.SourceFile, m.TargetFile)
		err := os.Rename(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateMove(source string, target string) (result map[string]string, valid bool) {
	result = make(map[string]string)
	sourceFiles, _ := filepath.Glob(source)
	info, err := os.Stat(target)

	if err != nil {
		if os.IsNotExist(err) && len(sourceFiles) > 1 {
			logging.Error.Printf("Target folder does not exist")
			valid = false
			return
		} else if !os.IsNotExist(err) {
			valid = false
			logging.Error.Printf("Error reading target file: %s", err)
			return
		}
	} else if info.IsDir() && len(sourceFiles) > 1 {
		logging.Error.Printf("Cannot move multiple files to single file target")
		valid = false
		return
	}

	if info != nil && info.IsDir() {
		for _, v := range sourceFiles {
			_, fileName := filepath.Split(v)
			result[v] = filepath.Join(target, fileName)
		}
	} else {
		for _, v := range sourceFiles {
			result[v] = target
		}
	}
	valid = true
	return
}
