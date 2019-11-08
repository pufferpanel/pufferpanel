package utils

import (
	"github.com/pufferpanel/pufferpanel/v2/shared"
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

	return shared.ContainsString(groups, expectedGroup.Gid)
}
