package pufferpanel

import "github.com/pufferpanel/pufferpanel/v3/config"

var useOpenat2 = false

func DetermineKernelSupport() {
	if config.SecurityForceOpenat2.Value() {
		useOpenat2 = true
	} else if config.SecurityForceOpenat.Value() {
		useOpenat2 = false
	} else {
		testOpenat2()
	}
}

func UseOpenat2() bool {
	return useOpenat2
}
