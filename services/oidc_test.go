package services

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

func TestOidcStateRoundTrip(t *testing.T) {
	state, err := NewOidcState(true)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := state.Encode()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeOidcState(encoded)
	if err != nil || decoded != state {
		t.Fatalf("round trip failed: %v %v", decoded, err)
	}
	if !decoded.Matches(state.State) || decoded.Matches(state.State+"x") {
		t.Fatal("Matches must only accept the original state")
	}
}

func TestDecodeOidcStateRejectsBrokenCookies(t *testing.T) {
	encode := func(s OidcState) string {
		data, _ := json.Marshal(s)
		return base64.URLEncoding.EncodeToString(data)
	}

	for name, raw := range map[string]string{
		"not base64":    "%%%",
		"not json":      base64.URLEncoding.EncodeToString([]byte("nope")),
		"missing nonce": encode(OidcState{State: "s", Verifier: "v"}),
		"missing state": encode(OidcState{Nonce: "n", Verifier: "v"}),
		"missing pkce":  encode(OidcState{State: "s", Nonce: "n"}),
		"empty cookie":  "",
	} {
		if _, err := DecodeOidcState(raw); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}
