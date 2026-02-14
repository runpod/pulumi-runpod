package runpod

import (
	"context"
	"fmt"
)

// SaveTemplateInput holds the parameters for creating/updating a template.
type SaveTemplateInput struct {
	ID                      string                `json:"id,omitempty"`
	Name                    string                `json:"name"`
	ImageName               string                `json:"imageName,omitempty"`
	ContainerDiskInGb       int                   `json:"containerDiskInGb"`
	DockerArgs              string                `json:"dockerArgs"`
	Env                     []EnvironmentVariable `json:"env"`
	Ports                   string                `json:"ports,omitempty"`
	VolumeMountPath         string                `json:"volumeMountPath,omitempty"`
	VolumeInGb              int                   `json:"volumeInGb"`
	StartJupyter            *bool                 `json:"startJupyter,omitempty"`
	StartSsh                *bool                 `json:"startSsh,omitempty"`
	StartScript             string                `json:"startScript,omitempty"`
	IsServerless            *bool                 `json:"isServerless,omitempty"`
	IsPublic                *bool                 `json:"isPublic,omitempty"`
	ContainerRegistryAuthID string                `json:"containerRegistryAuthId,omitempty"`
}

const saveTemplateQuery = `
mutation saveTemplate($input: SaveTemplateInput!) {
  saveTemplate(input: $input) {
    id
    name
    imageName
    dockerArgs
    env {
      key
      value
    }
    ports
    volumeMountPath
    volumeInGb
    containerDiskInGb
    startJupyter
    startSsh
    startScript
    isServerless
    isPublic
    containerRegistryAuthId
  }
}
`

const getTemplateQuery = `
query myself {
  myself {
    podTemplates {
      id
      name
      imageName
      dockerArgs
      env {
        key
        value
      }
      ports
      volumeMountPath
      volumeInGb
      containerDiskInGb
      startJupyter
      startSsh
      startScript
      isServerless
      isPublic
      containerRegistryAuthId
    }
  }
}
`

const deleteTemplateQuery = `
mutation deleteTemplate($templateName: String!) {
  deleteTemplate(templateName: $templateName)
}
`

// CreateTemplate creates a new pod template.
func (c *Client) CreateTemplate(ctx context.Context, input SaveTemplateInput) (*PodTemplate, error) {
	input.ID = "" // Ensure no ID for creation
	return c.saveTemplate(ctx, input)
}

// UpdateTemplate updates an existing pod template.
func (c *Client) UpdateTemplate(ctx context.Context, input SaveTemplateInput) (*PodTemplate, error) {
	if input.ID == "" {
		return nil, fmt.Errorf("template ID is required for update")
	}
	return c.saveTemplate(ctx, input)
}

func (c *Client) saveTemplate(ctx context.Context, input SaveTemplateInput) (*PodTemplate, error) {
	var result struct {
		SaveTemplate *PodTemplate `json:"saveTemplate"`
	}

	vars := map[string]any{"input": input}
	if err := c.Do(ctx, saveTemplateQuery, vars, &result); err != nil {
		return nil, FormatError("saving", "template", input.ID, err)
	}

	if result.SaveTemplate == nil {
		return nil, FormatError("saving", "template", input.ID, fmt.Errorf("API returned nil template"))
	}

	return result.SaveTemplate, nil
}

// GetTemplate retrieves a template by ID from the user's templates.
func (c *Client) GetTemplate(ctx context.Context, templateID string) (*PodTemplate, error) {
	var result struct {
		Myself struct {
			PodTemplates []PodTemplate `json:"podTemplates"`
		} `json:"myself"`
	}

	if err := c.Do(ctx, getTemplateQuery, nil, &result); err != nil {
		return nil, FormatError("reading", "template", templateID, err)
	}

	for _, t := range result.Myself.PodTemplates {
		if t.ID == templateID {
			return &t, nil
		}
	}

	return nil, nil // Not found
}

// DeleteTemplate deletes a template by name.
// Note: The RunPod API deleteTemplate mutation takes the template name, not ID.
func (c *Client) DeleteTemplate(ctx context.Context, templateName string) error {
	vars := map[string]any{"templateName": templateName}
	if err := c.Do(ctx, deleteTemplateQuery, vars, nil); err != nil {
		return FormatError("deleting", "template", templateName, err)
	}
	return nil
}
