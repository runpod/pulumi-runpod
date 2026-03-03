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

// Endpoint is the controller for the runpod:index:Endpoint resource.
type Endpoint struct{}

// EndpointArgs are the inputs for creating a serverless endpoint.
type EndpointArgs struct {
	Name            string            `pulumi:"name"`
	TemplateID      *string           `pulumi:"templateId,optional"`
	GpuIDs          *string           `pulumi:"gpuIds,optional"`
	WorkersMin      *int              `pulumi:"workersMin,optional"`
	WorkersMax      *int              `pulumi:"workersMax,optional"`
	IdleTimeout     *int              `pulumi:"idleTimeout,optional"`
	Locations       *string           `pulumi:"locations,optional"`
	ScalerType      *string           `pulumi:"scalerType,optional"`
	ScalerValue     *int              `pulumi:"scalerValue,optional"`
	NetworkVolumeID *string           `pulumi:"networkVolumeId,optional"`
	GpuCount        *int              `pulumi:"gpuCount,optional"`
	InstanceIDs     []string          `pulumi:"instanceIds,optional"`
	Env             map[string]string `pulumi:"env,optional"`
	// New fields
	FlashBootType       *string  `pulumi:"flashBootType,optional"`
	ExecutionTimeoutMs  *int     `pulumi:"executionTimeoutMs,optional"`
	AllowedCudaVersions *string  `pulumi:"allowedCudaVersions,optional"`
	MinCudaVersion      *string  `pulumi:"minCudaVersion,optional"`
	FlashEnvironmentID  *string  `pulumi:"flashEnvironmentId,optional"`
	BindEndpoint        *bool    `pulumi:"bindEndpoint,optional"`
	HubReleaseID        *string  `pulumi:"hubReleaseId,optional"`
	Type                *string  `pulumi:"type,optional"`
	ModelName           *string  `pulumi:"modelName,optional"`
	ModelReferences     []string `pulumi:"modelReferences,optional"`
}

// Annotate provides descriptions for EndpointArgs fields.
func (a *EndpointArgs) Annotate(an infer.Annotator) {
	an.Describe(&a.Name, "A name for the endpoint.")
	an.Describe(&a.TemplateID,
		"The template ID to use for the endpoint workers.")
	an.Describe(&a.GpuIDs,
		"The GPU type IDs to use (e.g. \"AMPERE_16\").")
	an.Describe(&a.WorkersMin,
		"The minimum number of workers to keep running.")
	an.Describe(&a.WorkersMax,
		"The maximum number of workers to scale up to.")
	an.Describe(&a.IdleTimeout,
		"The number of seconds a worker can remain idle "+
			"before being scaled down.")
	an.Describe(&a.Locations,
		"Comma-separated data center locations for worker deployment.")
	an.Describe(&a.ScalerType,
		"The autoscaler type (e.g. \"QUEUE_DELAY\", \"REQUEST_COUNT\").")
	an.Describe(&a.ScalerValue,
		"The autoscaler target value.")
	an.Describe(&a.NetworkVolumeID,
		"The network volume ID to attach to endpoint workers.")
	an.Describe(&a.GpuCount,
		"The number of GPUs per worker.")
	an.Describe(&a.InstanceIDs,
		"Specific instance IDs to use for workers.")
	an.Describe(&a.Env,
		"Environment variables as key-value pairs.")
	an.Describe(&a.FlashBootType,
		"The flash boot type.")
	an.Describe(&a.ExecutionTimeoutMs,
		"Maximum execution time in milliseconds "+
			"before a request is terminated.")
	an.Describe(&a.AllowedCudaVersions,
		"Comma-separated list of allowed CUDA versions.")
	an.Describe(&a.MinCudaVersion,
		"The minimum CUDA version required.")
	an.Describe(&a.FlashEnvironmentID,
		"The flash environment ID.")
	an.Describe(&a.BindEndpoint,
		"Whether to bind the endpoint to specific workers.")
	an.Describe(&a.Type, "The endpoint type.")
	an.Describe(&a.ModelName,
		"The model name for the endpoint.")
	an.Describe(&a.HubReleaseID,
		"The hub release ID for the endpoint.")
	an.Describe(&a.ModelReferences,
		"Model references for the endpoint.")
}

// EndpointNetworkVolumeBinding represents a network volume attached to an endpoint in a specific data center.
type EndpointNetworkVolumeBinding struct {
	NetworkVolumeID string `pulumi:"networkVolumeId"`
	DataCenterID    string `pulumi:"dataCenterId"`
}

// EndpointState is the persisted state of an endpoint resource.
type EndpointState struct {
	EndpointArgs
	EndpointID       string                         `pulumi:"endpointId"`
	NetworkVolumeIDs []EndpointNetworkVolumeBinding `pulumi:"networkVolumeIds,optional"`
}

// Annotate provides descriptions for EndpointState fields.
func (s *EndpointState) Annotate(a infer.Annotator) {
	a.Describe(&s.EndpointID,
		"The unique identifier of the endpoint.")
	a.Describe(&s.NetworkVolumeIDs,
		"Network volumes attached to the endpoint, returned by the API.")
}

// Create creates a new serverless endpoint.
func (Endpoint) Create(
	ctx context.Context,
	req infer.CreateRequest[EndpointArgs],
) (infer.CreateResponse[EndpointState], error) {
	input := req.Inputs
	if req.DryRun {
		return infer.CreateResponse[EndpointState]{
			ID:     req.Name,
			Output: EndpointState{EndpointArgs: input},
		}, nil
	}

	client := getClient(ctx)
	saveInput := endpointArgsToInput(nil, input)

	resp, err := runpod.SaveEndpoint(ctx, client, saveInput)
	if err != nil {
		return infer.CreateResponse[EndpointState]{}, err
	}

	ep := &resp.SaveEndpoint
	state := endpointResponseToState(input, ep)
	return infer.CreateResponse[EndpointState]{
		ID:     runpod.PtrString(ep.Id),
		Output: state,
	}, nil
}

// Read refreshes the endpoint state from the API.
func (Endpoint) Read(
	ctx context.Context,
	req infer.ReadRequest[EndpointArgs, EndpointState],
) (infer.ReadResponse[EndpointArgs, EndpointState], error) {
	client := getClient(ctx)

	resp, err := runpod.GetMyEndpoints(ctx, client)
	if err != nil {
		return infer.ReadResponse[EndpointArgs, EndpointState]{}, err
	}

	if resp.Myself == nil {
		return infer.ReadResponse[EndpointArgs, EndpointState]{
			ID: "",
		}, nil
	}

	for _, e := range resp.Myself.Endpoints {
		if e != nil && runpod.PtrString(e.Id) == req.ID {
			state := endpointResponseToState(req.Inputs, e)
			return infer.ReadResponse[EndpointArgs, EndpointState]{
				ID:     req.ID,
				Inputs: req.Inputs,
				State:  state,
			}, nil
		}
	}

	return infer.ReadResponse[EndpointArgs, EndpointState]{
		ID: "",
	}, nil
}

// Update modifies an endpoint using the upsert pattern (saveEndpoint with id).
func (Endpoint) Update(
	ctx context.Context,
	req infer.UpdateRequest[EndpointArgs, EndpointState],
) (infer.UpdateResponse[EndpointState], error) {
	if req.DryRun {
		return infer.UpdateResponse[EndpointState]{
			Output: EndpointState{EndpointArgs: req.Inputs},
		}, nil
	}

	client := getClient(ctx)
	id := req.ID
	saveInput := endpointArgsToInput(&id, req.Inputs)

	resp, err := runpod.SaveEndpoint(ctx, client, saveInput)
	if err != nil {
		return infer.UpdateResponse[EndpointState]{}, err
	}

	ep := &resp.SaveEndpoint
	state := endpointResponseToState(req.Inputs, ep)
	return infer.UpdateResponse[EndpointState]{Output: state}, nil
}

// Delete removes an endpoint.
func (Endpoint) Delete(
	ctx context.Context,
	req infer.DeleteRequest[EndpointState],
) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	if _, err := runpod.DeleteEndpoint(ctx, client, req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func endpointArgsToInput(
	id *string, args EndpointArgs,
) runpod.EndpointInput {
	input := runpod.EndpointInput{
		Id:                  id,
		Name:                args.Name,
		TemplateId:          args.TemplateID,
		GpuIds:              args.GpuIDs,
		WorkersMin:          args.WorkersMin,
		WorkersMax:          args.WorkersMax,
		IdleTimeout:         args.IdleTimeout,
		Locations:           args.Locations,
		ScalerType:          args.ScalerType,
		ScalerValue:         args.ScalerValue,
		NetworkVolumeId:     args.NetworkVolumeID,
		GpuCount:            args.GpuCount,
		InstanceIds:         runpod.StringPtrSlice(args.InstanceIDs),
		Env:                 runpod.EnvMapToGQL(args.Env),
		ExecutionTimeoutMs:  args.ExecutionTimeoutMs,
		AllowedCudaVersions: args.AllowedCudaVersions,
		MinCudaVersion:      args.MinCudaVersion,
		FlashEnvironmentId:  args.FlashEnvironmentID,
		BindEndpoint:        args.BindEndpoint,
		HubReleaseId:        args.HubReleaseID,
		Type:                args.Type,
		ModelName:           args.ModelName,
		ModelReferences:     runpod.StringPtrSlice(args.ModelReferences),
	}

	if args.FlashBootType != nil {
		fbt := runpod.FlashBootType(*args.FlashBootType)
		input.FlashBootType = &fbt
	}

	return input
}

func endpointResponseToState(
	input EndpointArgs, ep *runpod.EndpointResponse,
) EndpointState {
	state := EndpointState{
		EndpointArgs: input,
		EndpointID:   runpod.PtrString(ep.Id),
	}

	if len(ep.NetworkVolumeIds) > 0 {
		bindings := make([]EndpointNetworkVolumeBinding, 0, len(ep.NetworkVolumeIds))
		for _, nv := range ep.NetworkVolumeIds {
			if nv != nil {
				bindings = append(bindings, EndpointNetworkVolumeBinding{
					NetworkVolumeID: runpod.PtrString(nv.NetworkVolumeId),
					DataCenterID:    runpod.PtrString(nv.DataCenterId),
				})
			}
		}
		state.NetworkVolumeIDs = bindings
	}

	return state
}
