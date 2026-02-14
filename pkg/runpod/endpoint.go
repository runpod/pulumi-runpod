package runpod

import (
	"context"
	"fmt"
)

// SaveEndpointInput holds the parameters for creating/updating an endpoint.
type SaveEndpointInput struct {
	ID              string                `json:"id,omitempty"`
	Name            string                `json:"name"`
	TemplateID      string                `json:"templateId,omitempty"`
	GpuIds          string                `json:"gpuIds,omitempty"`
	WorkersMin      int                   `json:"workersMin,omitempty"`
	WorkersMax      int                   `json:"workersMax,omitempty"`
	IdleTimeout     int                   `json:"idleTimeout,omitempty"`
	Locations       string                `json:"locations,omitempty"`
	ScalerType      string                `json:"scalerType,omitempty"`
	ScalerValue     int                   `json:"scalerValue,omitempty"`
	NetworkVolumeID string                `json:"networkVolumeId,omitempty"`
	GpuCount        int                   `json:"gpuCount,omitempty"`
	InstanceIds     []string              `json:"instanceIds,omitempty"`
	Env             []EnvironmentVariable `json:"env,omitempty"`
}

const saveEndpointQuery = `
mutation saveEndpoint($input: EndpointInput!) {
  saveEndpoint(input: $input) {
    id
    name
    templateId
    gpuIds
    workersMin
    workersMax
    idleTimeout
    locations
    scalerType
    scalerValue
    networkVolumeId
    gpuCount
    instanceIds
    workersPFBTarget
  }
}
`

const getEndpointsQuery = `
query myself {
  myself {
    endpoints {
      id
      name
      templateId
      gpuIds
      workersMin
      workersMax
      idleTimeout
      locations
      scalerType
      scalerValue
      networkVolumeId
      gpuCount
      instanceIds
      workersPFBTarget
    }
  }
}
`

const deleteEndpointQuery = `
mutation deleteEndpoint($id: String!) {
  deleteEndpoint(id: $id)
}
`

// CreateEndpoint creates a new serverless endpoint.
func (c *Client) CreateEndpoint(ctx context.Context, input SaveEndpointInput) (*Endpoint, error) {
	input.ID = "" // Ensure no ID for creation
	return c.saveEndpoint(ctx, input)
}

// UpdateEndpoint updates an existing serverless endpoint.
func (c *Client) UpdateEndpoint(ctx context.Context, input SaveEndpointInput) (*Endpoint, error) {
	if input.ID == "" {
		return nil, fmt.Errorf("endpoint ID is required for update")
	}
	return c.saveEndpoint(ctx, input)
}

func (c *Client) saveEndpoint(ctx context.Context, input SaveEndpointInput) (*Endpoint, error) {
	var result struct {
		SaveEndpoint *Endpoint `json:"saveEndpoint"`
	}

	vars := map[string]any{"input": input}
	if err := c.Do(ctx, saveEndpointQuery, vars, &result); err != nil {
		return nil, FormatError("saving", "endpoint", input.ID, err)
	}

	if result.SaveEndpoint == nil {
		return nil, FormatError("saving", "endpoint", input.ID, fmt.Errorf("API returned nil endpoint"))
	}

	return result.SaveEndpoint, nil
}

// GetEndpoint retrieves an endpoint by ID.
func (c *Client) GetEndpoint(ctx context.Context, endpointID string) (*Endpoint, error) {
	var result struct {
		Myself struct {
			Endpoints []Endpoint `json:"endpoints"`
		} `json:"myself"`
	}

	if err := c.Do(ctx, getEndpointsQuery, nil, &result); err != nil {
		return nil, FormatError("reading", "endpoint", endpointID, err)
	}

	for _, e := range result.Myself.Endpoints {
		if e.ID == endpointID {
			return &e, nil
		}
	}

	return nil, nil // Not found
}

// DeleteEndpoint deletes an endpoint by ID.
func (c *Client) DeleteEndpoint(ctx context.Context, endpointID string) error {
	vars := map[string]any{"id": endpointID}
	if err := c.Do(ctx, deleteEndpointQuery, vars, nil); err != nil {
		return FormatError("deleting", "endpoint", endpointID, err)
	}
	return nil
}
