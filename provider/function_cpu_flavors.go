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

// GetCpuFlavors is the controller for the runpod:index:getCpuFlavors function (invoke).
type GetCpuFlavors struct{}

// GetCpuFlavorsArgs are the inputs for the CPU flavors query.
type GetCpuFlavorsArgs struct {
	// SlsOnly filters to serverless-only flavors.
	SlsOnly *bool `pulumi:"slsOnly,optional"`
	// IsSls filters by serverless eligibility.
	IsSls *bool `pulumi:"isSls,optional"`
}

// Annotate provides descriptions for GetCpuFlavorsArgs fields.
func (a *GetCpuFlavorsArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.SlsOnly, "When true, return only serverless-eligible CPU flavors.")
	ann.Describe(&a.IsSls, "Filter by serverless eligibility.")
}

// GetCpuFlavorsResult is the output of the CPU flavors query.
type GetCpuFlavorsResult struct {
	CpuFlavors []CpuFlavorOutput `pulumi:"cpuFlavors"`
}

// Annotate provides descriptions for GetCpuFlavorsResult fields.
func (r *GetCpuFlavorsResult) Annotate(a infer.Annotator) {
	a.Describe(&r.CpuFlavors, "The list of available CPU instance flavors.")
}

// CpuFlavorOutput represents a single CPU flavor in the output.
type CpuFlavorOutput struct {
	ID               string  `pulumi:"id"`
	GroupID          string  `pulumi:"groupId"`
	GroupName        string  `pulumi:"groupName"`
	DisplayName      string  `pulumi:"displayName"`
	MinVcpu          float64 `pulumi:"minVcpu"`
	MaxVcpu          int     `pulumi:"maxVcpu"`
	VcpuBurstable    bool    `pulumi:"vcpuBurstable"`
	RamMultiplier    float64 `pulumi:"ramMultiplier"`
	DiskLimitPerVcpu int     `pulumi:"diskLimitPerVcpu"`
}

// Annotate provides descriptions for CpuFlavorOutput fields.
func (c *CpuFlavorOutput) Annotate(a infer.Annotator) {
	a.Describe(&c.ID, "The unique identifier of the CPU flavor (used as flavorId in instanceId).")
	a.Describe(&c.GroupID, "The group this flavor belongs to.")
	a.Describe(&c.GroupName, "The display name of the flavor group.")
	a.Describe(&c.DisplayName, "The human-readable name of the CPU flavor.")
	a.Describe(&c.MinVcpu, "The minimum number of vCPUs for this flavor.")
	a.Describe(&c.MaxVcpu, "The maximum number of vCPUs for this flavor.")
	a.Describe(&c.VcpuBurstable, "Whether vCPUs are burstable.")
	a.Describe(&c.RamMultiplier, "RAM allocated per vCPU (in GB).")
	a.Describe(&c.DiskLimitPerVcpu, "Disk limit per vCPU (in GB).")
}

// Invoke executes the CPU flavors query.
func (GetCpuFlavors) Invoke(
	ctx context.Context,
	req infer.FunctionRequest[GetCpuFlavorsArgs],
) (infer.FunctionResponse[GetCpuFlavorsResult], error) {
	client := getClient(ctx)

	var input *runpod.CpuFlavorInput
	if req.Input.SlsOnly != nil || req.Input.IsSls != nil {
		input = &runpod.CpuFlavorInput{
			SlsOnly: req.Input.SlsOnly,
			IsSls:   req.Input.IsSls,
		}
	}

	resp, err := runpod.GetCpuFlavors(ctx, client, input)
	if err != nil {
		return infer.FunctionResponse[GetCpuFlavorsResult]{}, err
	}

	result := make([]CpuFlavorOutput, 0, len(resp.CpuFlavors))
	for _, f := range resp.CpuFlavors {
		if f == nil {
			continue
		}
		result = append(result, CpuFlavorOutput{
			ID:               runpod.PtrString(f.Id),
			GroupID:          runpod.PtrString(f.GroupId),
			GroupName:        runpod.PtrString(f.GroupName),
			DisplayName:      runpod.PtrString(f.DisplayName),
			MinVcpu:          runpod.PtrFloat64(f.MinVcpu),
			MaxVcpu:          runpod.PtrInt(f.MaxVcpu),
			VcpuBurstable:    runpod.PtrBool(f.VcpuBurstable),
			RamMultiplier:    runpod.PtrFloat64(f.RamMultiplier),
			DiskLimitPerVcpu: runpod.PtrInt(f.DiskLimitPerVcpu),
		})
	}

	return infer.FunctionResponse[GetCpuFlavorsResult]{
		Output: GetCpuFlavorsResult{CpuFlavors: result},
	}, nil
}
