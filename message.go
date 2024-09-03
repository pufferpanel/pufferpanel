package pufferpanel

type Transmission struct {
	Message interface{}      `json:"data"`
	Type    TransmissionType `json:"type"`
}

type TransmissionType string

const (
	MessageTypeLog    = "console"
	MessageTypeStats  = "stat"
	MessageTypeStatus = "status"
)
