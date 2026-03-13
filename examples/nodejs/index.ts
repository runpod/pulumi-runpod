import * as pulumi from "@pulumi/pulumi";
import * as runpod from "@runpod/pulumi";

const myTemplate = new runpod.Template("myTemplate", {
    name: `pulumi-nodejs-test-${pulumi.getProject()}-${pulumi.getStack()}`,
    imageName: "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04",
    containerDiskInGb: 20,
    volumeInGb: 20,
    startSsh: true,
});

export const templateId = myTemplate.templateId;
