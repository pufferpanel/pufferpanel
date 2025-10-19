package utils

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
		//older jvms start the lines we want with a space, newer ones don't
		line = strings.TrimPrefix(line, " ")

		if z, had := strings.CutPrefix(line, "garbage-first heap"); had {
			//heap could have array stuff in it, remove it
			results := parseLine(z)
			if num, exists := results["used"]; exists {
				stats.HeapUsed += num
			}
			if num, exists := results["total"]; exists {
				stats.HeapTotal += num
			}
			if num, exists := results["committed"]; exists {
				stats.HeapTotal += num
			}
		} else if z, had := strings.CutPrefix(line, "def new generation"); had {
			//heap could have array stuff in it, remove it
			results := parseLine(z)
			if num, exists := results["used"]; exists {
				stats.HeapUsed += num
			}
			if num, exists := results["total"]; exists {
				stats.HeapTotal += num
			}
		} else if z, had := strings.CutPrefix(line, "tenured generation"); had {
			//heap could have array stuff in it, remove it
			results := parseLine(z)
			if num, exists := results["used"]; exists {
				stats.HeapUsed += num
			}
			if num, exists := results["total"]; exists {
				stats.HeapTotal += num
			}
		} else if z, had = strings.CutPrefix(line, "Metaspace"); had {
			results := parseLine(z)
			if num, exists := results["used"]; exists {
				stats.MetaspaceUsed += num
			}
			if num, exists := results["committed"]; exists {
				stats.MetaspaceTotal += num
			}
		}
	}

	return stats
}

func parseLine(line string) map[string]int64 {
	result := make(map[string]int64)
	z := strings.Split(line, "[")[0]
	z = strings.TrimSpace(z)
	parts := strings.Split(z, ", ")
	for _, v := range parts {
		if strings.HasPrefix(v, "used ") {
			d := strings.TrimPrefix(v, "used ")
			d = strings.TrimSuffix(d, "K")
			result["used"] = cast.ToInt64(d) * 1024
		} else if strings.HasPrefix(v, "total ") {
			d := strings.TrimPrefix(v, "total ")
			d = strings.TrimSuffix(d, "K")
			result["total"] = cast.ToInt64(d) * 1024
		} else if strings.HasPrefix(v, "reserved ") {
			d := strings.TrimPrefix(v, "reserved ")
			d = strings.TrimSuffix(d, "K")
			result["reserved"] = cast.ToInt64(d) * 1024
		} else if strings.HasPrefix(v, "committed ") {
			d := strings.TrimPrefix(v, "committed ")
			d = strings.TrimSuffix(d, "K")
			result["committed"] = cast.ToInt64(d) * 1024
		}
	}
	return result
}

type JvmStats struct {
	HeapUsed       int64 `json:"heapUsed"`
	HeapTotal      int64 `json:"heapTotal"`
	MetaspaceUsed  int64 `json:"metaspaceUsed"`
	MetaspaceTotal int64 `json:"metaspaceTotal"`
}
