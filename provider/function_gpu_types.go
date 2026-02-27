// Copyright 2025, RunPod, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/runpod/pulumi-runpod/provider/pkg/runpod"
)

// GetGpuTypes is the controller for the runpod:index:getGpuTypes function (invoke).
type GetGpuTypes struct{}

// GetGpuTypesArgs are the (empty) inputs for the GPU types query.
type GetGpuTypesArgs struct{}

// GetGpuTypesResult is the output of the GPU types query.
type GetGpuTypesResult struct {
	GpuTypes []GpuTypeOutput `pulumi:"gpuTypes"`
}

// LowestPriceOutput represents the lowest available pricing for a GPU type.
type LowestPriceOutput struct {
	MinimumBidPrice    float64 `pulumi:"minimumBidPrice"`
	UninterruptablePrice float64 `pulumi:"uninterruptablePrice"`
	RentedCount        int     `pulumi:"rentedCount"`
	TotalCount         int     `pulumi:"totalCount"`
	StockStatus        string  `pulumi:"stockStatus"`
}

// GpuTypeOutput represents a single GPU type in the output.
type GpuTypeOutput struct {
	ID             string             `pulumi:"id"`
	DisplayName    string             `pulumi:"displayName"`
	MemoryInGb     int                `pulumi:"memoryInGb"`
	SecureCloud    bool               `pulumi:"secureCloud"`
	CommunityCloud bool               `pulumi:"communityCloud"`
	SecurePrice    float64            `pulumi:"securePrice"`
	CommunityPrice float64            `pulumi:"communityPrice"`
	MaxGpuCount    int                `pulumi:"maxGpuCount"`
	LowestPrice    *LowestPriceOutput `pulumi:"lowestPrice,optional"`
}

// Annotate provides descriptions for GetGpuTypesResult fields.
func (r *GetGpuTypesResult) Annotate(a infer.Annotator) {
	a.Describe(&r.GpuTypes, "The list of available GPU types.")
}

// Annotate provides descriptions for GpuTypeOutput fields.
func (g *GpuTypeOutput) Annotate(a infer.Annotator) {
	a.Describe(&g.ID, "The unique identifier of the GPU type.")
	a.Describe(&g.DisplayName, "The display name of the GPU type.")
	a.Describe(&g.MemoryInGb, "The amount of VRAM in GB.")
	a.Describe(&g.SecureCloud,
		"Whether the GPU is available in secure cloud.")
	a.Describe(&g.CommunityCloud,
		"Whether the GPU is available in community cloud.")
	a.Describe(&g.SecurePrice,
		"The price per hour in secure cloud (USD).")
	a.Describe(&g.CommunityPrice,
		"The price per hour in community cloud (USD).")
	a.Describe(&g.MaxGpuCount,
		"The maximum number of this GPU type that can be allocated.")
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

	result := make([]GpuTypeOutput, 0, len(resp.GpuTypes))
	for _, g := range resp.GpuTypes {
		if g == nil {
			continue
		}
		out := GpuTypeOutput{
			ID:             runpod.PtrString(g.Id),
			DisplayName:    runpod.PtrString(g.DisplayName),
			MemoryInGb:     runpod.PtrInt(g.MemoryInGb),
			SecureCloud:    runpod.PtrBool(g.SecureCloud),
			CommunityCloud: runpod.PtrBool(g.CommunityCloud),
			SecurePrice:    runpod.PtrFloat64(g.SecurePrice),
			CommunityPrice: runpod.PtrFloat64(g.CommunityPrice),
			MaxGpuCount:    runpod.PtrInt(g.MaxGpuCount),
		}
		if g.LowestPrice != nil {
			out.LowestPrice = &LowestPriceOutput{
				MinimumBidPrice:      runpod.PtrFloat64(g.LowestPrice.MinimumBidPrice),
				UninterruptablePrice: runpod.PtrFloat64(g.LowestPrice.UninterruptablePrice),
				RentedCount:          runpod.PtrInt(g.LowestPrice.RentedCount),
				TotalCount:           runpod.PtrInt(g.LowestPrice.TotalCount),
				StockStatus:          runpod.PtrString(g.LowestPrice.StockStatus),
			}
		}
		result = append(result, out)
	}

	return infer.FunctionResponse[GetGpuTypesResult]{
		Output: GetGpuTypesResult{GpuTypes: result},
	}, nil
}
