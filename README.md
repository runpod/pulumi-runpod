# Pulumi Runpod Provider

> [!WARNING]
> **This provider is deprecated and unmaintained.** It calls the Runpod GraphQL API and will get no further releases. Use the [Runpod Terraform provider](https://registry.terraform.io/providers/runpod/runpod/latest) 2.0.0, which runs on the Runpod REST v2 API. You can keep writing Pulumi programs: see [Migrating to the Terraform provider](#migrating-to-the-terraform-provider).

## Migrating to the Terraform provider

Pulumi can consume a Terraform provider directly. From your Pulumi project:

```bash
pulumi package add terraform-provider registry.terraform.io/runpod/runpod
```

That generates a local SDK for your project's language and pins Runpod provider 2.0.0 in `Pulumi.yaml`. Set your API key:

```bash
pulumi config set --secret runpod:apiKey YOUR_API_KEY
```

The provider also reads `RUNPOD_API_KEY` from the environment.

### Resource mapping

| This provider | Terraform provider 2.0.0, through Pulumi |
|---|---|
| `runpod:Pod` | `runpod.Pod` |
| `runpod:Template` | `runpod.Template` |
| `runpod:Endpoint` | `runpod.ServerlessEndpoint` |
| `runpod:NetworkVolume` | `runpod.NetworkVolume` |
| `runpod:ContainerRegistryAuth` | `runpod.Registry` |
| `runpod:Secret` | no counterpart yet, see below |
| `runpod:getGpuTypes` | `runpod.getGpuTypes` |
| `runpod:getCPUFlavors` | `runpod.getCpuTypes` |
| `runpod:getDataCenters` | `runpod.getDataCenters` |

The Terraform provider also covers global volumes, Instant Clusters, account SSH keys, and billing and catalog data sources that this provider never had.

### Secrets

The Terraform provider has no secret resource yet. Until it does, manage secrets in the [Runpod console](https://www.runpod.io/console/user/settings) or through the REST v2 `/v2/account/secrets` endpoints, and reference them from pod and endpoint environment variables.

---

**Everything below describes the deprecated provider and is kept for existing users.**

The Pulumi Runpod provider lets you manage [Runpod](https://www.runpod.io/) GPU cloud infrastructure using infrastructure as code.

<!-- runpod-agent-setup:start -->
**For AI agents:** setup instructions for the Runpod skills and MCP server are at https://agents.runpod.io/
<!-- runpod-agent-setup:end -->

## Resources

| Resource | Description |
|----------|-------------|
| `runpod:Template` | Container templates for pods and serverless endpoints |
| `runpod:Pod` | GPU and CPU pod instances |
| `runpod:Endpoint` | Serverless GPU endpoints |
| `runpod:NetworkVolume` | Persistent network-attached storage volumes |
| `runpod:Secret` | Encrypted secrets for use in pods and endpoints |
| `runpod:ContainerRegistryAuth` | Container registry authentication credentials |

## Functions

| Function | Description |
|----------|-------------|
| `runpod:getGpuTypes` | List available GPU types with pricing and availability |
| `runpod:getCPUFlavors` | List available CPU configurations |
| `runpod:getDataCenters` | List data centers with GPU availability |

## Installation

The provider plugin is installed automatically when you use it in a Pulumi program.

### Node.js (TypeScript/JavaScript)

```bash
npm install @runpod/pulumi
```

### Python

```bash
pip install pulumi_runpod
```

### Go

```bash
go get github.com/runpod/pulumi-runpod/sdk/go/runpod
```

### .NET

```bash
dotnet add package Pulumi.Runpod
```

## Configuration

Set your RunPod API key:

```bash
pulumi config set runpod:apiKey --secret YOUR_API_KEY
```

Or use the `RUNPOD_API_KEY` environment variable.

Optionally set a custom API URL (defaults to `https://api.runpod.io/graphql`):

```bash
pulumi config set runpod:apiUrl https://api.runpod.io/graphql
```

## Examples

### TypeScript

```typescript
import * as runpod from "@runpod/pulumi";

const template = new runpod.Template("myTemplate", {
    name: "my-template",
    imageName: "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04",
    containerDiskInGb: 20,
    volumeInGb: 20,
    startSsh: true,
});

const endpoint = new runpod.Endpoint("myEndpoint", {
    name: "my-endpoint",
    templateId: template.templateId,
    gpuIds: "AMPERE_16",
    workersMin: 0,
    workersMax: 3,
    idleTimeout: 5,
});

export const endpointId = endpoint.endpointId;
```

### Python

```python
import pulumi
import pulumi_runpod as runpod

template = runpod.Template("myTemplate",
    name="my-template",
    image_name="runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04",
    container_disk_in_gb=20,
    volume_in_gb=20,
    start_ssh=True,
)

endpoint = runpod.Endpoint("myEndpoint",
    name="my-endpoint",
    template_id=template.template_id,
    gpu_ids="AMPERE_16",
    workers_min=0,
    workers_max=3,
    idle_timeout=5,
)

pulumi.export("endpointId", endpoint.endpoint_id)
```

### Go

```go
package main

import (
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    runpod "github.com/runpod/pulumi-runpod/sdk/go/runpod"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        template, err := runpod.NewTemplate(ctx, "myTemplate", &runpod.TemplateArgs{
            Name:              pulumi.String("my-template"),
            ImageName:         pulumi.String("runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04"),
            ContainerDiskInGb: pulumi.Int(20),
            VolumeInGb:        pulumi.Int(20),
            StartSsh:          pulumi.Bool(true),
        })
        if err != nil {
            return err
        }
        ctx.Export("templateId", template.TemplateId)
        return nil
    })
}
```

### YAML

```yaml
name: runpod-example
runtime: yaml

# Set runpod:apiKey via: pulumi config set runpod:apiKey --secret YOUR_KEY
# Or use the RUNPOD_API_KEY environment variable.

resources:
  myTemplate:
    type: runpod:Template
    properties:
      name: my-template
      imageName: runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04
      containerDiskInGb: 20
      volumeInGb: 20
      startSsh: true

  myEndpoint:
    type: runpod:Endpoint
    properties:
      name: my-endpoint
      templateId: ${myTemplate.templateId}
      gpuIds: AMPERE_16
      workersMin: 0
      workersMax: 3
      idleTimeout: 5

outputs:
  endpointId: ${myEndpoint.endpointId}
```

## Development

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Pulumi CLI](https://www.pulumi.com/docs/install/)
- [pulumictl](https://github.com/pulumi/pulumictl#installation)
- [Node.js 14+](https://nodejs.org/)
- [Python 3](https://www.python.org/downloads/)
- [.NET 8+](https://dotnet.microsoft.com/download)

### Build

```bash
make build install
```

### Test

```bash
make test_provider
```

## License

Apache 2.0. See [LICENSE](LICENSE).
