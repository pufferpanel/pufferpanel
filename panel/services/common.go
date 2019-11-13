package services

import "strings"

func ParseAllowedTags(source string, allowed []string) []string {
	includeTags := make([]string, 0)

	args := strings.Split(strings.ToLower(source), ",")

	for _, test := range args {
		for _, v := range allowed {
			if test == v {
				includeTags = append(includeTags, v)
			}
		}
	}

	return includeTags
}
