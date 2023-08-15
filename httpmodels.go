/*
 Copyright 2022 (c) PufferPanel
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
} //@name ServerId

type ServerStats struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
} //@name ServerStats

type ServerLogs struct {
	Epoch int64  `json:"epoch"`
	Logs  string `json:"logs"`
} //@name ServerLogs

type ServerRunning struct {
	Running bool `json:"running"`
} //@name ServerRunning

type ServerData struct {
	Variables map[string]Variable `json:"data"`
} //@name ServerData

type ServerDataAdmin struct {
	*Server
}

type DaemonRunning struct {
	Message string `json:"message"`
} //@name DaemonRunning

type ServerTasks struct {
	Tasks map[string]ServerTask
} //@name ServerTasks

type ServerTask struct {
	IsRunning bool `json:"isRunning"`
	Task
} //@name ServerTask

type ErrorResponse struct {
	Error *Error `json:"error"`
} //@name ErrorResponse

type Metadata struct {
	Paging *Paging `json:"paging"`
} //@name Metadata

type Paging struct {
	Page    uint  `json:"page"`
	Size    uint  `json:"pageSize"`
	MaxSize uint  `json:"maxSize"`
	Total   int64 `json:"total"`
} //@name Paging
