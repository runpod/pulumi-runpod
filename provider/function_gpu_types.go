package provider

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// GetGpuTypes is the controller for the runpod:index:getGpuTypes function (invoke).
type GetGpuTypes struct{}

// GetGpuTypesArgs are the (empty) inputs for the GPU types query.
type GetGpuTypesArgs struct{}

// GetGpuTypesResult is the output of the GPU types query.
type GetGpuTypesResult struct {
	GpuTypes []GpuTypeOutput `pulumi:"gpuTypes"`
}

// GpuTypeOutput represents a single GPU type in the output.
type GpuTypeOutput struct {
	ID             string  `pulumi:"id"`
	DisplayName    string  `pulumi:"displayName"`
	MemoryInGb     int     `pulumi:"memoryInGb"`
	SecureCloud    bool    `pulumi:"secureCloud"`
	CommunityCloud bool    `pulumi:"communityCloud"`
	SecurePrice    float64 `pulumi:"securePrice"`
	CommunityPrice float64 `pulumi:"communityPrice"`
	MaxGpuCount    int     `pulumi:"maxGpuCount"`
}

// Invoke executes the GPU types query.
func (GetGpuTypes) Invoke(
	ctx context.Context,
	req infer.FunctionRequest[GetGpuTypesArgs],
) (infer.FunctionResponse[GetGpuTypesResult], error) {
	client := getClient(ctx)

	gpuTypes, err := client.GetGpuTypes(ctx)
	if err != nil {
		return infer.FunctionResponse[GetGpuTypesResult]{}, err
	}

	result := make([]GpuTypeOutput, len(gpuTypes))
	for i, g := range gpuTypes {
		result[i] = GpuTypeOutput{
			ID:             g.ID,
			DisplayName:    g.DisplayName,
			MemoryInGb:     g.MemoryInGb,
			SecureCloud:    g.SecureCloud,
			CommunityCloud: g.CommunityCloud,
			SecurePrice:    g.SecurePrice,
			CommunityPrice: g.CommunityPrice,
			MaxGpuCount:    g.MaxGpuCount,
		}
	}

	return infer.FunctionResponse[GetGpuTypesResult]{
		Output: GetGpuTypesResult{GpuTypes: result},
	}, nil
}
