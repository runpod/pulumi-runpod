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
