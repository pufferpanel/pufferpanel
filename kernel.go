package pufferpanel

import (
	"bufio"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"os"
	"runtime"
	"strings"
)

var useOpenat2 = false

func DetermineKernelSupport() {
	if runtime.GOOS == "linux" {
		f, err := os.Open("/proc/kallsyms")
		if err != nil {
			panic("could not open /proc/kallsyms to validate kernel support, hard failing.\n" + err.Error())
		}
		defer Close(f)
		reader := bufio.NewScanner(f)
		var line string
		for reader.Scan() {
			line = reader.Text()
			if strings.Contains(line, " t do_sys_openat2") {
				useOpenat2 = true
				break
			}
		}
		if !useOpenat2 {
			logging.Info.Printf("WARNING: OPENAT2 SUPPORT NOT ENABLED")
		}
	}
}

func UseOpenat2() bool {
	return useOpenat2
}
