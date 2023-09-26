package pufferpanel

type EnvironmentFactory interface {
	Create(id string) Environment

	Key() string
}
