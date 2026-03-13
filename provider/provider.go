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

// Package provider implements the RunPod Pulumi provider.
package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/Khan/genqlient/graphql"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	"github.com/runpod/pulumi-runpod/provider/pkg/runpod"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

// Name controls how this provider is referenced in package names and elsewhere.
const Name string = "runpod"

// Provider creates a new instance of the RunPod provider.
func Provider() p.Provider {
	prov, err := infer.NewProviderBuilder().
		WithDisplayName("RunPod").
		WithDescription("Manage RunPod GPU cloud resources.").
		WithHomepage("https://www.runpod.io").
		WithNamespace("runpod").
		WithPublisher("RunPod").
		WithKeywords("category/infrastructure", "kind/native").
		WithLicense("Apache-2.0").
		WithPluginDownloadURL("github://api.github.com/runpod/pulumi-runpod").
		WithRepository("https://github.com/runpod/pulumi-runpod").
		WithLanguageMap(map[string]any{
			"go": map[string]any{
				"importBasePath":                 "github.com/runpod/pulumi-runpod/sdk/go/runpod",
				"generateResourceContainerTypes": true,
				"respectSchemaVersion":           true,
			},
			"nodejs": map[string]any{
				"packageName":          "pulumi-runpod",
				"respectSchemaVersion": true,
			},
			"python": map[string]any{
				"packageName": "pulumi_runpod",
				"pyproject": map[string]any{
					"enabled": true,
				},
				"respectSchemaVersion": true,
			},
			"csharp": map[string]any{
				"rootNamespace":        "Pulumi",
				"respectSchemaVersion": true,
			},
		}).
		WithResources(
			infer.Resource(Pod{}),
			infer.Resource(Template{}),
			infer.Resource(Endpoint{}),
			infer.Resource(NetworkVolume{}),
			infer.Resource(Secret{}),
			infer.Resource(ContainerRegistryAuth{}),
		).
		WithFunctions(
			infer.Function(GetGpuTypes{}),
			infer.Function(GetCPUFlavors{}),
			infer.Function(GetDataCenters{}),
		).
		WithConfig(infer.Config(&Config{})).
		WithModuleMap(map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		}).Build()
	if err != nil {
		panic(fmt.Errorf("unable to build provider: %w", err))
	}
	return prov
}

// Config defines provider-level configuration.
type Config struct {
	APIKey string `pulumi:"apiKey,optional" provider:"secret"`
	APIURL string `pulumi:"apiUrl,optional"`
}

// Annotate provides descriptions for Config fields.
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.APIKey,
		"The RunPod API key for authentication. "+
			"Can also be set via the RUNPOD_API_KEY environment variable.")
	a.SetDefault(&c.APIKey, "", "RUNPOD_API_KEY")
	a.Describe(&c.APIURL,
		"The RunPod API URL. Defaults to https://api.runpod.io/graphql. "+
			"Can also be set via the RUNPOD_API_URL environment variable.")
	a.SetDefault(&c.APIURL, nil, "RUNPOD_API_URL")
}

// Configure validates the provider configuration.
func (c *Config) Configure(_ context.Context) error {
	if c.APIKey == "" {
		c.APIKey = os.Getenv("RUNPOD_API_KEY")
	}
	if c.APIURL == "" {
		c.APIURL = os.Getenv("RUNPOD_API_URL")
	}
	return nil
}

// getClient creates a genqlient GraphQL client from the provider config in context.
func getClient(ctx context.Context) graphql.Client {
	config := infer.GetConfig[Config](ctx)
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("RUNPOD_API_KEY")
	}
	apiURL := config.APIURL
	if apiURL == "" {
		apiURL = os.Getenv("RUNPOD_API_URL")
	}
	return runpod.NewClient(apiKey, apiURL)
}
