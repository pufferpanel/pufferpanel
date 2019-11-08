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

import (
	"os"
	"path"
	"time"
)

func WithLogDirectory(directory string, lvl *Level, ignore *Level) (err error) {
	if directory == "" {
		directory = "logs"
	}

	err = os.MkdirAll(directory, 0755)
	if err != nil && !os.IsExist(err) {
		return
	}

	file, err := os.OpenFile(path.Join(directory, time.Now().Format("2006-01-02T15-04-05.log")), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	WithWriterIgnore(file, lvl, ignore)
	return
}
