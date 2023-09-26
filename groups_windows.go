//go:build !linux

package pufferpanel

func UserInGroup(groups ...string) bool {
	return true
}
