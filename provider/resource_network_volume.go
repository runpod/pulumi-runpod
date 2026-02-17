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
	"errors"

	"github.com/runpod/pulumi-runpod/pkg/runpod"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// NetworkVolume is the controller for the runpod:index:NetworkVolume resource.
type NetworkVolume struct{}

// NetworkVolumeArgs are the inputs for creating a network volume.
type NetworkVolumeArgs struct {
	Name             string `pulumi:"name"`
	Size             int    `pulumi:"size"`
	DataCenterID     string `pulumi:"dataCenterId" provider:"replaceOnChanges"`
	IsNextGenStorage *bool  `pulumi:"isNextGenStorage,optional" provider:"replaceOnChanges"`
}

// Annotate provides descriptions for NetworkVolumeArgs fields.
func (a *NetworkVolumeArgs) Annotate(an infer.Annotator) {
	an.Describe(&a.Name,
		"A name for the network volume.")
	an.Describe(&a.Size,
		"The size of the network volume in GB.")
	an.Describe(&a.DataCenterID,
		"The data center ID where the volume will be created "+
			"(e.g. \"US-TX-3\").")
	an.Describe(&a.IsNextGenStorage,
		"Whether to use next-generation storage.")
}

// NetworkVolumeState is the persisted state of a network volume resource.
type NetworkVolumeState struct {
	NetworkVolumeArgs
	NetworkVolumeID string `pulumi:"networkVolumeId"`
}

// Annotate provides descriptions for NetworkVolumeState fields.
func (s *NetworkVolumeState) Annotate(a infer.Annotator) {
	a.Describe(&s.NetworkVolumeID,
		"The unique identifier of the network volume.")
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
		Name:             &input.Name,
		Size:             &input.Size,
		DataCenterId:     input.DataCenterID,
		IsNextGenStorage: input.IsNextGenStorage,
	}

	resp, err := runpod.CreateNetworkVolume(ctx, client, createInput)
	if err != nil {
		return infer.CreateResponse[NetworkVolumeState]{}, err
	}

	if resp.CreateNetworkVolume == nil {
		return infer.CreateResponse[NetworkVolumeState]{},
			errors.New("API returned nil network volume")
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
		return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{},
			err
	}

	if resp.Myself == nil {
		return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{
			ID: "",
		}, nil
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

	return infer.ReadResponse[NetworkVolumeArgs, NetworkVolumeState]{
		ID: "",
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
		Id:   req.ID,
		Name: &req.Inputs.Name,
		Size: &req.Inputs.Size,
	}

	resp, err := runpod.UpdateNetworkVolume(ctx, client, updateInput)
	if err != nil {
		return infer.UpdateResponse[NetworkVolumeState]{}, err
	}

	if resp.UpdateNetworkVolume == nil {
		return infer.UpdateResponse[NetworkVolumeState]{},
			errors.New("API returned nil network volume on update")
	}

	state := networkVolumeResponseToState(
		req.Inputs, resp.UpdateNetworkVolume,
	)
	return infer.UpdateResponse[NetworkVolumeState]{Output: state}, nil
}

// Delete removes a network volume.
func (NetworkVolume) Delete(
	ctx context.Context,
	req infer.DeleteRequest[NetworkVolumeState],
) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	_, err := runpod.DeleteNetworkVolume(
		ctx, client, runpod.DeleteNetworkVolumeInput{Id: req.ID},
	)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func networkVolumeResponseToState(
	input NetworkVolumeArgs, vol *runpod.NetworkVolumeResponse,
) NetworkVolumeState {
	return NetworkVolumeState{
		NetworkVolumeArgs: input,
		NetworkVolumeID:   runpod.PtrString(vol.Id),
	}
}
