---
title: RunPod
meta_desc: Provides an overview of the RunPod Provider for Pulumi.
layout: overview
---

The RunPod provider for Pulumi can be used to provision and manage [RunPod](https://www.runpod.io) GPU cloud resources including pods, templates, serverless endpoints, network volumes, secrets, and container registry credentials.

To manage RunPod resources with Pulumi, you need a RunPod API key. You can create one from the [RunPod Console](https://www.runpod.io/console/user/settings).

## Example

{{< chooser language "typescript,go,python,csharp" >}}
{{% choosable language typescript %}}

```typescript
import * as runpod from "@runpod/pulumi";

const myTemplate = new runpod.Template("myTemplate", {
    name: "my-pulumi-template",
    imageName: "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04",
    containerDiskInGb: 20,
    volumeInGb: 20,
    startSsh: true,
});

export const templateId = myTemplate.templateId;
```

{{% /choosable %}}
{{% choosable language go %}}

```go
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	runpod "github.com/runpod/pulumi-runpod/sdk/go/runpod"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myTemplate, err := runpod.NewTemplate(ctx, "myTemplate", &runpod.TemplateArgs{
			Name:              pulumi.String("my-pulumi-template"),
			ImageName:         pulumi.String("runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04"),
			ContainerDiskInGb: pulumi.Int(20),
			VolumeInGb:        pulumi.Int(20),
			StartSsh:          pulumi.Bool(true),
		})
		if err != nil {
			return err
		}
		ctx.Export("templateId", myTemplate.TemplateId)
		return nil
	})
}
```

{{% /choosable %}}
{{% choosable language python %}}

```python
import pulumi
import pulumi_runpod as runpod

my_template = runpod.Template("myTemplate",
    name="my-pulumi-template",
    image_name="runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04",
    container_disk_in_gb=20,
    volume_in_gb=20,
    start_ssh=True,
)

pulumi.export("templateId", my_template.template_id)
```

{{% /choosable %}}
{{% choosable language csharp %}}

```csharp
using System.Collections.Generic;
using Pulumi;
using Runpod = Pulumi.Runpod;

return await Deployment.RunAsync(() =>
{
    var myTemplate = new Runpod.Template("myTemplate", new()
    {
        Name = "my-pulumi-template",
        ImageName = "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04",
        ContainerDiskInGb = 20,
        VolumeInGb = 20,
        StartSsh = true,
    });

    return new Dictionary<string, object?>
    {
        ["templateId"] = myTemplate.TemplateId,
    };
});
```

{{% /choosable %}}
{{< /chooser >}}
