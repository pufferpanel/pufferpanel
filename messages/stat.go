package messages

type Stat struct {
	Cpu    float64   `json:"cpu"`
	Memory float64   `json:"memory"`
	Jvm    *JvmStats `json:"jvm,omitempty"`
}

type JvmStats struct {
	HeapUsed      int64 `json:"heapUsed"`
	HeapTotal     int64 `json:"heapTotal"`
	MetaspaceUsed int64 `json:"metaspaceUsed"`
}

func (m Stat) Key() string {
	return "stat"
}
