package pufferpanel

import (
	"os/user"
)

func UserInGroup() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}

	groups, err := u.GroupIds()
	if err != nil {
		return false
	}

	allowedIds := make([]string, 0)

	if expectedGroup, err := user.LookupGroup("pufferpanel"); err == nil {
		allowedIds = append(allowedIds, expectedGroup.Gid)
	}

	if rootGroup, err := user.LookupGroup("root"); err == nil {
		allowedIds = append(allowedIds, rootGroup.Gid)
	}

	for _, v := range groups {
		for _, t := range allowedIds {
			if v == t {
				return true
			}
		}
	}

	return false
}
