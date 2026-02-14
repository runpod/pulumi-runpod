package runpod

import "context"

const getGpuTypesQuery = `
query gpuTypes {
  gpuTypes {
    id
    displayName
    memoryInGb
    secureCloud
    communityCloud
    securePrice
    communityPrice
    maxGpuCount
  }
}
`

// GetGpuTypes retrieves the list of available GPU types.
func (c *Client) GetGpuTypes(ctx context.Context) ([]GpuType, error) {
	var result struct {
		GpuTypes []GpuType `json:"gpuTypes"`
	}

	if err := c.Do(ctx, getGpuTypesQuery, nil, &result); err != nil {
		return nil, FormatError("reading", "GPU types", "", err)
	}

	return result.GpuTypes, nil
}
