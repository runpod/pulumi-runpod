//go:generate genqlient genqlient.yaml
//go:generate go run ../genqlient_fixup.go
package runpod

import (
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
)

const (
	DefaultAPIURL = "https://api.runpod.io/graphql"
	DevAPIURL     = "https://api.runpod.dev/graphql"
)

// authTransport adds Bearer token authentication to HTTP requests.
type authTransport struct {
	apiKey    string
	transport http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	return t.transport.RoundTrip(req)
}

// NewClient creates a genqlient GraphQL client with Bearer token auth.
func NewClient(apiKey, apiURL string) graphql.Client {
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &authTransport{
			apiKey:    apiKey,
			transport: http.DefaultTransport,
		},
	}
	return graphql.NewClient(apiURL, httpClient)
}
