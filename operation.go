package pufferpanel

type Operation interface {
	Run(args RunOperatorArgs) OperationResult
}

type OperationFactory interface {
	Create(CreateOperation) (Operation, error)

	Key() string
}

type CreateOperation struct {
	OperationArgs        map[string]interface{}
	EnvironmentVariables map[string]string
	DataMap              map[string]interface{}
}

type RunOperatorArgs struct {
	Environment Environment
	Server      DaemonServer
}

type OperationResult struct {
	Error             error
	VariableOverrides map[string]interface{}
}
