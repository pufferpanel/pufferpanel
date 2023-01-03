package steamgamedl

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type SteamMetadata struct {
	Version   string
	BinsLinux SteamFile
}

type SteamFile struct {
	File                  string
	Size                  string
	Sha2                  string
	ZipVz                 string
	Sha2Vz                string
	IsBootstrapperPackage string
}

func Parse(os string, body io.Reader) (string, error) {
	expectedRegex, err := regexp.CompilePOSIX("(\"steamcmd_bins_" + os + ".zip.[a-zA-Z0-9]+\")")
	if err != nil {
		return "", err
	}

	fileScanner := bufio.NewScanner(body)

	//first line indicates "key"
	var line string
	var results []string
	for fileScanner.Scan() {
		line = fileScanner.Text()
		results = expectedRegex.FindAllString(line, 1)
		if len(results) == 1 {
			return strings.Trim(results[0], "\""), nil
		}
	}

	return "", nil
}
