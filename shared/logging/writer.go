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

import "io"

type writeLogger struct {
	io.Writer

	level *Level
}

// This is simply a wrapper around the logger system so that interfaces which
// expect a Writer as a target will be able to use the logging structure we have
// located here.
func AsWriter(level *Level) io.Writer {
	return &writeLogger{level: level}
}

func (wl *writeLogger) Write(p []byte) (int, error) {
	Log(wl.level, string(p))
	return len(p), nil
}
