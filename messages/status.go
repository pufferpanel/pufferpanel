package messages

type Status struct {
	Running    bool `json:"running"`
	Installing bool `json:"installing"`
}

func (m Status) Key() string {
	return "status"
}
