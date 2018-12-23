package error

type PanelError interface {
	error

	GetCode() int

	GetMachineMessage() string

	GetHumanMessage() string
}

type GenericError struct {
	PanelError

	code int
	machineMsg string
	humanMsg string
}

func FromError(err error) {

}