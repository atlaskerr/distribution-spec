package client

import (
	"net/http"
	"strings"

	ispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Client is an interface that contains methods to interact with an registry
// compliant with the OCI Distribution Specification.
type Client interface {
	Verify() bool

	// GetManifest retrieves a manifest from a repository.
	GetManifest(repo string, ref string) (ispec.Manifest, error)

	// VerifyManifest checks to see if a manifest exists in a repository.
	VerifyManifest(repo string, ref string) (bool, error)

	// UploadManifest uploads a manifest to a repository.
	UploadManifest(repo string, ref string, manifest ispec.Manifest) error

	// DeleteManifest deletes a manifest from a repository.
	DeleteManifest(repo string, ref string) error
}

// Config defines configuration parameters for the client.
type Config struct {

	// Endpoint is the URL of the registry.
	Endpoint string

	Credential Credential
}

// New returns a new Client.
func New(conf Config) (Client, error) {
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
