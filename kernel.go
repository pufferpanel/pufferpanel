package pufferpanel

import (
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/logging"
)

var useOpenat2 = false

func DetermineKernelSupport() {
	if config.SecurityForceOpenat2.Value() {
		useOpenat2 = true
	} else if config.SecurityForceOpenat.Value() {
		useOpenat2 = false
	} else {
		testOpenat2()
	}

	if useOpenat2 {
		logging.Debug.Printf("openat2 enabled")
	} else {
		logging.Info.Printf("WARNING: OPENAT2 SUPPORT NOT ENABLED")
	}
}

func UseOpenat2() bool {
	return useOpenat2
}
