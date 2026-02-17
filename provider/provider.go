package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/runpod/pulumi-runpod/pkg/runpod"
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
		WithPluginDownloadURL("github://api.github.com/runpod/pulumi-runpod").
		WithRepository("https://github.com/runpod/pulumi-runpod").
		WithLanguageMap(map[string]any{
			"go": map[string]any{
				"importBasePath":                 "github.com/runpod/pulumi-runpod/sdk/go/runpod",
				"generateResourceContainerTypes": true,
				"respectSchemaVersion":           true,
			},
			"nodejs": map[string]any{
				"packageName":         "@runpod/pulumi",
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
	a.Describe(&c.APIKey, "The RunPod API key for authentication. Can also be set via the RUNPOD_API_KEY environment variable.")
	a.SetDefault(&c.APIKey, nil, "RUNPOD_API_KEY")
	a.Describe(&c.APIURL, "The RunPod API URL. Defaults to https://api.runpod.io/graphql. Can also be set via the RUNPOD_API_URL environment variable.")
	a.SetDefault(&c.APIURL, nil, "RUNPOD_API_URL")
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
