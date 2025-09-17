package utils

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

var useOpenat2 = false

func DetermineKernelSupport() {
	if config.SecurityForceOpenat.Value() {
		useOpenat2 = false
		logging.Info.Printf("WARNING: OPENAT2 SUPPORT NOT ENABLED. OVERRIDE IS BEING USED.")
	} else {
		passes := testOpenat2()
		if !passes {
			panic(fmt.Sprintf("OpenAt2 is not supported. Cowardly not starting to avoid security issues."))
		}
	}
}

func UseOpenat2() bool {
	return useOpenat2
}
