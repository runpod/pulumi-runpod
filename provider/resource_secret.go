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

	"github.com/runpod/pulumi-runpod/pkg/runpod"
)

// Secret is the controller for the runpod:index:Secret resource.
type Secret struct{}

// SecretArgs are the inputs for creating a secret.
type SecretArgs struct {
	Name        string  `pulumi:"name" provider:"replaceOnChanges"`
	Value       string  `pulumi:"value" provider:"secret"`
	Description *string `pulumi:"description,optional"`
}

// Annotate provides descriptions for SecretArgs fields.
func (a *SecretArgs) Annotate(an infer.Annotator) {
	an.Describe(&a.Name, "A name for the secret.")
	an.Describe(&a.Value, "The secret value.")
	an.Describe(&a.Description,
		"A human-readable description of the secret.")
}

// SecretState is the persisted state of a secret resource.
type SecretState struct {
	SecretArgs
	SecretID string `pulumi:"secretId"`
}

// Annotate provides descriptions for SecretState fields.
func (s *SecretState) Annotate(a infer.Annotator) {
	a.Describe(&s.SecretID,
		"The unique identifier of the secret.")
}

// Create creates a new secret.
func (Secret) Create(
	ctx context.Context,
	req infer.CreateRequest[SecretArgs],
) (infer.CreateResponse[SecretState], error) {
	input := req.Inputs
	if req.DryRun {
		return infer.CreateResponse[SecretState]{
			ID:     req.Name,
			Output: SecretState{SecretArgs: input},
		}, nil
	}

	client := getClient(ctx)

	createInput := runpod.SecretCreateInput{
		Name:        input.Name,
		Value:       input.Value,
		Description: input.Description,
	}

	resp, err := runpod.SecretCreate(ctx, client, createInput)
	if err != nil {
		return infer.CreateResponse[SecretState]{}, err
	}

	if resp.SecretCreate == nil {
		return infer.CreateResponse[SecretState]{},
			errors.New("API returned nil secret")
	}

	secret := resp.SecretCreate
	state := secretResponseToState(input, secret)
	return infer.CreateResponse[SecretState]{
		ID:     secret.Id,
		Output: state,
	}, nil
}

// Read refreshes the secret state from the API.
func (Secret) Read(
	ctx context.Context,
	req infer.ReadRequest[SecretArgs, SecretState],
) (infer.ReadResponse[SecretArgs, SecretState], error) {
	client := getClient(ctx)

	resp, err := runpod.GetMySecrets(ctx, client)
	if err != nil {
		return infer.ReadResponse[SecretArgs, SecretState]{}, err
	}

	if resp.Myself == nil {
		return infer.ReadResponse[SecretArgs, SecretState]{
			ID: "",
		}, nil
	}

	for _, s := range resp.Myself.Secrets {
		if s.Id == req.ID {
			state := secretResponseToState(req.Inputs, &s)
			return infer.ReadResponse[SecretArgs, SecretState]{
				ID:     req.ID,
				Inputs: req.Inputs,
				State:  state,
			}, nil
		}
	}

	return infer.ReadResponse[SecretArgs, SecretState]{
		ID: "",
	}, nil
}

// Update modifies a secret's value and/or description.
func (Secret) Update(
	ctx context.Context,
	req infer.UpdateRequest[SecretArgs, SecretState],
) (infer.UpdateResponse[SecretState], error) {
	if req.DryRun {
		return infer.UpdateResponse[SecretState]{
			Output: SecretState{SecretArgs: req.Inputs},
		}, nil
	}

	client := getClient(ctx)

	// Update value
	if req.Inputs.Value != req.State.Value {
		_, err := runpod.SecretValueUpdate(
			ctx, client, runpod.SecretValueUpdateInput{
				Id:    req.ID,
				Value: req.Inputs.Value,
			},
		)
		if err != nil {
			return infer.UpdateResponse[SecretState]{}, err
		}
	}

	// Update description
	if req.Inputs.Description != nil {
		oldDesc := ""
		if req.State.Description != nil {
			oldDesc = *req.State.Description
		}
		if *req.Inputs.Description != oldDesc {
			_, err := runpod.SecretDescriptionUpdate(
				ctx, client, runpod.SecretDescriptionUpdateInput{
					Id:          req.ID,
					Description: *req.Inputs.Description,
				},
			)
			if err != nil {
				return infer.UpdateResponse[SecretState]{}, err
			}
		}
	}

	// Re-read to get latest state
	readResp, err := runpod.GetMySecrets(ctx, client)
	if err != nil {
		return infer.UpdateResponse[SecretState]{}, err
	}

	if readResp.Myself != nil {
		for _, s := range readResp.Myself.Secrets {
			if s.Id == req.ID {
				state := secretResponseToState(req.Inputs, &s)
				return infer.UpdateResponse[SecretState]{
					Output: state,
				}, nil
			}
		}
	}

	// Fallback: return input-based state
	return infer.UpdateResponse[SecretState]{
		Output: SecretState{
			SecretArgs: req.Inputs,
			SecretID:   req.ID,
		},
	}, nil
}

// Delete removes a secret.
func (Secret) Delete(
	ctx context.Context,
	req infer.DeleteRequest[SecretState],
) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	if _, err := runpod.SecretDelete(ctx, client, req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func secretResponseToState(
	input SecretArgs, secret *runpod.SecretResponse,
) SecretState {
	state := SecretState{
		SecretArgs: input,
		SecretID:   secret.Id,
	}
	// Sync name and description from API
	state.Name = secret.Name
	if secret.Description != nil {
		state.Description = secret.Description
	}
	return state
}
