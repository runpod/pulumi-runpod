package runpod

// EnvironmentVariable is the input format for env vars in mutations.
type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Pod represents a RunPod GPU pod.
type Pod struct {
	ID                      string   `json:"id"`
	Name                    string   `json:"name"`
	MachineID               string   `json:"machineId"`
	ImageName               string   `json:"imageName"`
	DockerArgs              string   `json:"dockerArgs"`
	GpuCount                int      `json:"gpuCount"`
	VcpuCount               float64  `json:"vcpuCount"`
	MemoryInGb              float64  `json:"memoryInGb"`
	ContainerDiskInGb       int      `json:"containerDiskInGb"`
	VolumeInGb              float64  `json:"volumeInGb"`
	VolumeMountPath         string   `json:"volumeMountPath"`
	DesiredStatus           string   `json:"desiredStatus"`
	CostPerHr               float64  `json:"costPerHr"`
	Env                     []string `json:"env"`
	Ports                   string   `json:"ports"`
	TemplateID              string   `json:"templateId"`
	NetworkVolumeID         string   `json:"networkVolumeId"`
	ContainerRegistryAuthID string   `json:"containerRegistryAuthId"`
	PodType                 string   `json:"podType"`
}

// PodTemplate represents a RunPod pod template.
type PodTemplate struct {
	ID                      string                `json:"id"`
	Name                    string                `json:"name"`
	ImageName               string                `json:"imageName"`
	DockerArgs              string                `json:"dockerArgs"`
	Env                     []EnvironmentVariable `json:"env"`
	Ports                   string                `json:"ports"`
	VolumeMountPath         string                `json:"volumeMountPath"`
	VolumeInGb              int                   `json:"volumeInGb"`
	ContainerDiskInGb       int                   `json:"containerDiskInGb"`
	StartJupyter            bool                  `json:"startJupyter"`
	StartSsh                bool                  `json:"startSsh"`
	StartScript             string                `json:"startScript"`
	IsServerless            bool                  `json:"isServerless"`
	IsPublic                bool                  `json:"isPublic"`
	ContainerRegistryAuthID string                `json:"containerRegistryAuthId"`
}

// Endpoint represents a RunPod serverless endpoint.
type Endpoint struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	TemplateID       string   `json:"templateId"`
	GpuIds           string   `json:"gpuIds"`
	WorkersMin       int      `json:"workersMin"`
	WorkersMax       int      `json:"workersMax"`
	IdleTimeout      int      `json:"idleTimeout"`
	Locations        string   `json:"locations"`
	ScalerType       string   `json:"scalerType"`
	ScalerValue      int      `json:"scalerValue"`
	NetworkVolumeID  string   `json:"networkVolumeId"`
	GpuCount         int      `json:"gpuCount"`
	InstanceIds      []string `json:"instanceIds"`
	WorkersPFBTarget int      `json:"workersPFBTarget"`
}

// NetworkVolume represents a RunPod network volume.
type NetworkVolume struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	DataCenterID string `json:"dataCenterId"`
}

// GpuType represents a GPU type available on RunPod.
type GpuType struct {
	ID             string  `json:"id"`
	DisplayName    string  `json:"displayName"`
	MemoryInGb     int     `json:"memoryInGb"`
	SecureCloud    bool    `json:"secureCloud"`
	CommunityCloud bool    `json:"communityCloud"`
	SecurePrice    float64 `json:"securePrice"`
	CommunityPrice float64 `json:"communityPrice"`
	MaxGpuCount    int     `json:"maxGpuCount"`
}
