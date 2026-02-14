package provider

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/runpod/pulumi-runpod/pkg/runpod"
)

// Pod is the controller for the runpod:index:Pod resource.
type Pod struct{}

// PodArgs are the inputs for creating a pod.
type PodArgs struct {
	Name                    string            `pulumi:"name"`
	GpuTypeID               string            `pulumi:"gpuTypeId"`
	GpuCount                *int              `pulumi:"gpuCount,optional"`
	CloudType               *string           `pulumi:"cloudType,optional"`
	ImageName               *string           `pulumi:"imageName,optional"`
	DockerArgs              *string           `pulumi:"dockerArgs,optional"`
	Env                     map[string]string `pulumi:"env,optional"`
	Ports                   *string           `pulumi:"ports,optional"`
	VolumeInGb              *int              `pulumi:"volumeInGb,optional"`
	VolumeMountPath         *string           `pulumi:"volumeMountPath,optional"`
	ContainerDiskInGb       *int              `pulumi:"containerDiskInGb,optional"`
	TemplateID              *string           `pulumi:"templateId,optional"`
	NetworkVolumeID         *string           `pulumi:"networkVolumeId,optional"`
	ContainerRegistryAuthID *string           `pulumi:"containerRegistryAuthId,optional"`
	DataCenterID            *string           `pulumi:"dataCenterId,optional"`
	StartJupyter            *bool             `pulumi:"startJupyter,optional"`
	StartSsh                *bool             `pulumi:"startSsh,optional"`
	SupportPublicIP         *bool             `pulumi:"supportPublicIp,optional"`
	MinVcpuCount            *int              `pulumi:"minVcpuCount,optional"`
	MinMemoryInGb           *int              `pulumi:"minMemoryInGb,optional"`
	CudaVersion             *string           `pulumi:"cudaVersion,optional"`
}

// PodState is the persisted state of a pod resource.
type PodState struct {
	PodArgs
	// Outputs
	PodID         string  `pulumi:"podId"`
	MachineID     string  `pulumi:"machineId"`
	CostPerHr     float64 `pulumi:"costPerHr"`
	DesiredStatus string  `pulumi:"desiredStatus"`
	VcpuCount     float64 `pulumi:"vcpuCount"`
	MemoryInGb    float64 `pulumi:"memoryInGb"`
}

// Create creates a new pod using podFindAndDeployOnDemand.
func (Pod) Create(
	ctx context.Context,
	req infer.CreateRequest[PodArgs],
) (infer.CreateResponse[PodState], error) {
	input := req.Inputs
	if req.DryRun {
		return infer.CreateResponse[PodState]{
			ID:     req.Name,
			Output: PodState{PodArgs: input},
		}, nil
	}

	client := getClient(ctx)

	createInput := runpod.PodFindAndDeployOnDemandInput{
		Name:                    &input.Name,
		GpuTypeId:               &input.GpuTypeID,
		GpuCount:                input.GpuCount,
		ImageName:               input.ImageName,
		DockerArgs:              input.DockerArgs,
		Env:                     runpod.EnvMapToGQL(input.Env),
		Ports:                   input.Ports,
		VolumeInGb:              input.VolumeInGb,
		VolumeMountPath:         input.VolumeMountPath,
		ContainerDiskInGb:       input.ContainerDiskInGb,
		TemplateId:              input.TemplateID,
		NetworkVolumeId:         input.NetworkVolumeID,
		ContainerRegistryAuthId: input.ContainerRegistryAuthID,
		DataCenterId:            input.DataCenterID,
		StartJupyter:            input.StartJupyter,
		StartSsh:                input.StartSsh,
		SupportPublicIp:         input.SupportPublicIP,
		MinVcpuCount:            input.MinVcpuCount,
		MinMemoryInGb:           input.MinMemoryInGb,
		CudaVersion:             input.CudaVersion,
	}

	// CloudType needs conversion to the enum
	if input.CloudType != nil {
		ct := runpod.CloudTypeEnum(*input.CloudType)
		createInput.CloudType = &ct
	}

	resp, err := runpod.CreatePod(ctx, client, createInput)
	if err != nil {
		return infer.CreateResponse[PodState]{}, err
	}

	if resp.PodFindAndDeployOnDemand == nil {
		return infer.CreateResponse[PodState]{}, fmt.Errorf("API returned nil pod")
	}

	pod := resp.PodFindAndDeployOnDemand
	state := podResponseToState(input, pod)
	return infer.CreateResponse[PodState]{
		ID:     pod.Id,
		Output: state,
	}, nil
}

// Read refreshes the pod state from the API.
func (Pod) Read(
	ctx context.Context,
	req infer.ReadRequest[PodArgs, PodState],
) (infer.ReadResponse[PodArgs, PodState], error) {
	client := getClient(ctx)

	resp, err := runpod.GetPod(ctx, client, runpod.PodFilter{PodId: req.ID})
	if err != nil {
		return infer.ReadResponse[PodArgs, PodState]{}, err
	}

	if resp.Pod == nil {
		return infer.ReadResponse[PodArgs, PodState]{}, fmt.Errorf("pod %q not found", req.ID)
	}

	state := podResponseToState(req.Inputs, resp.Pod)
	return infer.ReadResponse[PodArgs, PodState]{
		ID:     resp.Pod.Id,
		Inputs: req.Inputs,
		State:  state,
	}, nil
}

// Update modifies mutable pod fields (imageName, dockerArgs, env).
func (Pod) Update(
	ctx context.Context,
	req infer.UpdateRequest[PodArgs, PodState],
) (infer.UpdateResponse[PodState], error) {
	if req.DryRun {
		return infer.UpdateResponse[PodState]{
			Output: PodState{PodArgs: req.Inputs},
		}, nil
	}

	client := getClient(ctx)

	imageName := ""
	if req.Inputs.ImageName != nil {
		imageName = *req.Inputs.ImageName
	}

	updateInput := runpod.PodEditJobInput{
		PodId:     req.ID,
		ImageName: imageName,
		DockerArgs: req.Inputs.DockerArgs,
		Env:       runpod.EnvMapToGQL(req.Inputs.Env),
	}

	resp, err := runpod.UpdatePod(ctx, client, updateInput)
	if err != nil {
		return infer.UpdateResponse[PodState]{}, err
	}

	if resp.PodEditJob == nil {
		return infer.UpdateResponse[PodState]{}, fmt.Errorf("API returned nil pod on update")
	}

	state := podResponseToState(req.Inputs, resp.PodEditJob)
	return infer.UpdateResponse[PodState]{Output: state}, nil
}

// Delete terminates the pod.
func (Pod) Delete(ctx context.Context, req infer.DeleteRequest[PodState]) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	if _, err := runpod.TerminatePod(ctx, client, runpod.PodTerminateInput{PodId: req.ID}); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func podResponseToState(input PodArgs, pod *runpod.PodResponse) PodState {
	state := PodState{
		PodArgs:       input,
		PodID:         pod.Id,
		MachineID:     pod.MachineId,
		CostPerHr:     pod.CostPerHr,
		DesiredStatus: string(pod.DesiredStatus),
		VcpuCount:     pod.VcpuCount,
		MemoryInGb:    pod.MemoryInGb,
	}
	// Sync mutable fields from API response
	if pod.ImageName != nil {
		state.ImageName = pod.ImageName
	}
	if pod.DockerArgs != nil {
		state.DockerArgs = pod.DockerArgs
	}
	if env := runpod.EnvSliceToMap(pod.Env); len(env) > 0 {
		state.Env = env
	}
	return state
}
