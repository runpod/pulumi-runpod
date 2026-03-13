//go:build go || all
// +build go all

package examples

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/stretchr/testify/require"
)

func TestGoExampleLifecycle(t *testing.T) {
	if os.Getenv("RUNPOD_API_KEY") == "" {
		t.Skip("RUNPOD_API_KEY not set, skipping integration test")
	}

	cwd, err := os.Getwd()
	require.NoError(t, err)

	module := filepath.Join(cwd, "../sdk/go/runpod")
	pt := pulumitest.NewPulumiTest(t, "go",
		opttest.GoModReplacement("github.com/runpod/pulumi-runpod/sdk/go/runpod", module),
		opttest.AttachProviderServer("runpod", providerFactory),
		opttest.SkipInstall(),
	)

	// The GoModReplacement option changes the replace directive, which can
	// invalidate go.sum entries. Run go mod tidy so the copied project
	// compiles cleanly under any Go toolchain version.
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = pt.WorkingDir()
	out, tidyErr := cmd.CombinedOutput()
	require.NoError(t, tidyErr, "go mod tidy failed: %s", string(out))

	pt.Preview(t)
	pt.Up(t)
	pt.Destroy(t)
}
