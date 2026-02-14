package provider

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/runpod/pulumi-runpod/pkg/runpod"
)

// Endpoint is the controller for the runpod:index:Endpoint resource.
type Endpoint struct{}

// EndpointArgs are the inputs for creating a serverless endpoint.
type EndpointArgs struct {
	Name            string            `pulumi:"name"`
	TemplateID      *string           `pulumi:"templateId,optional"`
	GpuIds          *string           `pulumi:"gpuIds,optional"`
	WorkersMin      *int              `pulumi:"workersMin,optional"`
	WorkersMax      *int              `pulumi:"workersMax,optional"`
	IdleTimeout     *int              `pulumi:"idleTimeout,optional"`
	Locations       *string           `pulumi:"locations,optional"`
	ScalerType      *string           `pulumi:"scalerType,optional"`
	ScalerValue     *int              `pulumi:"scalerValue,optional"`
	NetworkVolumeID *string           `pulumi:"networkVolumeId,optional"`
	GpuCount        *int              `pulumi:"gpuCount,optional"`
	InstanceIds     []string          `pulumi:"instanceIds,optional"`
	Env             map[string]string `pulumi:"env,optional"`
}

// EndpointState is the persisted state of an endpoint resource.
type EndpointState struct {
	EndpointArgs
	EndpointID string `pulumi:"endpointId"`
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

	saveInput := endpointArgsToSaveInput("", input)
	ep, err := client.CreateEndpoint(ctx, saveInput)
	if err != nil {
		return infer.CreateResponse[EndpointState]{}, err
	}

	state := endpointToState(input, ep)
	return infer.CreateResponse[EndpointState]{
		ID:     ep.ID,
		Output: state,
	}, nil
}

// Read refreshes the endpoint state from the API.
func (Endpoint) Read(
	ctx context.Context,
	req infer.ReadRequest[EndpointArgs, EndpointState],
) (infer.ReadResponse[EndpointArgs, EndpointState], error) {
	client := getClient(ctx)

	ep, err := client.GetEndpoint(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[EndpointArgs, EndpointState]{}, err
	}

	if ep == nil {
		return infer.ReadResponse[EndpointArgs, EndpointState]{},
			fmt.Errorf("endpoint %q not found", req.ID)
	}

	state := endpointToState(req.Inputs, ep)
	return infer.ReadResponse[EndpointArgs, EndpointState]{
		ID:     ep.ID,
		Inputs: req.Inputs,
		State:  state,
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

	saveInput := endpointArgsToSaveInput(req.ID, req.Inputs)
	ep, err := client.UpdateEndpoint(ctx, saveInput)
	if err != nil {
		return infer.UpdateResponse[EndpointState]{}, err
	}

	state := endpointToState(req.Inputs, ep)
	return infer.UpdateResponse[EndpointState]{Output: state}, nil
}

// Delete removes an endpoint.
func (Endpoint) Delete(ctx context.Context, req infer.DeleteRequest[EndpointState]) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	if err := client.DeleteEndpoint(ctx, req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func endpointArgsToSaveInput(id string, args EndpointArgs) runpod.SaveEndpointInput {
	input := runpod.SaveEndpointInput{
		ID:          id,
		Name:        args.Name,
		InstanceIds: args.InstanceIds,
		Env:         runpod.EnvMapToGQL(args.Env),
	}
	if args.TemplateID != nil {
		input.TemplateID = *args.TemplateID
	}
	if args.GpuIds != nil {
		input.GpuIds = *args.GpuIds
	}
	if args.WorkersMin != nil {
		input.WorkersMin = *args.WorkersMin
	}
	if args.WorkersMax != nil {
		input.WorkersMax = *args.WorkersMax
	}
	if args.IdleTimeout != nil {
		input.IdleTimeout = *args.IdleTimeout
	}
	if args.Locations != nil {
		input.Locations = *args.Locations
	}
	if args.ScalerType != nil {
		input.ScalerType = *args.ScalerType
	}
	if args.ScalerValue != nil {
		input.ScalerValue = *args.ScalerValue
	}
	if args.NetworkVolumeID != nil {
		input.NetworkVolumeID = *args.NetworkVolumeID
	}
	if args.GpuCount != nil {
		input.GpuCount = *args.GpuCount
	}
	return input
}

func endpointToState(input EndpointArgs, ep *runpod.Endpoint) EndpointState {
	return EndpointState{
		EndpointArgs: input,
		EndpointID:   ep.ID,
	}
}
