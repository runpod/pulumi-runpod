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
		return infer.ReadResponse[EndpointArgs, EndpointState]{},
			fmt.Errorf("endpoint %q not found", req.ID)
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

	return infer.ReadResponse[EndpointArgs, EndpointState]{},
		fmt.Errorf("endpoint %q not found", req.ID)
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
func (Endpoint) Delete(ctx context.Context, req infer.DeleteRequest[EndpointState]) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	if _, err := runpod.DeleteEndpoint(ctx, client, req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func endpointArgsToInput(id *string, args EndpointArgs) runpod.EndpointInput {
	input := runpod.EndpointInput{
		Id:              id,
		Name:            args.Name,
		TemplateId:      args.TemplateID,
		GpuIds:          args.GpuIds,
		WorkersMin:      args.WorkersMin,
		WorkersMax:      args.WorkersMax,
		IdleTimeout:     args.IdleTimeout,
		Locations:       args.Locations,
		ScalerType:      args.ScalerType,
		ScalerValue:     args.ScalerValue,
		NetworkVolumeId: args.NetworkVolumeID,
		GpuCount:        args.GpuCount,
		InstanceIds:     runpod.StringPtrSlice(args.InstanceIds),
		Env:             runpod.EnvMapToGQL(args.Env),
	}
	return input
}

func endpointResponseToState(input EndpointArgs, ep *runpod.EndpointResponse) EndpointState {
	return EndpointState{
		EndpointArgs: input,
		EndpointID:   runpod.PtrString(ep.Id),
	}
}
