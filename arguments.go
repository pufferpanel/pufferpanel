package pufferpanel

import (
	"fmt"
	"strings"
)

func ReplaceTokens(msg string, mapping map[string]interface{}) string {
	newmsg := msg
	for key, value := range mapping {
		newmsg = strings.Replace(newmsg, "${"+key+"}", fmt.Sprint(value), -1)
	}
	return newmsg
}

func ReplaceTokensInArr(msg []string, mapping map[string]interface{}) []string {
	newarr := make([]string, len(msg))
	for index, element := range msg {
		newarr[index] = ReplaceTokens(element, mapping)
	}
	return newarr
}

func ReplaceTokensInMap(msg map[string]string, mapping map[string]interface{}) map[string]string {
	newarr := make(map[string]string, len(msg))
	for index, element := range msg {
		newarr[index] = ReplaceTokens(element, mapping)
	}
	return newarr
}

func SplitArguments(source string) (cmd string, arguments []string) {
	if source == "" {
		return "", []string{}
	}

	results := []string{""}

	skip := false //if this is set, the next char is always added to the current string
	inQuote := false
	for _, v := range source {
		if skip {
			skip = false
			results[len(results)-1] += string(v)
			continue
		}
		switch v {
		case '\\':
			{
				skip = true
			}
		case '"':
			{
				inQuote = !inQuote
				results[len(results)-1] += "\""
			}
		case ' ':
			{
				if inQuote {
					results[len(results)-1] += string(v)
				} else {
					results = append(results, "")
				}
			}
		default:
			results[len(results)-1] += string(v)
		}
	}

	//remove any "empty" items
	i := 0 // output index
	for _, x := range results {
		if x != "" {

			results[i] = x
			i++
		}
	}
	results = results[:i]

	cmd = results[0]
	arguments = results[1:]
	return
}
