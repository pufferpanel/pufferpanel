package models

type Deployment struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	PublicKey    string `json:"publicKey"`
}
