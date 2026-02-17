import * as runpod from "@pulumi/runpod";

const myTemplate = new runpod.Template("myTemplate", {
    name: "my-pulumi-template",
    imageName: "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04",
    containerDiskInGb: 20,
    volumeInGb: 20,
    startSsh: true,
});

export const templateId = myTemplate.templateId;
