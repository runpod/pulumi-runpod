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

// GetDataCenters is the controller for the runpod:index:getDataCenters function (invoke).
type GetDataCenters struct{}

// GetDataCentersArgs are the (empty) inputs for the data centers query.
type GetDataCentersArgs struct{}

// GetDataCentersResult is the output of the data centers query.
type GetDataCentersResult struct {
	DataCenters []DataCenterOutput `pulumi:"dataCenters"`
}

// Annotate provides descriptions for GetDataCentersResult fields.
func (r *GetDataCentersResult) Annotate(a infer.Annotator) {
	a.Describe(&r.DataCenters, "The list of available RunPod data centers.")
}

// DataCenterOutput represents a single data center in the output.
type DataCenterOutput struct {
	ID              string                `pulumi:"id"`
	Name            string                `pulumi:"name"`
	Location        string                `pulumi:"location"`
	Region          string                `pulumi:"region"`
	Listed          bool                  `pulumi:"listed"`
	StorageSupport  bool                  `pulumi:"storageSupport"`
	GlobalNetwork   bool                  `pulumi:"globalNetwork"`
	Compliance      []string              `pulumi:"compliance"`
	GpuAvailability []GpuAvailabilityItem `pulumi:"gpuAvailability"`
}

// Annotate provides descriptions for DataCenterOutput fields.
func (d *DataCenterOutput) Annotate(a infer.Annotator) {
	a.Describe(&d.ID, "The unique identifier of the data center (used as dataCenterId).")
	a.Describe(&d.Name, "The display name of the data center.")
	a.Describe(&d.Location, "The geographic location of the data center.")
	a.Describe(&d.Region, "The broad region (e.g. NORTH_AMERICA, EUROPE).")
	a.Describe(&d.Listed, "Whether this data center is publicly listed.")
	a.Describe(&d.StorageSupport, "Whether this data center supports network volumes.")
	a.Describe(&d.GlobalNetwork, "Whether this data center is part of the global network.")
	a.Describe(&d.Compliance, "Compliance certifications held by this data center.")
	a.Describe(&d.GpuAvailability, "GPU availability within this data center.")
}

// GpuAvailabilityItem represents GPU availability for a specific GPU type at a data center.
type GpuAvailabilityItem struct {
	GpuTypeID          string `pulumi:"gpuTypeId"`
	GpuTypeDisplayName string `pulumi:"gpuTypeDisplayName"`
	Available          bool   `pulumi:"available"`
	StockStatus        string `pulumi:"stockStatus"`
}

// Annotate provides descriptions for GpuAvailabilityItem fields.
func (g *GpuAvailabilityItem) Annotate(a infer.Annotator) {
	a.Describe(&g.GpuTypeID, "The GPU type identifier.")
	a.Describe(&g.GpuTypeDisplayName, "The human-readable GPU type name.")
	a.Describe(&g.Available, "Whether this GPU type is currently available at this data center.")
	a.Describe(&g.StockStatus, "Current stock status (e.g. High, Medium, Low).")
}

// ptrRegionString safely dereferences a *runpod.DataCenterRegion to string.
func ptrRegionString(r *runpod.DataCenterRegion) string {
	if r == nil {
		return ""
	}
	return string(*r)
}

// Invoke executes the data centers query.
func (GetDataCenters) Invoke(
	ctx context.Context,
	_ infer.FunctionRequest[GetDataCentersArgs],
) (infer.FunctionResponse[GetDataCentersResult], error) {
	client := getClient(ctx)

	resp, err := runpod.GetDataCenters(ctx, client)
	if err != nil {
		return infer.FunctionResponse[GetDataCentersResult]{}, err
	}

	result := make([]DataCenterOutput, 0, len(resp.DataCenters))
	for _, dc := range resp.DataCenters {
		if dc == nil {
			continue
		}

		compliance := make([]string, 0, len(dc.Compliance))
		for _, c := range dc.Compliance {
			if c != nil {
				compliance = append(compliance, string(*c))
			}
		}

		gpuAvail := make([]GpuAvailabilityItem, 0, len(dc.GpuAvailability))
		for _, g := range dc.GpuAvailability {
			if g == nil {
				continue
			}
			gpuAvail = append(gpuAvail, GpuAvailabilityItem{
				GpuTypeID:          runpod.PtrString(g.GpuTypeId),
				GpuTypeDisplayName: runpod.PtrString(g.GpuTypeDisplayName),
				Available:          runpod.PtrBool(g.Available),
				StockStatus:        runpod.PtrString(g.StockStatus),
			})
		}

		result = append(result, DataCenterOutput{
			ID:              runpod.PtrString(dc.Id),
			Name:            runpod.PtrString(dc.Name),
			Location:        runpod.PtrString(dc.Location),
			Region:          ptrRegionString(dc.Region),
			Listed:          runpod.PtrBool(dc.Listed),
			StorageSupport:  runpod.PtrBool(dc.StorageSupport),
			GlobalNetwork:   runpod.PtrBool(dc.GlobalNetwork),
			Compliance:      compliance,
			GpuAvailability: gpuAvail,
		})
	}

	return infer.FunctionResponse[GetDataCentersResult]{
		Output: GetDataCentersResult{DataCenters: result},
	}, nil
}
