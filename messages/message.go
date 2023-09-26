package messages

type Message interface {
	Key() string
}

type Transmission struct {
	Message Message `json:"data"`
	Type    string  `json:"type"`
}
