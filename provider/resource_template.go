package provider

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/runpod/pulumi-runpod/pkg/runpod"
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
	StartSsh                *bool             `pulumi:"startSsh,optional"`
	StartScript             *string           `pulumi:"startScript,optional"`
	IsServerless            *bool             `pulumi:"isServerless,optional"`
	IsPublic                *bool             `pulumi:"isPublic,optional"`
	ContainerRegistryAuthID *string           `pulumi:"containerRegistryAuthId,optional"`
}

// TemplateState is the persisted state of a template resource.
type TemplateState struct {
	TemplateArgs
	TemplateID string `pulumi:"templateId"`
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

	saveInput := templateArgsToSaveInput("", input)
	tmpl, err := client.CreateTemplate(ctx, saveInput)
	if err != nil {
		return infer.CreateResponse[TemplateState]{}, err
	}

	state := templateToState(input, tmpl)
	return infer.CreateResponse[TemplateState]{
		ID:     tmpl.ID,
		Output: state,
	}, nil
}

// Read refreshes the template state from the API.
func (Template) Read(
	ctx context.Context,
	req infer.ReadRequest[TemplateArgs, TemplateState],
) (infer.ReadResponse[TemplateArgs, TemplateState], error) {
	client := getClient(ctx)

	tmpl, err := client.GetTemplate(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[TemplateArgs, TemplateState]{}, err
	}

	if tmpl == nil {
		return infer.ReadResponse[TemplateArgs, TemplateState]{},
			fmt.Errorf("template %q not found", req.ID)
	}

	state := templateToState(req.Inputs, tmpl)
	return infer.ReadResponse[TemplateArgs, TemplateState]{
		ID:     tmpl.ID,
		Inputs: req.Inputs,
		State:  state,
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

	saveInput := templateArgsToSaveInput(req.ID, req.Inputs)
	tmpl, err := client.UpdateTemplate(ctx, saveInput)
	if err != nil {
		return infer.UpdateResponse[TemplateState]{}, err
	}

	state := templateToState(req.Inputs, tmpl)
	return infer.UpdateResponse[TemplateState]{Output: state}, nil
}

// Delete removes a template.
func (Template) Delete(ctx context.Context, req infer.DeleteRequest[TemplateState]) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	if err := client.DeleteTemplate(ctx, req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func templateArgsToSaveInput(id string, args TemplateArgs) runpod.SaveTemplateInput {
	input := runpod.SaveTemplateInput{
		ID:                id,
		Name:              args.Name,
		ImageName:         args.ImageName,
		ContainerDiskInGb: args.ContainerDiskInGb,
		VolumeInGb:        args.VolumeInGb,
		Env:               runpod.EnvMapToGQL(args.Env),
		StartJupyter:      args.StartJupyter,
		StartSsh:          args.StartSsh,
		IsServerless:      args.IsServerless,
		IsPublic:          args.IsPublic,
	}
	if args.DockerArgs != nil {
		input.DockerArgs = *args.DockerArgs
	}
	if args.Ports != nil {
		input.Ports = *args.Ports
	}
	if args.VolumeMountPath != nil {
		input.VolumeMountPath = *args.VolumeMountPath
	}
	if args.StartScript != nil {
		input.StartScript = *args.StartScript
	}
	if args.ContainerRegistryAuthID != nil {
		input.ContainerRegistryAuthID = *args.ContainerRegistryAuthID
	}
	return input
}

func templateToState(input TemplateArgs, tmpl *runpod.PodTemplate) TemplateState {
	state := TemplateState{
		TemplateArgs: input,
		TemplateID:   tmpl.ID,
	}
	if env := runpod.EnvGQLToMap(tmpl.Env); len(env) > 0 {
		state.Env = env
	}
	return state
}
