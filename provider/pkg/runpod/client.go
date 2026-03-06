// Copyright 2025, RunPod, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package runpod provides a GraphQL client for the RunPod API.
//
//go:generate genqlient genqlient.yaml
//go:generate go run ../genqlient_fixup.go
package runpod

import (
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
)

// DefaultAPIURL is the production RunPod GraphQL endpoint.
const DefaultAPIURL = "https://api.runpod.io/graphql"

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
