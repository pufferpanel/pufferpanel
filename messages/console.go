package messages

type Console struct {
	Logs []byte `json:"logs"`
}

func (m Console) Key() string {
	return "console"
}
