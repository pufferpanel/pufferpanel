package utils

import (
	"github.com/pufferpanel/apufferi/v4"
	"os/user"
)

func UserInGroup() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}

	expectedGroup, err := user.LookupGroup("pufferpanel")
	if err != nil {
		return false
	}
	groups, err := u.GroupIds()
	if err != nil {
		return false
	}

	return apufferi.ContainsString(groups, expectedGroup.Gid)
}
