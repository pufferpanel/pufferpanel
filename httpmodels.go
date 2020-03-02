/*
 Copyright 2020 Padduck, LLC
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

type ServerIdResponse struct {
	Id string `json:"id"`
}

type ServerStats struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

type ServerLogs struct {
	Epoch int64  `json:"epoch"`
	Logs  string `json:"logs"`
}

type ServerRunning struct {
	Running bool `json:"running"`
}

type ServerData struct {
	Variables map[string]Variable `json:"data"`
}

type ServerDataAdmin struct {
	*Server
}

type PufferdRunning struct {
	Message string `json:"message"`
}
