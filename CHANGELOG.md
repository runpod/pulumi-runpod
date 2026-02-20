# Changelog

## v0.1.0-alpha.1

Initial pre-release of the RunPod Pulumi provider.

### Resources

- `runpod:index:Pod` — Create and manage GPU pods
- `runpod:index:Template` — Create and manage pod templates
- `runpod:index:Endpoint` — Create and manage serverless endpoints
- `runpod:index:NetworkVolume` — Create and manage network storage volumes
- `runpod:index:Secret` — Create and manage secrets
- `runpod:index:ContainerRegistryAuth` — Manage container registry authentication

### Functions

- `runpod:index:getGpuTypes` — Query available GPU types and pricing

### SDK Packages

- npm: `@runpod/pulumi`
- PyPI: `pulumi_runpod`
- NuGet: `Pulumi.Runpod`
- Go: `github.com/runpod/pulumi-runpod/sdk/go/runpod`
