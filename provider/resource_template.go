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

// Template is the controller for the runpod:index:Template resource.
type Template struct{}

// TemplateArgs are the inputs for creating a template.
type TemplateArgs struct {
	Name                    string            `pulumi:"name"`
	ImageName               string            `pulumi:"imageName"`
	ContainerDiskInGb       int               `pulumi:"containerDiskInGb"`
	VolumeInGb              int               `pulumi:"volumeInGb"`
	DockerArgs              *string           `pulumi:"dockerArgs,optional"`
	Env                     map[string]string `pulumi:"env,optional"`
	Ports                   *string           `pulumi:"ports,optional"`
	VolumeMountPath         *string           `pulumi:"volumeMountPath,optional"`
	StartJupyter            *bool             `pulumi:"startJupyter,optional"`
	StartSSH                *bool             `pulumi:"startSsh,optional"`
	StartScript             *string           `pulumi:"startScript,optional"`
	IsServerless            *bool             `pulumi:"isServerless,optional"`
	IsPublic                *bool             `pulumi:"isPublic,optional"`
	ContainerRegistryAuthID *string           `pulumi:"containerRegistryAuthId,optional"`
	// New fields
	Readme        *string `pulumi:"readme,optional"`
	AdvancedStart *bool   `pulumi:"advancedStart,optional"`
	Category      *string `pulumi:"category,optional"`
}

// Annotate provides descriptions for TemplateArgs fields.
func (a *TemplateArgs) Annotate(an infer.Annotator) {
	an.Describe(&a.Name, "A name for the template.")
	an.Describe(&a.ImageName,
		"The Docker image to use for the template.")
	an.Describe(&a.ContainerDiskInGb,
		"The size of the container disk in GB.")
	an.Describe(&a.VolumeInGb,
		"The size of the persistent volume in GB. Use 0 for no volume.")
	an.Describe(&a.DockerArgs,
		"Docker arguments to pass to the container.")
	an.Describe(&a.Env,
		"Environment variables as key-value pairs.")
	an.Describe(&a.Ports,
		"Ports to expose (e.g. \"8080/http,22/tcp\").")
	an.Describe(&a.VolumeMountPath,
		"The path to mount the persistent volume.")
	an.Describe(&a.StartJupyter,
		"Whether to start Jupyter notebook server.")
	an.Describe(&a.StartSSH,
		"Whether to start an SSH server.")
	an.Describe(&a.StartScript,
		"A bash script to run on container start.")
	an.Describe(&a.IsServerless,
		"Whether this template is for serverless endpoints.")
	an.Describe(&a.IsPublic,
		"Whether this template is publicly visible.")
	an.Describe(&a.ContainerRegistryAuthID,
		"The ID of the container registry auth credentials to use.")
	an.Describe(&a.Readme,
		"A readme/description for the template in Markdown.")
	an.Describe(&a.AdvancedStart,
		"Whether to use advanced start mode.")
	an.Describe(&a.Category,
		"The category of the template.")
}

// TemplateState is the persisted state of a template resource.
type TemplateState struct {
	TemplateArgs
	TemplateID string `pulumi:"templateId"`
}

// Annotate provides descriptions for TemplateState fields.
func (s *TemplateState) Annotate(a infer.Annotator) {
	a.Describe(&s.TemplateID,
		"The unique identifier of the template.")
}

// Create creates a new template.
func (Template) Create(
	ctx context.Context,
	req infer.CreateRequest[TemplateArgs],
) (infer.CreateResponse[TemplateState], error) {
	input := req.Inputs
	if req.DryRun {
		return infer.CreateResponse[TemplateState]{
			ID:     req.Name,
			Output: TemplateState{TemplateArgs: input},
		}, nil
	}

	client := getClient(ctx)
	saveInput := templateArgsToSaveInput(nil, input)

	resp, err := runpod.SaveTemplate(ctx, client, &saveInput)
	if err != nil {
		return infer.CreateResponse[TemplateState]{}, err
	}

	if resp.SaveTemplate == nil {
		return infer.CreateResponse[TemplateState]{},
			errors.New("API returned nil template")
	}

	state := templateResponseToState(input, resp.SaveTemplate)
	return infer.CreateResponse[TemplateState]{
		ID:     runpod.PtrString(resp.SaveTemplate.Id),
		Output: state,
	}, nil
}

// Read refreshes the template state from the API.
func (Template) Read(
	ctx context.Context,
	req infer.ReadRequest[TemplateArgs, TemplateState],
) (infer.ReadResponse[TemplateArgs, TemplateState], error) {
	client := getClient(ctx)

	resp, err := runpod.GetMyTemplates(ctx, client)
	if err != nil {
		return infer.ReadResponse[TemplateArgs, TemplateState]{}, err
	}

	if resp.Myself == nil {
		return infer.ReadResponse[TemplateArgs, TemplateState]{
			ID: "",
		}, nil
	}

	for _, t := range resp.Myself.PodTemplates {
		if t != nil && runpod.PtrString(t.Id) == req.ID {
			state := templateResponseToState(req.Inputs, t)
			return infer.ReadResponse[TemplateArgs, TemplateState]{
				ID:     req.ID,
				Inputs: req.Inputs,
				State:  state,
			}, nil
		}
	}

	// Resource was deleted externally — return empty ID so Pulumi removes it from state.
	return infer.ReadResponse[TemplateArgs, TemplateState]{
		ID: "",
	}, nil
}

// Update modifies a template using the upsert pattern (saveTemplate with id).
func (Template) Update(
	ctx context.Context,
	req infer.UpdateRequest[TemplateArgs, TemplateState],
) (infer.UpdateResponse[TemplateState], error) {
	if req.DryRun {
		return infer.UpdateResponse[TemplateState]{
			Output: TemplateState{TemplateArgs: req.Inputs},
		}, nil
	}

	client := getClient(ctx)
	id := req.ID
	saveInput := templateArgsToSaveInput(&id, req.Inputs)

	resp, err := runpod.SaveTemplate(ctx, client, &saveInput)
	if err != nil {
		return infer.UpdateResponse[TemplateState]{}, err
	}

	if resp.SaveTemplate == nil {
		return infer.UpdateResponse[TemplateState]{},
			errors.New("API returned nil template")
	}

	state := templateResponseToState(req.Inputs, resp.SaveTemplate)
	return infer.UpdateResponse[TemplateState]{Output: state}, nil
}

// Delete removes a template.
// Note: RunPod's deleteTemplate mutation takes the template name, not ID.
func (Template) Delete(
	ctx context.Context,
	req infer.DeleteRequest[TemplateState],
) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	name := req.State.Name
	if _, err := runpod.DeleteTemplate(ctx, client, &name); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func templateArgsToSaveInput(
	id *string, args TemplateArgs,
) runpod.SaveTemplateInput {
	dockerArgs := ""
	if args.DockerArgs != nil {
		dockerArgs = *args.DockerArgs
	}

	input := runpod.SaveTemplateInput{
		Id:                      id,
		Name:                    args.Name,
		ImageName:               &args.ImageName,
		ContainerDiskInGb:       args.ContainerDiskInGb,
		VolumeInGb:              args.VolumeInGb,
		DockerArgs:              dockerArgs,
		Env:                     runpod.EnvMapToGQL(args.Env),
		Ports:                   args.Ports,
		VolumeMountPath:         args.VolumeMountPath,
		StartJupyter:            args.StartJupyter,
		StartSsh:                args.StartSSH,
		StartScript:             args.StartScript,
		IsServerless:            args.IsServerless,
		IsPublic:                args.IsPublic,
		ContainerRegistryAuthId: args.ContainerRegistryAuthID,
		Readme:                  args.Readme,
		AdvancedStart:           args.AdvancedStart,
	}

	if args.Category != nil {
		cat := runpod.TemplateCategory(*args.Category)
		input.Category = &cat
	}

	return input
}

func templateResponseToState(
	input TemplateArgs, tmpl *runpod.TemplateResponse,
) TemplateState {
	state := TemplateState{
		TemplateArgs: input,
		TemplateID:   runpod.PtrString(tmpl.Id),
	}
	if env := runpod.EnvGQLToMap(tmpl.Env); len(env) > 0 {
		state.Env = env
	}
	return state
}
