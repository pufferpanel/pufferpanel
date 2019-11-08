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

package logging

import "strings"

type Level struct {
	scale   byte
	display string
}

func (l Level) GetScale() byte {
	return l.scale
}

func (l Level) GetName() string {
	return l.display
}

var (
	DEBUG    = &Level{scale: 7, display: "DEBUG"}
	INFO     = &Level{scale: 31, display: "INFO"}
	WARN     = &Level{scale: 63, display: "WARN"}
	ERROR    = &Level{scale: 127, display: "ERROR"}
	CRITICAL = &Level{scale: 255, display: "CRITICAL"}
	DEVEL    = &Level{scale: 0, display: "DEVEL"}

	levels = []*Level{
		DEBUG, INFO, WARN, ERROR, CRITICAL, DEVEL,
	}
)

func GetLevel(name string) *Level {
	name = strings.ToUpper(name)
	for _, v := range levels {
		if v.GetName() == name {
			return v
		}
	}

	return nil
}
