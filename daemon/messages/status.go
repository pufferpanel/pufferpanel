package messages

type Status struct {
	Running bool `json:"running"`
}

func (m Status) Key() string {
	return "status"
}
