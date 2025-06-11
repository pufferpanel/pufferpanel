package pufferpanel

type EnvironmentFactory interface {
	Create() EnvironmentImpl

	Key() string
}
