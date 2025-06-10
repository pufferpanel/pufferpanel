package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/sys"
	"golang.org/x/sys/unix"
	"os"
	"strings"
)

func testOpenat2() {
	f, err := os.Open("/proc/kallsyms")
	if errors.Is(err, os.ErrNotExist) {
		logging.Info.Printf("Could not open /proc/kallsyms to validate kernel support, falling back to temp file test\n%s", err.Error())

		var testPath string
		testPath, err = os.MkdirTemp(os.TempDir(), "puffer-openat2-test-*")
		if err != nil {
			panic(fmt.Errorf("failed to validate kernel support with test file\n%s", err.Error()))
		}
		defer func(tar string) {
			_ = os.Remove(tar)
		}(testPath)

		var testFile *os.File
		testFile, err = os.Open(testPath)
		if err != nil {
			panic(fmt.Errorf("failed to validate kernel support with test file\n%s", err.Error()))
		}
		defer Close(testFile)

		//we have a file now, let's see if we can... read it with openat2
		var fd int
		fd, err = unix.Openat2(int(testFile.Fd()), "validate", &unix.OpenHow{
			Flags: uint64(os.O_CREATE),
			Mode:  uint64(sys.SyscallMode(0644)),
		})
		if err == nil {
			_ = unix.Close(fd)
			useOpenat2 = true
		} else if errors.Is(err, unix.EOPNOTSUPP) || errors.Is(err, unix.ENOSYS) {
			useOpenat2 = false
		} else {
			panic(fmt.Errorf("failed to validate kernel support with test file\n%s", err.Error()))
		}
	} else if err == nil {
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
	} else {
		panic(fmt.Errorf("Could not open /proc/kallsyms to validate kernel support\n%s", err.Error()))
	}
}
