package runpod

import (
	"context"
	"fmt"
)

// CreateNetworkVolumeInput holds the parameters for creating a network volume.
type CreateNetworkVolumeInput struct {
	Name         string `json:"name"`
	Size         int    `json:"size"`
	DataCenterID string `json:"dataCenterId"`
}

// UpdateNetworkVolumeInput holds the parameters for updating a network volume.
type UpdateNetworkVolumeInput struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	Size int    `json:"size,omitempty"`
}

const createNetworkVolumeQuery = `
mutation createNetworkVolume($input: NetworkVolumeInput!) {
  createNetworkVolume(input: $input) {
    id
    name
    size
    dataCenterId
  }
}
`

const getNetworkVolumesQuery = `
query myself {
  myself {
    networkVolumes {
      id
      name
      size
      dataCenterId
    }
  }
}
`

const updateNetworkVolumeQuery = `
mutation updateNetworkVolume($input: UpdateNetworkVolumeInput!) {
  updateNetworkVolume(input: $input) {
    id
    name
    size
    dataCenterId
  }
}
`

const deleteNetworkVolumeQuery = `
mutation deleteNetworkVolume($id: String!) {
  deleteNetworkVolume(id: $id)
}
`

// CreateNetworkVolume creates a new network volume.
func (c *Client) CreateNetworkVolume(ctx context.Context, input CreateNetworkVolumeInput) (*NetworkVolume, error) {
	var result struct {
		CreateNetworkVolume *NetworkVolume `json:"createNetworkVolume"`
	}

	vars := map[string]any{"input": input}
	if err := c.Do(ctx, createNetworkVolumeQuery, vars, &result); err != nil {
		return nil, FormatError("creating", "network volume", "", err)
	}

	if result.CreateNetworkVolume == nil {
		return nil, FormatError("creating", "network volume", "", fmt.Errorf("API returned nil network volume"))
	}

	return result.CreateNetworkVolume, nil
}

// GetNetworkVolume retrieves a network volume by ID.
func (c *Client) GetNetworkVolume(ctx context.Context, volumeID string) (*NetworkVolume, error) {
	var result struct {
		Myself struct {
			NetworkVolumes []NetworkVolume `json:"networkVolumes"`
		} `json:"myself"`
	}

	if err := c.Do(ctx, getNetworkVolumesQuery, nil, &result); err != nil {
		return nil, FormatError("reading", "network volume", volumeID, err)
	}

	for _, v := range result.Myself.NetworkVolumes {
		if v.ID == volumeID {
			return &v, nil
		}
	}

	return nil, nil // Not found
}

// UpdateNetworkVolume updates a network volume.
func (c *Client) UpdateNetworkVolume(ctx context.Context, input UpdateNetworkVolumeInput) (*NetworkVolume, error) {
	var result struct {
		UpdateNetworkVolume *NetworkVolume `json:"updateNetworkVolume"`
	}

	vars := map[string]any{"input": input}
	if err := c.Do(ctx, updateNetworkVolumeQuery, vars, &result); err != nil {
		return nil, FormatError("updating", "network volume", input.ID, err)
	}

	return result.UpdateNetworkVolume, nil
}

// DeleteNetworkVolume deletes a network volume by ID.
func (c *Client) DeleteNetworkVolume(ctx context.Context, volumeID string) error {
	vars := map[string]any{"id": volumeID}
	if err := c.Do(ctx, deleteNetworkVolumeQuery, vars, nil); err != nil {
		return FormatError("deleting", "network volume", volumeID, err)
	}
	return nil
}
