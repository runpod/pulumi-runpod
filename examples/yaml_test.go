//go:build yaml || all
// +build yaml all

package examples

import (
	"os"
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

func TestYAMLExampleLifecycle(t *testing.T) {
	if os.Getenv("RUNPOD_API_KEY") == "" {
		t.Skip("RUNPOD_API_KEY not set, skipping integration test")
	}

	pt := pulumitest.NewPulumiTest(t, "yaml",
		opttest.AttachProviderServer("runpod", providerFactory),
		opttest.SkipInstall(),
	)

	pt.Preview(t)
	pt.Up(t)
	pt.Destroy(t)
}
