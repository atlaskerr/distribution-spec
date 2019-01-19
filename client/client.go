package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// ErrUserAndToken is the error returned when username and token fields are
	// defined in the Config struct at the same time.
	ErrUserAndToken = errors.New("client cannot be configured with both user and token authentication")

	// ErrNoEndpoint is the error returned when no endpoint is defined in the
	// Config struct
	ErrNoEndpoint      = errors.New("no endpoint provided to client")
	ErrMissingUsername = errors.New("no username defined")
	ErrMissingPassword = errors.New("no password defined")
)

var DefaultRequestTimeout = 5 * time.Second

var DefaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
}

// Client is an interface that contains methods to interact with an registry
// compliant with the OCI Distribution Specification.
type Client interface {
	SetEndpoint(url string) error
	NewCredential(cred Credential)
	SetCredential(req *http.Request)
	//httpClient
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
func New(conf Config) (Client, error) {
	c := new(client)

	if conf.Endpoint == "" {
		return nil, ErrNoEndpoint
	}
	c.SetEndpoint(conf.Endpoint)

	basicAuth := (conf.Username != "") || (conf.Password != "")
	tokenAuth := conf.Token != ""
	switch {
	case basicAuth && tokenAuth:
		return nil, ErrUserAndToken
	case basicAuth:
		if conf.Username == "" {
			return nil, ErrMissingUsername
		}
		if conf.Password == "" {
			return nil, ErrMissingPassword
		}
		cred := &BasicCredential{
			Username: conf.Username,
			Password: conf.Password,
		}
		c.NewCredential(cred)
		c.authEnabled = true
	case tokenAuth:
		cred := &TokenCredential{
			Token: conf.Token,
		}
		c.NewCredential(cred)
		c.authEnabled = true
	}
	return c, nil
}

type httpClient interface {
	Do(context.Context, operation) (*http.Response, error)
}

type operation interface {
	HTTPRequest(url.URL) *http.Request
}

type client struct {
	endpoint    url.URL
	authEnabled bool
	credential  Credential
}

func (c *client) NewCredential(cred Credential) {
	c.credential = cred
}

func (c *client) SetCredential(req *http.Request) {
	c.credential.setCredential(req)
}

func (c *client) SetEndpoint(endpoint string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	c.endpoint = *u
	return nil
}

//func (c *client) Do(ctx context.Context, op operation) (resp *http.Response, err error) {
//}

// Credential defines a methods to inject credentials into an HTTP request.
type Credential interface {
	setCredential(*http.Request)
}

// TokenCredential implements the Credential interface. Uses an OAuth2 token to
// authenticate requests to the registry.
type TokenCredential struct {
	Token string
}

// SetCredential sets the authorization header of a request with an OAuth2 token.
func (c *TokenCredential) setCredential(req *http.Request) {
	bearer := strings.Join([]string{"Bearer", c.Token}, "")
	req.Header.Set("Authorization", bearer)
}

// BasicCredential implements the Credential interface. Uses username and
// password to authenticate requests to the registry.
type BasicCredential struct {
	Username string
	Password string
}

// SetCredential sets the authorization header of a request with a username and password.
func (c *BasicCredential) setCredential(req *http.Request) {
	req.SetBasicAuth(c.Username, c.Password)
}
