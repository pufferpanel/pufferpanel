package pufferpanel

import (
	"regexp"
	"strings"
)

var ipv6Regex = regexp.MustCompile(`(?P<host>\\[[a-f0-9:%]+\\])(:[0-9]+)?`)

func GetHostname(requestHost string) string {
	//ipv6
	if strings.HasPrefix(requestHost, "[") {
		return ipv6Regex.FindStringSubmatch(requestHost)[1]
	} else {
		return strings.SplitN(requestHost, ":", 2)[0]
	}
}
