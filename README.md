# Pulumi RunPod Provider

The Pulumi RunPod provider lets you manage [RunPod](https://www.runpod.io/) GPU cloud infrastructure using infrastructure as code.

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

config:
  runpod:apiKey:
    secret: true

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
