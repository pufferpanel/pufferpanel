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

import "time"

type Cache interface {
	Read() (cache []string, epoch int64)

	ReadFrom(startTime int64) (cache []string, epoch int64)

	Write(b []byte) (n int, err error)
}

type Message struct {
	msg  string
	time int64
}

type MemoryCache struct {
	Cache
	Buffer   []Message
	Capacity int
}

func (c *MemoryCache) Read() (msg []string, lastTime int64) {
	msg, lastTime = c.ReadFrom(0)
	return
}

func (c *MemoryCache) ReadFrom(startTime int64) (msg []string, lastTime int64) {
	result := make([]string, 0)

	var endTime int64 = 0

	for _, v := range c.Buffer {
		if v.time > startTime {
			result = append(result, v.msg)
			endTime = v.time
		}
	}

	if endTime == 0 {
		endTime = time.Now().Unix()
	}
	return result, endTime
}

func (c *MemoryCache) Write(b []byte) (n int, err error) {
	if len(c.Buffer) == c.Capacity {
		c.Buffer = c.Buffer[1:]
	}
	c.Buffer = append(c.Buffer, Message{msg: string(b), time: time.Now().Unix()})
	n = len(b)
	return
}
