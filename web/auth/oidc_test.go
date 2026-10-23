package auth

import (
	"encoding/base64"
	"encoding/json"
	"slices"
	"testing"
)

func TestOidcScopes(t *testing.T) {
	for value, want := range map[string][]string{
		"openid email":         {"openid", "email"},
		"email profile":        {"openid", "email", "profile"},
		"  email   openid  ":   {"email", "openid"},
		"":                     {"openid"},
		"openid email profile": {"openid", "email", "profile"},
	} {
		if got := oidcScopes(value); !slices.Equal(got, want) {
			t.Errorf("oidcScopes(%q) = %v, want %v", value, got, want)
		}
	}
}

func TestDecodeOidcState(t *testing.T) {
	encode := func(s oidcState) string {
		data, _ := json.Marshal(s)
		return base64.RawURLEncoding.EncodeToString(data)
	}

	state, err := decodeOidcState(encode(oidcState{State: "s", Nonce: "n", Verifier: "v"}))
	if err != nil || state.State != "s" || state.Nonce != "n" || state.Verifier != "v" {
		t.Fatalf("round trip failed: %v %v", state, err)
	}

	for name, raw := range map[string]string{
		"not base64":    "%%%",
		"not json":      base64.RawURLEncoding.EncodeToString([]byte("nope")),
		"missing nonce": encode(oidcState{State: "s", Verifier: "v"}),
		"missing state": encode(oidcState{Nonce: "n", Verifier: "v"}),
		"missing pkce":  encode(oidcState{State: "s", Nonce: "n"}),
		"empty cookie":  "",
	} {
		if _, err := decodeOidcState(raw); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestRandomOidcValue(t *testing.T) {
	a, b := randomOidcValue(), randomOidcValue()
	if a == b || len(a) < 40 {
		t.Fatalf("values must be long and unique: %q %q", a, b)
	}
}
