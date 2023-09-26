package messages

type Stat struct {
	Memory float64 `json:"memory"`
	Cpu    float64 `json:"cpu"`
}

func (m Stat) Key() string {
	return "stat"
}
