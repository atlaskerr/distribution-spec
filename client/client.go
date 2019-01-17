package client

import (
	"errors"
	"net/http"
	"strings"
)

var (
	// ErrUserAndToken is the error returned when username and token fields are
	// defined in the Config struct at the same time.
	ErrUserAndToken = errors.New("client cannot be configured with both user and token authentication")

	// ErrNoEndpoint is the error returned when no endpoint is defined in the
	// Config struct
	ErrNoEndpoint = errors.New("no endpoint provided to client")
)

// Client is an interface that contains methods to interact with an registry
// compliant with the OCI Distribution Specification.
type Client interface {
	SetEndpoint(url string)
	SetCredential(cred Credential)
}

// Config defines configuration parameters for the client.
type Config struct {

	// Endpoint is the URL of the registry.
	Endpoint string

	Username string
	Password string
	Token    string
}

// New returns a new Client.
func New(conf *Config) (Client, error) {
	var c *client

	if &conf.Endpoint == nil {
		return nil, ErrNoEndpoint
	}
	c.SetEndpoint(conf.Endpoint)

	basicAuth := (&conf.Username != nil) && (&conf.Password != nil)
	tokenAuth := &conf.Token != nil
	switch {
	case basicAuth && tokenAuth:
		return nil, ErrUserAndToken
	case basicAuth:
		cred := &BasicCredential{}
		c.SetCredential(cred)
	case tokenAuth:
		cred := &TokenCredential{}
		c.SetCredential(cred)
	}
	return c, nil
}

type client struct {
	endpoint   string
	credential *Credential
}

func (c *client) SetCredential(cred Credential) {
	c.credential = &cred
}

func (c *client) SetEndpoint(url string) {
	c.endpoint = url
}

// Credential defines a methods to inject credentials into an HTTP request.
type Credential interface {
	Set(*http.Request)
}

// TokenCredential implements the Credential interface. Uses an OAuth2 token to
// authenticate requests to the registry.
type TokenCredential struct {
	Token string
}

// Set sets the authorization header of a request with an OAuth2 token.
func (c *TokenCredential) Set(req *http.Request) {
	bearer := strings.Join([]string{"Bearer", c.Token}, "")
	req.Header.Set("Authorization", bearer)
}

// BasicCredential implements the Credential interface. Uses username and
// password to authenticate requests to the registry.
type BasicCredential struct {
	Username string
	Password string
}

// Set sets the authorization header of a request with a username and password.
func (c *BasicCredential) Set(req *http.Request) {
	req.SetBasicAuth(c.Username, c.Password)
}
