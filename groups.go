//go:build !windows

package pufferpanel

import (
	"fmt"
	"os/user"
)

func UserInGroup(groups ...string) bool {
	//add root as an allowed group
	groups = append(groups, "root")

	u, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	allowedIds := make([]string, 0)
	for _, v := range groups {
		rootGroup, err := user.LookupGroup(v)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			allowedIds = append(allowedIds, rootGroup.Gid)
		}
	}

	g, err := u.GroupIds()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for _, v := range g {
		for _, t := range allowedIds {
			if v == t {
				return true
			}
		}
	}

	return false
}
