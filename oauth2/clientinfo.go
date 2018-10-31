package oauth2

type ClientInfo struct {
}

func (ci *ClientInfo) GetSecret() string {
	return "test"
}

func (ci *ClientInfo) GetID() string {
	return "test"
}

func (ci *ClientInfo) GetDomain() string {
	return "*"
}

func (ci *ClientInfo) GetUserID() string {
	return "test"
}