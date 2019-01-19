package client

import (
	"net/http"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	host := "http://localhost"

	tt := []struct {
		name string
		pass bool
		Config
	}{
		{"basic auth", true, Config{host, "user", "pass", ""}},
		{"basic auth no pass", false, Config{host, "user", "", ""}},
		{"basic auth no user", false, Config{host, "", "pass", ""}},
		{"token auth", true, Config{host, "", "", "token"}},
		{"basic and token auth", false, Config{host, "user", "pass", "token"}},
		{"no endpoint", false, Config{"", "", "", "token"}},
	}

	for _, tc := range tt {
		tf := func(t *testing.T) {
			_, err := New(tc.Config)
			if (err != nil) && (tc.pass) {
				t.Fatalf("expected valid config is invalid: %s", err)
			}
			if err == nil && (!tc.pass) {
				t.Fatalf("expected invalid config is valid")
			}
		}
		t.Run(tc.name, tf)
	}
}

func TestCredential(t *testing.T) {
	tt := []struct {
		name string
		Credential
	}{
		{"basic auth", &TokenCredential{"token"}},
		{"token auth", &BasicCredential{"user", "pass"}},
	}

	for _, tc := range tt {
		tf := func(t *testing.T) {
			host, _ := url.Parse("http://localhost")
			c := client{
				endpoint:    *host,
				authEnabled: true,
				credential:  tc.Credential,
			}
			req := new(http.Request)
			req.Header = make(http.Header)
			c.SetCredential(req)
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				t.Fatalf("auth token not set in request")
			}
		}
		t.Run(tc.name, tf)
	}

}
