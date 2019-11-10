package pufferpanel

import (
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

	return ContainsString(groups, expectedGroup.Gid)
}
