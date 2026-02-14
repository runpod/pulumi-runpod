package provider

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/runpod/pulumi-runpod/pkg/runpod"
)

// NetworkVolume is the controller for the runpod:index:NetworkVolume resource.
type NetworkVolume struct{}

// NetworkVolumeArgs are the inputs for creating a network volume.
type NetworkVolumeArgs struct {
	Name         string `pulumi:"name"`
	Size         int    `pulumi:"size"`
	DataCenterID string `pulumi:"dataCenterId"`
}

// NetworkVolumeState is the persisted state of a network volume resource.
type NetworkVolumeState struct {
	NetworkVolumeArgs
	NetworkVolumeID string `pulumi:"networkVolumeId"`
}

// Create creates a new network volume.
func (NetworkVolume) Create(
	ctx context.Context,
	req infer.CreateRequest[NetworkVolumeArgs],
) (infer.CreateResponse[NetworkVolumeState], error) {
	input := req.Inputs
	if req.DryRun {
		return infer.CreateResponse[NetworkVolumeState]{
			ID:     req.Name,
			Output: NetworkVolumeState{NetworkVolumeArgs: input},
		}, nil
	}

	client := getClient(ctx)

	createInput := runpod.CreateNetworkVolumeInput{
		Name:         input.Name,
		Size:         input.Size,
		DataCenterID: input.DataCenterID,
	}

	vol, err := client.CreateNetworkVolume(ctx, createInput)
	if err != nil {
		return infer.CreateResponse[NetworkVolumeState]{}, err
	}

	state := networkVolumeToState(input, vol)
	return infer.CreateResponse[NetworkVolumeState]{
		ID:     vol.ID,
		Output: state,
	}, nil
}

// Read refreshes the network volume state from the API.
func (NetworkVolume) Read(
	ctx context.Context,
	req infer.ReadRequest[NetworkVolumeArgs, NetworkVolumeState],
) (infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState], error) {
	client := getClient(ctx)

	vol, err := client.GetNetworkVolume(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{}, err
	}

	if vol == nil {
		return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{},
			fmt.Errorf("network volume %q not found", req.ID)
	}

	state := networkVolumeToState(req.Inputs, vol)
	return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{
		ID:     vol.ID,
		Inputs: req.Inputs,
		State:  state,
	}, nil
}

// Update modifies a network volume (name and size are mutable).
func (NetworkVolume) Update(
	ctx context.Context,
	req infer.UpdateRequest[NetworkVolumeArgs, NetworkVolumeState],
) (infer.UpdateResponse[NetworkVolumeState], error) {
	if req.DryRun {
		return infer.UpdateResponse[NetworkVolumeState]{
			Output: NetworkVolumeState{NetworkVolumeArgs: req.Inputs},
		}, nil
	}

	client := getClient(ctx)

	updateInput := runpod.UpdateNetworkVolumeInput{
		ID:   req.ID,
		Name: req.Inputs.Name,
		Size: req.Inputs.Size,
	}

	vol, err := client.UpdateNetworkVolume(ctx, updateInput)
	if err != nil {
		return infer.UpdateResponse[NetworkVolumeState]{}, err
	}

	state := networkVolumeToState(req.Inputs, vol)
	return infer.UpdateResponse[NetworkVolumeState]{Output: state}, nil
}

// Delete removes a network volume.
func (NetworkVolume) Delete(ctx context.Context, req infer.DeleteRequest[NetworkVolumeState]) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	if err := client.DeleteNetworkVolume(ctx, req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func networkVolumeToState(input NetworkVolumeArgs, vol *runpod.NetworkVolume) NetworkVolumeState {
	return NetworkVolumeState{
		NetworkVolumeArgs: input,
		NetworkVolumeID:   vol.ID,
	}
}
