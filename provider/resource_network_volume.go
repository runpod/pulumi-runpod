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
		Name:         &input.Name,
		Size:         &input.Size,
		DataCenterId: input.DataCenterID,
	}

	resp, err := runpod.CreateNetworkVolume(ctx, client, createInput)
	if err != nil {
		return infer.CreateResponse[NetworkVolumeState]{}, err
	}

	if resp.CreateNetworkVolume == nil {
		return infer.CreateResponse[NetworkVolumeState]{}, fmt.Errorf("API returned nil network volume")
	}

	vol := resp.CreateNetworkVolume
	state := networkVolumeResponseToState(input, vol)
	return infer.CreateResponse[NetworkVolumeState]{
		ID:     runpod.PtrString(vol.Id),
		Output: state,
	}, nil
}

// Read refreshes the network volume state from the API.
func (NetworkVolume) Read(
	ctx context.Context,
	req infer.ReadRequest[NetworkVolumeArgs, NetworkVolumeState],
) (infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState], error) {
	client := getClient(ctx)

	resp, err := runpod.GetMyNetworkVolumes(ctx, client)
	if err != nil {
		return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{}, err
	}

	if resp.Myself == nil {
		return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{},
			fmt.Errorf("network volume %q not found", req.ID)
	}

	for _, v := range resp.Myself.NetworkVolumes {
		if v != nil && runpod.PtrString(v.Id) == req.ID {
			state := networkVolumeResponseToState(req.Inputs, v)
			return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{
				ID:     req.ID,
				Inputs: req.Inputs,
				State:  state,
			}, nil
		}
	}

	return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{},
		fmt.Errorf("network volume %q not found", req.ID)
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
		Id:   req.ID,
		Name: &req.Inputs.Name,
		Size: &req.Inputs.Size,
	}

	resp, err := runpod.UpdateNetworkVolume(ctx, client, updateInput)
	if err != nil {
		return infer.UpdateResponse[NetworkVolumeState]{}, err
	}

	if resp.UpdateNetworkVolume == nil {
		return infer.UpdateResponse[NetworkVolumeState]{}, fmt.Errorf("API returned nil network volume on update")
	}

	state := networkVolumeResponseToState(req.Inputs, resp.UpdateNetworkVolume)
	return infer.UpdateResponse[NetworkVolumeState]{Output: state}, nil
}

// Delete removes a network volume.
func (NetworkVolume) Delete(ctx context.Context, req infer.DeleteRequest[NetworkVolumeState]) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	if _, err := runpod.DeleteNetworkVolume(ctx, client, runpod.DeleteNetworkVolumeInput{Id: req.ID}); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func networkVolumeResponseToState(input NetworkVolumeArgs, vol *runpod.NetworkVolumeResponse) NetworkVolumeState {
	return NetworkVolumeState{
		NetworkVolumeArgs: input,
		NetworkVolumeID:   runpod.PtrString(vol.Id),
	}
}
