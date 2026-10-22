package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"

	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
)

type fakeUserLookup struct {
	users []*models.User
}

func (f *fakeUserLookup) Get(username string) (*models.User, error) {
	for _, u := range f.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeUserLookup) GetByEmail(email string) (*models.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func TestMatchOidcUser(t *testing.T) {
	lookup := &fakeUserLookup{users: []*models.User{
		{ID: 1, Username: "alice", Email: "alice@example.com"},
		{ID: 2, Username: "bob", Email: "bob@example.com"},
	}}

	tests := []struct {
		name            string
		claims          map[string]interface{}
		mode            string
		requireVerified bool
		wantId          uint
		wantErr         error
	}{
		{"email verified", map[string]interface{}{"email": "alice@example.com", "email_verified": true}, "email", true, 1, nil},
		{"email verified as string", map[string]interface{}{"email": "alice@example.com", "email_verified": "true"}, "email", true, 1, nil},
		{"email not verified", map[string]interface{}{"email": "alice@example.com", "email_verified": false}, "email", true, 0, errOidcEmailNotVerified},
		{"email verification missing", map[string]interface{}{"email": "alice@example.com"}, "email", true, 0, errOidcEmailNotVerified},
		{"email trusted without verification", map[string]interface{}{"email": "alice@example.com"}, "email", false, 1, nil},
		{"email unknown", map[string]interface{}{"email": "eve@example.com", "email_verified": true}, "email", true, 0, errOidcNoUser},
		{"email claim missing", map[string]interface{}{"preferred_username": "alice"}, "email", true, 0, errOidcNoUser},
		{"username", map[string]interface{}{"preferred_username": "bob"}, "username", true, 2, nil},
		{"username mode is case-insensitive", map[string]interface{}{"preferred_username": "bob"}, "Username", true, 2, nil},
		{"username ignores email", map[string]interface{}{"email": "alice@example.com", "email_verified": true, "preferred_username": "bob"}, "username", true, 2, nil},
		{"username unknown", map[string]interface{}{"preferred_username": "eve"}, "username", true, 0, errOidcNoUser},
		{"username claim missing", map[string]interface{}{"email": "bob@example.com"}, "username", true, 0, errOidcNoUser},
		{"username claim not a string", map[string]interface{}{"preferred_username": 42}, "username", true, 0, errOidcNoUser},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := matchOidcUser(lookup, tt.claims, tt.mode, "preferred_username", tt.requireVerified)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && (user == nil || user.ID != tt.wantId) {
				t.Fatalf("user = %v, want id %d", user, tt.wantId)
			}
		})
	}
}

func TestMatchOidcUserLookupError(t *testing.T) {
	lookup := &failingLookup{}
	_, err := matchOidcUser(lookup, map[string]interface{}{"email": "a@example.com", "email_verified": true}, "email", "preferred_username", true)
	if err == nil || errors.Is(err, errOidcNoUser) {
		t.Fatalf("a database error must not look like an unknown user, got %v", err)
	}
}

type failingLookup struct{}

func (f *failingLookup) Get(string) (*models.User, error) { return nil, errors.New("db down") }

func (f *failingLookup) GetByEmail(string) (*models.User, error) { return nil, errors.New("db down") }

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
