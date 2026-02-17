//go:build yaml || all
// +build yaml all

package examples

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

func TestYAMLExampleLifecycle(t *testing.T) {
	pt := pulumitest.NewPulumiTest(t, "yaml",
		opttest.AttachProviderServer("runpod", providerFactory),
		opttest.SkipInstall(),
	)

	pt.Preview(t)
	pt.Up(t)
	pt.Destroy(t)
}
