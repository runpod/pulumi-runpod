package examples

import (
	"github.com/pulumi/providertest/providers"
	goprovider "github.com/pulumi/pulumi-go-provider"
	"github.com/runpod/pulumi-runpod/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var providerFactory = func(_ providers.PulumiTest) (pulumirpc.ResourceProviderServer, error) {
	return goprovider.RawServer("runpod", "1.0.0", provider.Provider())(nil)
}
