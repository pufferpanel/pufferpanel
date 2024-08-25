package email

import "strings"

type Provider interface {
	Send(to, subject, body string) error
}

var providers = make(map[string]Provider)

func GetProvider(provider string) Provider {
	return providers[strings.ToLower(provider)]
}
