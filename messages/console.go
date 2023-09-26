package messages

type Console struct {
	Logs []string `json:"logs"`
}

func (m Console) Key() string {
	return "console"
}
