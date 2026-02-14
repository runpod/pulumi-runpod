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

	createInput := runpod.CreatePodInput{
		Name:            input.Name,
		GpuTypeID:       input.GpuTypeID,
		Env:             runpod.EnvMapToGQL(input.Env),
		StartJupyter:    input.StartJupyter,
		StartSsh:        input.StartSsh,
		SupportPublicIP: input.SupportPublicIP,
	}
	if input.GpuCount != nil {
		createInput.GpuCount = *input.GpuCount
	}
	if input.CloudType != nil {
		createInput.CloudType = *input.CloudType
	}
	if input.ImageName != nil {
		createInput.ImageName = *input.ImageName
	}
	if input.DockerArgs != nil {
		createInput.DockerArgs = *input.DockerArgs
	}
	if input.Ports != nil {
		createInput.Ports = *input.Ports
	}
	if input.VolumeInGb != nil {
		createInput.VolumeInGb = *input.VolumeInGb
	}
	if input.VolumeMountPath != nil {
		createInput.VolumeMountPath = *input.VolumeMountPath
	}
	if input.ContainerDiskInGb != nil {
		createInput.ContainerDiskInGb = *input.ContainerDiskInGb
	}
	if input.TemplateID != nil {
		createInput.TemplateID = *input.TemplateID
	}
	if input.NetworkVolumeID != nil {
		createInput.NetworkVolumeID = *input.NetworkVolumeID
	}
	if input.ContainerRegistryAuthID != nil {
		createInput.ContainerRegistryAuthID = *input.ContainerRegistryAuthID
	}
	if input.DataCenterID != nil {
		createInput.DataCenterID = *input.DataCenterID
	}
	if input.MinVcpuCount != nil {
		createInput.MinVcpuCount = *input.MinVcpuCount
	}
	if input.MinMemoryInGb != nil {
		createInput.MinMemoryInGb = *input.MinMemoryInGb
	}
	if input.CudaVersion != nil {
		createInput.CudaVersion = *input.CudaVersion
	}

	pod, err := client.CreatePod(ctx, createInput)
	if err != nil {
		return infer.CreateResponse[PodState]{}, err
	}

	state := podToState(input, pod)
	return infer.CreateResponse[PodState]{
		ID:     pod.ID,
		Output: state,
	}, nil
}

// Read refreshes the pod state from the API.
func (Pod) Read(
	ctx context.Context,
	req infer.ReadRequest[PodArgs, PodState],
) (infer.ReadResponse[PodArgs, PodState], error) {
	client := getClient(ctx)

	pod, err := client.GetPod(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[PodArgs, PodState]{}, err
	}

	if pod == nil {
		return infer.ReadResponse[PodArgs, PodState]{}, fmt.Errorf("pod %q not found", req.ID)
	}

	state := podToState(req.Inputs, pod)
	return infer.ReadResponse[PodArgs, PodState]{
		ID:     pod.ID,
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

	updateInput := runpod.UpdatePodInput{
		PodID: req.ID,
		Env:   runpod.EnvMapToGQL(req.Inputs.Env),
	}
	if req.Inputs.ImageName != nil {
		updateInput.ImageName = *req.Inputs.ImageName
	}
	if req.Inputs.DockerArgs != nil {
		updateInput.DockerArgs = *req.Inputs.DockerArgs
	}

	pod, err := client.UpdatePod(ctx, updateInput)
	if err != nil {
		return infer.UpdateResponse[PodState]{}, err
	}

	state := podToState(req.Inputs, pod)
	return infer.UpdateResponse[PodState]{Output: state}, nil
}

// Delete terminates the pod.
func (Pod) Delete(ctx context.Context, req infer.DeleteRequest[PodState]) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	if err := client.TerminatePod(ctx, req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func podToState(input PodArgs, pod *runpod.Pod) PodState {
	state := PodState{
		PodArgs:       input,
		PodID:         pod.ID,
		MachineID:     pod.MachineID,
		CostPerHr:     pod.CostPerHr,
		DesiredStatus: pod.DesiredStatus,
		VcpuCount:     pod.VcpuCount,
		MemoryInGb:    pod.MemoryInGb,
	}
	// Sync mutable fields from API response
	if pod.ImageName != "" {
		state.ImageName = &pod.ImageName
	}
	if pod.DockerArgs != "" {
		state.DockerArgs = &pod.DockerArgs
	}
	if env := runpod.EnvSliceToMap(pod.Env); len(env) > 0 {
		state.Env = env
	}
	return state
}
