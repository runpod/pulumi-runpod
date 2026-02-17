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

// Pod is the controller for the runpod:index:Pod resource.
type Pod struct{}

// SavingsPlan represents a savings plan configuration for a pod.
type SavingsPlan struct {
	PlanLength  *string  `pulumi:"planLength,optional"`
	UpfrontCost *float64 `pulumi:"upfrontCost,optional"`
}

// PodArgs are the inputs for creating a pod.
// Fields tagged with replaceOnChanges are immutable — changing them requires replacing the pod.
type PodArgs struct {
	// Immutable fields (not in PodEditJobInput — require replacement)
	Name            string  `pulumi:"name" provider:"replaceOnChanges"`
	GpuTypeID       string  `pulumi:"gpuTypeId" provider:"replaceOnChanges"`
	GpuCount        *int    `pulumi:"gpuCount,optional" provider:"replaceOnChanges"`
	CloudType       *string `pulumi:"cloudType,optional" provider:"replaceOnChanges"`
	TemplateID      *string `pulumi:"templateId,optional" provider:"replaceOnChanges"`
	NetworkVolumeID *string `pulumi:"networkVolumeId,optional" provider:"replaceOnChanges"`
	DataCenterID    *string `pulumi:"dataCenterId,optional" provider:"replaceOnChanges"`
	StartJupyter    *bool   `pulumi:"startJupyter,optional" provider:"replaceOnChanges"`
	StartSSH        *bool   `pulumi:"startSsh,optional" provider:"replaceOnChanges"`
	SupportPublicIP *bool   `pulumi:"supportPublicIp,optional" provider:"replaceOnChanges"`
	MinVcpuCount    *int    `pulumi:"minVcpuCount,optional" provider:"replaceOnChanges"`
	MinMemoryInGb   *int    `pulumi:"minMemoryInGb,optional" provider:"replaceOnChanges"`
	CudaVersion     *string `pulumi:"cudaVersion,optional" provider:"replaceOnChanges"`
	ComputeType     *string `pulumi:"computeType,optional" provider:"replaceOnChanges"`
	GlobalNetwork   *bool   `pulumi:"globalNetwork,optional" provider:"replaceOnChanges"`
	CountryCode     *string `pulumi:"countryCode,optional" provider:"replaceOnChanges"`
	StopAfter       *string `pulumi:"stopAfter,optional" provider:"replaceOnChanges"`
	TerminateAfter  *string `pulumi:"terminateAfter,optional" provider:"replaceOnChanges"`

	GpuTypeIDList       []string `pulumi:"gpuTypeIdList,optional" provider:"replaceOnChanges"`
	AllowedCudaVersions []string `pulumi:"allowedCudaVersions,optional" provider:"replaceOnChanges"`

	MinCudaVersion *string  `pulumi:"minCudaVersion,optional" provider:"replaceOnChanges"`
	DeployCost     *float64 `pulumi:"deployCost,optional" provider:"replaceOnChanges"`
	MinDisk        *int     `pulumi:"minDisk,optional" provider:"replaceOnChanges"`
	MinDownload    *int     `pulumi:"minDownload,optional" provider:"replaceOnChanges"`
	MinUpload      *int     `pulumi:"minUpload,optional" provider:"replaceOnChanges"`
	VolumeKey      *string  `pulumi:"volumeKey,optional" provider:"replaceOnChanges"`

	AiAPIID    *string `pulumi:"aiApiId,optional" provider:"replaceOnChanges"`
	IdeAiAPIID *string `pulumi:"ideAiApiId,optional" provider:"replaceOnChanges"`

	InstanceIDs     []string `pulumi:"instanceIds,optional" provider:"replaceOnChanges"`
	ModelReferences []string `pulumi:"modelReferences,optional" provider:"replaceOnChanges"`

	SavingsPlan *SavingsPlan `pulumi:"savingsPlan,optional" provider:"replaceOnChanges"`

	// Mutable fields (in PodEditJobInput — can be updated in-place)
	ImageName               *string           `pulumi:"imageName,optional"`
	DockerArgs              *string           `pulumi:"dockerArgs,optional"`
	Env                     map[string]string `pulumi:"env,optional"`
	Ports                   *string           `pulumi:"ports,optional"`
	VolumeInGb              *int              `pulumi:"volumeInGb,optional"`
	VolumeMountPath         *string           `pulumi:"volumeMountPath,optional"`
	ContainerDiskInGb       *int              `pulumi:"containerDiskInGb,optional"`
	ContainerRegistryAuthID *string           `pulumi:"containerRegistryAuthId,optional"`
}

// Annotate provides descriptions for PodArgs fields.
func (a *PodArgs) Annotate(an infer.Annotator) {
	an.Describe(&a.Name, "A name for the pod.")
	an.Describe(&a.GpuTypeID,
		"The GPU type ID to deploy (e.g. \"NVIDIA GeForce RTX 4090\").")
	an.Describe(&a.GpuCount, "The number of GPUs to allocate.")
	an.Describe(&a.CloudType,
		"The cloud type: SECURE, COMMUNITY, or ALL.")
	an.Describe(&a.ImageName,
		"The Docker image to run on the pod.")
	an.Describe(&a.DockerArgs,
		"Docker arguments to pass to the container.")
	an.Describe(&a.Env, "Environment variables as key-value pairs.")
	an.Describe(&a.Ports,
		"Ports to expose (e.g. \"8080/http,22/tcp\").")
	an.Describe(&a.VolumeInGb,
		"The size of the persistent volume in GB.")
	an.Describe(&a.VolumeMountPath,
		"The path to mount the persistent volume.")
	an.Describe(&a.ContainerDiskInGb,
		"The size of the container disk in GB.")
	an.Describe(&a.TemplateID,
		"The template ID to use for the pod.")
	an.Describe(&a.NetworkVolumeID,
		"The network volume ID to attach to the pod.")
	an.Describe(&a.ContainerRegistryAuthID,
		"The container registry auth ID for pulling private images.")
	an.Describe(&a.DataCenterID,
		"The data center ID to deploy the pod in.")
	an.Describe(&a.StartJupyter,
		"Whether to start a Jupyter notebook server.")
	an.Describe(&a.StartSSH,
		"Whether to start an SSH server.")
	an.Describe(&a.SupportPublicIP,
		"Whether to assign a public IP address.")
	an.Describe(&a.MinVcpuCount,
		"Minimum number of vCPUs required.")
	an.Describe(&a.MinMemoryInGb,
		"Minimum memory in GB required.")
	an.Describe(&a.CudaVersion, "The CUDA version to use.")
	an.Describe(&a.ComputeType,
		"The compute type: CPU or GPU.")
	an.Describe(&a.GlobalNetwork,
		"Whether to enable global networking.")
	an.Describe(&a.CountryCode,
		"The country code for data residency.")
	an.Describe(&a.StopAfter,
		"Duration after which the pod is automatically stopped.")
	an.Describe(&a.TerminateAfter,
		"Duration after which the pod is automatically terminated.")
	an.Describe(&a.GpuTypeIDList,
		"A list of acceptable GPU type IDs (fallback options).")
	an.Describe(&a.AllowedCudaVersions,
		"A list of allowed CUDA versions.")
	an.Describe(&a.MinCudaVersion,
		"The minimum CUDA version required.")
	an.Describe(&a.DeployCost,
		"The maximum bid price per GPU per hour for spot instances.")
	an.Describe(&a.MinDisk,
		"Minimum disk space in GB required on the host.")
	an.Describe(&a.MinDownload,
		"Minimum download bandwidth in Mbps.")
	an.Describe(&a.MinUpload,
		"Minimum upload bandwidth in Mbps.")
	an.Describe(&a.VolumeKey,
		"The volume key for persistent storage.")
	an.Describe(&a.AiAPIID, "The AI API ID for the pod.")
	an.Describe(&a.IdeAiAPIID,
		"The IDE AI API ID for the pod.")
	an.Describe(&a.InstanceIDs,
		"Specific instance IDs to deploy on.")
	an.Describe(&a.ModelReferences,
		"Model references for the pod.")
	an.Describe(&a.SavingsPlan,
		"Savings plan configuration for reduced pricing.")
}

// Annotate provides descriptions for SavingsPlan fields.
func (s *SavingsPlan) Annotate(a infer.Annotator) {
	a.Describe(&s.PlanLength,
		"The length of the savings plan.")
	a.Describe(&s.UpfrontCost,
		"The upfront cost for the savings plan.")
}

// PodState is the persisted state of a pod resource.
type PodState struct {
	PodArgs
	// Outputs
	PodID          string  `pulumi:"podId"`
	MachineID      string  `pulumi:"machineId"`
	CostPerHr      float64 `pulumi:"costPerHr"`
	DesiredStatus  string  `pulumi:"desiredStatus"`
	VcpuCount      float64 `pulumi:"vcpuCount"`
	MemoryInGb     float64 `pulumi:"memoryInGb"`
	OutputGpuCount int     `pulumi:"outputGpuCount"`

	OutputContainerDiskInGb *int     `pulumi:"outputContainerDiskInGb,optional"`
	OutputVolumeInGb        *float64 `pulumi:"outputVolumeInGb,optional"`
	OutputPorts             *string  `pulumi:"outputPorts,optional"`
	OutputTemplateID        *string  `pulumi:"outputTemplateId,optional"`
	OutputNetworkVolumeID   *string  `pulumi:"outputNetworkVolumeId,optional"`
	OutputPodType           *string  `pulumi:"outputPodType,optional"`

	OutputContainerRegistryAuthID *string `pulumi:"outputContainerRegistryAuthId,optional"`
}

// Annotate provides descriptions for PodState output fields.
func (s *PodState) Annotate(a infer.Annotator) {
	a.Describe(&s.PodID,
		"The unique identifier of the pod.")
	a.Describe(&s.MachineID,
		"The ID of the machine the pod is running on.")
	a.Describe(&s.CostPerHr,
		"The cost per hour for the pod in USD.")
	a.Describe(&s.DesiredStatus,
		"The desired status of the pod.")
	a.Describe(&s.VcpuCount,
		"The number of vCPUs allocated.")
	a.Describe(&s.MemoryInGb,
		"The amount of memory allocated in GB.")
	a.Describe(&s.OutputGpuCount,
		"The number of GPUs allocated (from API response).")
	a.Describe(&s.OutputContainerDiskInGb,
		"The container disk size in GB (from API response).")
	a.Describe(&s.OutputVolumeInGb,
		"The volume size in GB (from API response).")
	a.Describe(&s.OutputPorts,
		"The exposed ports (from API response).")
	a.Describe(&s.OutputTemplateID,
		"The template ID used (from API response).")
	a.Describe(&s.OutputNetworkVolumeID,
		"The network volume ID attached (from API response).")
	a.Describe(&s.OutputContainerRegistryAuthID,
		"The container registry auth ID (from API response).")
	a.Describe(&s.OutputPodType,
		"The pod type (from API response).")
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
		StartSsh:                input.StartSSH,
		SupportPublicIp:         input.SupportPublicIP,
		MinVcpuCount:            input.MinVcpuCount,
		MinMemoryInGb:           input.MinMemoryInGb,
		CudaVersion:             input.CudaVersion,
		GlobalNetwork:           input.GlobalNetwork,
		CountryCode:             input.CountryCode,
		StopAfter:               input.StopAfter,
		TerminateAfter:          input.TerminateAfter,
		GpuTypeIdList:           runpod.StringPtrSlice(input.GpuTypeIDList),
		AllowedCudaVersions:     runpod.StringPtrSlice(input.AllowedCudaVersions),
		MinCudaVersion:          input.MinCudaVersion,
		DeployCost:              input.DeployCost,
		MinDisk:                 input.MinDisk,
		MinDownload:             input.MinDownload,
		MinUpload:               input.MinUpload,
		VolumeKey:               input.VolumeKey,
		AiApiId:                 input.AiAPIID,
		IdeAiApiId:              input.IdeAiAPIID,
		InstanceIds:             runpod.StringPtrSlice(input.InstanceIDs),
		ModelReferences:         runpod.StringPtrSlice(input.ModelReferences),
	}

	// CloudType needs conversion to the enum
	if input.CloudType != nil {
		ct := runpod.CloudTypeEnum(*input.CloudType)
		createInput.CloudType = &ct
	}

	// ComputeType needs conversion to the enum
	if input.ComputeType != nil {
		ct := runpod.ComputeType(*input.ComputeType)
		createInput.ComputeType = &ct
	}

	// SavingsPlan nested type
	if input.SavingsPlan != nil {
		createInput.SavingsPlan = &runpod.SavingsPlanInput{
			PlanLength:  input.SavingsPlan.PlanLength,
			UpfrontCost: input.SavingsPlan.UpfrontCost,
		}
	}

	resp, err := runpod.CreatePod(ctx, client, createInput)
	if err != nil {
		return infer.CreateResponse[PodState]{}, err
	}

	if resp.PodFindAndDeployOnDemand == nil {
		return infer.CreateResponse[PodState]{},
			errors.New("API returned nil pod")
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
		// Resource was deleted externally — return empty ID so Pulumi removes it from state.
		return infer.ReadResponse[PodArgs, PodState]{ID: ""}, nil
	}

	state := podResponseToState(req.Inputs, resp.Pod)
	return infer.ReadResponse[PodArgs, PodState]{
		ID:     resp.Pod.Id,
		Inputs: req.Inputs,
		State:  state,
	}, nil
}

// Update modifies mutable pod fields.
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

	containerDiskInGb := 0
	if req.Inputs.ContainerDiskInGb != nil {
		containerDiskInGb = *req.Inputs.ContainerDiskInGb
	}

	updateInput := runpod.PodEditJobInput{
		PodId:                   req.ID,
		ImageName:               imageName,
		DockerArgs:              req.Inputs.DockerArgs,
		Env:                     runpod.EnvMapToGQL(req.Inputs.Env),
		Ports:                   req.Inputs.Ports,
		ContainerDiskInGb:       containerDiskInGb,
		VolumeInGb:              req.Inputs.VolumeInGb,
		VolumeMountPath:         req.Inputs.VolumeMountPath,
		ContainerRegistryAuthId: req.Inputs.ContainerRegistryAuthID,
	}

	resp, err := runpod.UpdatePod(ctx, client, updateInput)
	if err != nil {
		return infer.UpdateResponse[PodState]{}, err
	}

	if resp.PodEditJob == nil {
		return infer.UpdateResponse[PodState]{},
			errors.New("API returned nil pod on update")
	}

	state := podResponseToState(req.Inputs, resp.PodEditJob)
	return infer.UpdateResponse[PodState]{Output: state}, nil
}

// Delete terminates the pod.
func (Pod) Delete(
	ctx context.Context,
	req infer.DeleteRequest[PodState],
) (infer.DeleteResponse, error) {
	client := getClient(ctx)
	_, err := runpod.TerminatePod(
		ctx, client, runpod.PodTerminateInput{PodId: req.ID},
	)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func podResponseToState(
	input PodArgs, pod *runpod.PodResponse,
) PodState {
	state := PodState{
		PodArgs:        input,
		PodID:          pod.Id,
		MachineID:      pod.MachineId,
		CostPerHr:      pod.CostPerHr,
		DesiredStatus:  string(pod.DesiredStatus),
		VcpuCount:      pod.VcpuCount,
		MemoryInGb:     pod.MemoryInGb,
		OutputGpuCount: pod.GpuCount,

		OutputContainerDiskInGb:       pod.ContainerDiskInGb,
		OutputVolumeInGb:              pod.VolumeInGb,
		OutputPorts:                   pod.Ports,
		OutputTemplateID:              pod.TemplateId,
		OutputNetworkVolumeID:         pod.NetworkVolumeId,
		OutputContainerRegistryAuthID: pod.ContainerRegistryAuthId,
	}
	if pod.PodType != nil {
		pt := string(*pod.PodType)
		state.OutputPodType = &pt
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
