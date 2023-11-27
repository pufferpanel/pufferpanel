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
	Logs  []byte `json:"logs"`
} //@name ServerLogs

type ServerRunning struct {
	Running    bool `json:"running"`
	Installing bool `json:"installing"`
} //@name ServerRunning

type ServerData struct {
	Variables map[string]Variable `json:"data"`
	Groups    []Group             `json:"groups,omitempty"`
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

type ServerFlags struct {
	AutoStart             *bool `json:"autoStart,omitempty"`
	AutoRestartOnCrash    *bool `json:"autoRestartOnCrash,omitempty"`
	AutoRestartOnGraceful *bool `json:"autoRestartOnGraceful,omitempty"`
} //@name ServerFlags
