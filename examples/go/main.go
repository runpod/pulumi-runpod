package main

import (
	boilerplate "github.com/runpod/pulumi-runpod/sdk/go/pulumi-runpod"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myRandomResource, err := boilerplate.NewRandom(ctx, "myRandomResource", &boilerplate.RandomArgs{
			Length: pulumi.Int(24),
		})
		if err != nil {
			return err
		}
		_, err = boilerplate.NewRandomComponent(ctx, "myRandomComponent", &boilerplate.RandomComponentArgs{
			Length: pulumi.Int(24),
		})
		if err != nil {
			return err
		}
		ctx.Export("output", pulumi.StringMap{
			"value": myRandomResource.Result,
		})
		return nil
	})
}
