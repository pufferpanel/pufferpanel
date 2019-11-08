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

type Builder interface {
	WithMessage(string) Builder
	WithError(error) Builder
	WithArgs(...interface{}) Builder
	Log()
}

type builder struct {
	format string
	args   []interface{}
	err    error
	lvl    *Level
}

func createBuilder(lvl *Level) Builder {
	if lvl == nil {
		lvl = ERROR
	}
	return &builder{
		lvl: lvl,
	}
}

func (b *builder) WithMessage(msg string) Builder {
	b.format = msg
	return b
}

func (b *builder) WithError(err error) Builder {
	b.err = err
	return b
}

func (b *builder) WithArgs(args ...interface{}) Builder {
	b.args = args
	return b
}

func (b *builder) Log() {
	//if we have an error, alter the format and add it to the args
	if b.err != nil {
		b.format = b.format + "\n%s"
		if b.args == nil {
			b.args = []interface{}{b.err}
		} else {
			b.args = append(b.args, b.err)
		}
	}

	Log(b.lvl, b.format, b.args...)
}
