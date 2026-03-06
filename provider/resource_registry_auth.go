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

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/runpod/pulumi-runpod/provider/pkg/runpod"
)

// ContainerRegistryAuth is the controller for the runpod:index:ContainerRegistryAuth resource.
type ContainerRegistryAuth struct{}

// ContainerRegistryAuthArgs are the inputs for creating a container registry auth.
type ContainerRegistryAuthArgs struct {
	Name     string `pulumi:"name" provider:"replaceOnChanges"`
	Username string `pulumi:"username"`
	Password string `pulumi:"password" provider:"secret"`
}

// Annotate provides descriptions for ContainerRegistryAuthArgs fields.
func (a *ContainerRegistryAuthArgs) Annotate(an infer.Annotator) {
	an.Describe(&a.Name,
		"A name for the registry auth credentials.")
	an.Describe(&a.Username,
		"The username for the container registry.")
	an.Describe(&a.Password,
		"The password or access token for the container registry.")
}

// ContainerRegistryAuthState is the persisted state of a container registry auth resource.
type ContainerRegistryAuthState struct {
	ContainerRegistryAuthArgs
	RegistryAuthID string `pulumi:"registryAuthId"`
}

// Annotate provides descriptions for ContainerRegistryAuthState fields.
func (s *ContainerRegistryAuthState) Annotate(a infer.Annotator) {
	a.Describe(&s.RegistryAuthID,
		"The unique identifier of the registry auth.")
}

// Create creates a new container registry auth.
func (ContainerRegistryAuth) Create(
	ctx context.Context,
	req infer.CreateRequest[ContainerRegistryAuthArgs],
) (infer.CreateResponse[ContainerRegistryAuthState], error) {
	input := req.Inputs
	if req.DryRun {
		return infer.CreateResponse[ContainerRegistryAuthState]{
			ID: req.Name,
			Output: ContainerRegistryAuthState{
				ContainerRegistryAuthArgs: input,
			},
		}, nil
	}

	client := getClient(ctx)

	saveInput := &runpod.SaveRegistryAuthInput{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
	}

	resp, err := runpod.SaveRegistryAuth(ctx, client, saveInput)
	if err != nil {
		return infer.CreateResponse[ContainerRegistryAuthState]{}, err
	}

	if resp.SaveRegistryAuth == nil {
		return infer.CreateResponse[ContainerRegistryAuthState]{},
			errors.New("API returned nil registry auth")
	}

	ra := resp.SaveRegistryAuth
	state := registryAuthResponseToState(input, ra)
	return infer.CreateResponse[ContainerRegistryAuthState]{
		ID:     runpod.PtrString(ra.Id),
		Output: state,
	}, nil
}

// Read refreshes the registry auth state from the API.
func (ContainerRegistryAuth) Read(
	ctx context.Context,
	req infer.ReadRequest[
		ContainerRegistryAuthArgs, ContainerRegistryAuthState,
	],
) (infer.ReadResponse[
	ContainerRegistryAuthArgs, ContainerRegistryAuthState,
], error,
) {
	client := getClient(ctx)

	resp, err := runpod.GetMyRegistryAuths(ctx, client)
	if err != nil {
		return infer.ReadResponse[
			ContainerRegistryAuthArgs, ContainerRegistryAuthState,
		]{}, err
	}

	if resp.Myself == nil {
		return infer.ReadResponse[
			ContainerRegistryAuthArgs, ContainerRegistryAuthState,
		]{ID: ""}, nil
	}

	for _, ra := range resp.Myself.ContainerRegistryCreds {
		if ra != nil && runpod.PtrString(ra.Id) == req.ID {
			state := registryAuthResponseToState(req.Inputs, ra)
			return infer.ReadResponse[
				ContainerRegistryAuthArgs, ContainerRegistryAuthState,
			]{
				ID:     req.ID,
				Inputs: req.Inputs,
				State:  state,
			}, nil
		}
	}

	return infer.ReadResponse[
		ContainerRegistryAuthArgs, ContainerRegistryAuthState,
	]{ID: ""}, nil
}

// Update modifies a container registry auth's credentials.
func (ContainerRegistryAuth) Update(
	ctx context.Context,
	req infer.UpdateRequest[
		ContainerRegistryAuthArgs, ContainerRegistryAuthState,
	],
) (infer.UpdateResponse[ContainerRegistryAuthState], error) {
	if req.DryRun {
		return infer.UpdateResponse[ContainerRegistryAuthState]{
			Output: ContainerRegistryAuthState{
				ContainerRegistryAuthArgs: req.Inputs,
			},
		}, nil
	}

	client := getClient(ctx)

	updateInput := &runpod.UpdateRegistryAuthInput{
		Id:       req.ID,
		Username: req.Inputs.Username,
		Password: req.Inputs.Password,
	}

	resp, err := runpod.UpdateRegistryAuth(ctx, client, updateInput)
	if err != nil {
		return infer.UpdateResponse[ContainerRegistryAuthState]{}, err
	}

	if resp.UpdateRegistryAuth == nil {
		return infer.UpdateResponse[ContainerRegistryAuthState]{},
			errors.New("API returned nil registry auth on update")
	}

	state := registryAuthResponseToState(
		req.Inputs, resp.UpdateRegistryAuth,
	)
	return infer.UpdateResponse[ContainerRegistryAuthState]{
		Output: state,
	}, nil
}

// Delete removes a container registry auth.
func (ContainerRegistryAuth) Delete(
	ctx context.Context,
	req infer.DeleteRequest[ContainerRegistryAuthState],
) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	id := req.ID
	_, err := runpod.DeleteRegistryAuth(ctx, client, &id)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func registryAuthResponseToState(
	input ContainerRegistryAuthArgs,
	ra *runpod.RegistryAuthResponse,
) ContainerRegistryAuthState {
	state := ContainerRegistryAuthState{
		ContainerRegistryAuthArgs: input,
		RegistryAuthID:            runpod.PtrString(ra.Id),
	}
	if ra.Name != nil {
		state.Name = *ra.Name
	}
	return state
}
