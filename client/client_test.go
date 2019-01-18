package client

import (
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
