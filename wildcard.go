package pufferpanel

import (
	"regexp"
	"strings"
)

func CompareWildcard(source, match string) bool {
	if match == "" || match == "*" {
		return true
	}
	if strings.Contains(match, "*") {
		regex := WildCardToRegexp(match)
		return matchRegex(regex, source)
	} else {
		return source == match
	}
}

// WildCardToRegexp converts a wildcard pattern to a regular expression pattern.
func WildCardToRegexp(pattern string) string {
	var result strings.Builder
	for i, literal := range strings.Split(pattern, "*") {

		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return result.String()
}

func matchRegex(pattern string, value string) bool {
	result, _ := regexp.MatchString("^"+WildCardToRegexp(pattern)+"$", value)
	return result
}
