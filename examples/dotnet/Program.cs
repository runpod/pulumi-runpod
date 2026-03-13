using System.Collections.Generic;
using Pulumi;
using Runpod = Pulumi.Runpod;

return await Deployment.RunAsync(() =>
{
    var myTemplate = new Runpod.Template("myTemplate", new()
    {
        Name = $"my-pulumi-template-{Deployment.Instance.StackName}",
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
