package view

type OAuthTokenInfoViewModel struct {
	Active  bool                `json:"active"`
	Mapping map[string][]string `json:"mapping,omitempty"`
}
