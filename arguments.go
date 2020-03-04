/*
 Copyright 2019 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

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
	results := []string{""}

	skip := false //if this is set, the next char is always added to the current string
	inQuote := false
	for _, v := range source {
		if skip {
			skip = false
			results[len(results) - 1] += string(v)
		} else if v == '\\' {
			skip = true
		} else if v == '"' {
			inQuote = !inQuote
			results[len(results) - 1] += "\""
		} else if v == ' ' && !inQuote {
			results = append(results, "")
		} else {
			results[len(results) - 1] += string(v)
		}
	}

	if results[len(results) - 1] == "" {
		results = results[:len(results)-1]
	}

	cmd = results[0]
	arguments = results[1:]
	return
}
