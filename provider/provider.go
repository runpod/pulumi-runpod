package provider

import (
	"context"
	"fmt"
	"os"

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
		WithResources(
			infer.Resource(Pod{}),
			infer.Resource(Template{}),
			infer.Resource(Endpoint{}),
			infer.Resource(NetworkVolume{}),
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

// getClient creates a RunPod API client from the provider config in context.
func getClient(ctx context.Context) *runpod.Client {
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
