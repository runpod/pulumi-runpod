package runpod

import (
	"context"
	"fmt"
)

// CreatePodInput holds the parameters for creating a pod.
type CreatePodInput struct {
	Name                    string                `json:"name"`
	GpuTypeID               string                `json:"gpuTypeId"`
	GpuCount                int                   `json:"gpuCount,omitempty"`
	CloudType               string                `json:"cloudType,omitempty"`
	ImageName               string                `json:"imageName,omitempty"`
	DockerArgs              string                `json:"dockerArgs,omitempty"`
	Env                     []EnvironmentVariable `json:"env,omitempty"`
	Ports                   string                `json:"ports,omitempty"`
	VolumeInGb              int                   `json:"volumeInGb,omitempty"`
	VolumeMountPath         string                `json:"volumeMountPath,omitempty"`
	ContainerDiskInGb       int                   `json:"containerDiskInGb,omitempty"`
	TemplateID              string                `json:"templateId,omitempty"`
	NetworkVolumeID         string                `json:"networkVolumeId,omitempty"`
	ContainerRegistryAuthID string                `json:"containerRegistryAuthId,omitempty"`
	DataCenterID            string                `json:"dataCenterId,omitempty"`
	StartJupyter            *bool                 `json:"startJupyter,omitempty"`
	StartSsh                *bool                 `json:"startSsh,omitempty"`
	SupportPublicIP         *bool                 `json:"supportPublicIp,omitempty"`
	MinVcpuCount            int                   `json:"minVcpuCount,omitempty"`
	MinMemoryInGb           int                   `json:"minMemoryInGb,omitempty"`
	MinDisk                 int                   `json:"minDisk,omitempty"`
	StopAfter               string                `json:"stopAfter,omitempty"`
	TerminateAfter          string                `json:"terminateAfter,omitempty"`
	CudaVersion             string                `json:"cudaVersion,omitempty"`
}

// UpdatePodInput holds the parameters for updating a pod.
type UpdatePodInput struct {
	PodID      string                `json:"podId"`
	ImageName  string                `json:"imageName,omitempty"`
	DockerArgs string                `json:"dockerArgs,omitempty"`
	Env        []EnvironmentVariable `json:"env,omitempty"`
}

const createPodQuery = `
mutation podFindAndDeployOnDemand($input: PodFindAndDeployOnDemandInput!) {
  podFindAndDeployOnDemand(input: $input) {
    id
    name
    machineId
    imageName
    dockerArgs
    gpuCount
    vcpuCount
    memoryInGb
    containerDiskInGb
    volumeInGb
    volumeMountPath
    desiredStatus
    costPerHr
    env
    ports
    templateId
    networkVolumeId
    containerRegistryAuthId
    podType
  }
}
`

const getPodQuery = `
query pod($input: PodFilter!) {
  pod(input: $input) {
    id
    name
    machineId
    imageName
    dockerArgs
    gpuCount
    vcpuCount
    memoryInGb
    containerDiskInGb
    volumeInGb
    volumeMountPath
    desiredStatus
    costPerHr
    env
    ports
    templateId
    networkVolumeId
    containerRegistryAuthId
    podType
  }
}
`

const updatePodQuery = `
mutation podEditJob($input: PodEditJobInput!) {
  podEditJob(input: $input) {
    id
    name
    machineId
    imageName
    dockerArgs
    gpuCount
    vcpuCount
    memoryInGb
    containerDiskInGb
    volumeInGb
    volumeMountPath
    desiredStatus
    costPerHr
    env
    ports
    templateId
    networkVolumeId
    containerRegistryAuthId
    podType
  }
}
`

const terminatePodQuery = `
mutation podTerminate($input: PodFilter!) {
  podTerminate(input: $input) {
    id
  }
}
`

// CreatePod creates a new GPU pod using podFindAndDeployOnDemand.
func (c *Client) CreatePod(ctx context.Context, input CreatePodInput) (*Pod, error) {
	var result struct {
		PodFindAndDeployOnDemand *Pod `json:"podFindAndDeployOnDemand"`
	}

	vars := map[string]any{"input": input}
	if err := c.Do(ctx, createPodQuery, vars, &result); err != nil {
		return nil, FormatError("creating", "pod", "", err)
	}

	if result.PodFindAndDeployOnDemand == nil {
		return nil, FormatError("creating", "pod", "", fmt.Errorf("API returned nil pod"))
	}

	return result.PodFindAndDeployOnDemand, nil
}

// GetPod retrieves a pod by ID.
func (c *Client) GetPod(ctx context.Context, podID string) (*Pod, error) {
	var result struct {
		Pod *Pod `json:"pod"`
	}

	vars := map[string]any{
		"input": map[string]any{"podId": podID},
	}
	if err := c.Do(ctx, getPodQuery, vars, &result); err != nil {
		return nil, FormatError("reading", "pod", podID, err)
	}

	return result.Pod, nil
}

// UpdatePod updates a pod's mutable fields.
func (c *Client) UpdatePod(ctx context.Context, input UpdatePodInput) (*Pod, error) {
	var result struct {
		PodEditJob *Pod `json:"podEditJob"`
	}

	vars := map[string]any{"input": input}
	if err := c.Do(ctx, updatePodQuery, vars, &result); err != nil {
		return nil, FormatError("updating", "pod", input.PodID, err)
	}

	return result.PodEditJob, nil
}

// TerminatePod deletes a pod.
func (c *Client) TerminatePod(ctx context.Context, podID string) error {
	vars := map[string]any{
		"input": map[string]any{"podId": podID},
	}
	if err := c.Do(ctx, terminatePodQuery, vars, nil); err != nil {
		return FormatError("terminating", "pod", podID, err)
	}
	return nil
}
