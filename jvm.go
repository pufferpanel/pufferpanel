package pufferpanel

import (
	"bufio"
	"bytes"
	"github.com/spf13/cast"
	"strings"
	"unicode"
)

// ParseJCMDResponse Parses the results of a jcmd command that is executed
// As this data is retrieved from several places, this is the wrapper to
// decode the data
func ParseJCMDResponse(data []byte) *JvmStats {
	//data is as such:
	//line 1 is the PID
	//"garbage-first heap" is the actual heap results
	//"Metadata" is metaspace usage
	scanner := bufio.NewScanner(bytes.NewReader(data))

	stats := &JvmStats{}

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.Map(func(r rune) rune {
			if unicode.IsGraphic(r) {
				return r
			}
			return -1
		}, line)

		if z, had := strings.CutPrefix(line, " garbage-first heap"); had {
			//heap could have array stuff in it, remove it
			z = strings.Split(z, "[")[0]
			z = strings.TrimSpace(z)
			parts := strings.Split(z, ", ")
			for _, v := range parts {
				if strings.HasPrefix(v, "used ") {
					d := strings.TrimPrefix(v, "used ")
					d = strings.TrimSuffix(d, "K")
					stats.HeapUsed = cast.ToInt64(d)
					stats.HeapUsed *= 1024
				} else if strings.HasPrefix(v, "total ") {
					d := strings.TrimPrefix(v, "total ")
					d = strings.TrimSuffix(d, "K")
					stats.HeapTotal = cast.ToInt64(d)
					stats.HeapTotal *= 1024
				}
			}
		} else if z, had = strings.CutPrefix(line, " Metaspace"); had {
			line = strings.TrimSpace(z)
			parts := strings.Split(line, ", ")
			for _, v := range parts {
				if d, o := strings.CutPrefix(v, "used "); o {
					d = strings.TrimSuffix(d, "K")
					stats.MetaspaceUsed = cast.ToInt64(d)
					stats.MetaspaceUsed *= 1024
				}
			}
		}
	}

	return stats
}
