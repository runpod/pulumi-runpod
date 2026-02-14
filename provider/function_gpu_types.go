package provider

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/runpod/pulumi-runpod/pkg/runpod"
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

	resp, err := runpod.GetGpuTypes(ctx, client, nil)
	if err != nil {
		return infer.FunctionResponse[GetGpuTypesResult]{}, err
	}

	result := make([]GpuTypeOutput, len(resp.GpuTypes))
	for i, g := range resp.GpuTypes {
		if g == nil {
			continue
		}
		result[i] = GpuTypeOutput{
			ID:             runpod.PtrString(g.Id),
			DisplayName:    runpod.PtrString(g.DisplayName),
			MemoryInGb:     runpod.PtrInt(g.MemoryInGb),
			SecureCloud:    runpod.PtrBool(g.SecureCloud),
			CommunityCloud: runpod.PtrBool(g.CommunityCloud),
			SecurePrice:    runpod.PtrFloat64(g.SecurePrice),
			CommunityPrice: runpod.PtrFloat64(g.CommunityPrice),
			MaxGpuCount:    runpod.PtrInt(g.MaxGpuCount),
		}
	}

	return infer.FunctionResponse[GetGpuTypesResult]{
		Output: GetGpuTypesResult{GpuTypes: result},
	}, nil
}
