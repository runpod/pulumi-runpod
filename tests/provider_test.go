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

package tests

import (
	"context"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/require"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	xyz "github.com/runpod/pulumi-runpod/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
)

func TestTemplateDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("Template"),
		Properties: property.NewMap(map[string]property.Value{
			"name":              property.New("test-template"),
			"imageName":         property.New("runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04"),
			"containerDiskInGb": property.New(20.0),
			"volumeInGb":        property.New(20.0),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "test-template", response.Properties.Get("name").AsString())
	require.Equal(t, "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04", response.Properties.Get("imageName").AsString())
}

func TestTemplateWithNewFieldsDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("Template"),
		Properties: property.NewMap(map[string]property.Value{
			"name":              property.New("test-template-cpu"),
			"imageName":         property.New("runpod/pytorch:latest"),
			"containerDiskInGb": property.New(10.0),
			"volumeInGb":        property.New(10.0),
			"readme":            property.New("# My Template"),
			"advancedStart":     property.New(true),
			"category":          property.New("CPU"),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "test-template-cpu", response.Properties.Get("name").AsString())
	require.Equal(t, "# My Template", response.Properties.Get("readme").AsString())
	require.Equal(t, true, response.Properties.Get("advancedStart").AsBool())
	require.Equal(t, "CPU", response.Properties.Get("category").AsString())
}

func TestPodDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("Pod"),
		Properties: property.NewMap(map[string]property.Value{
			"name":      property.New("test-pod"),
			"gpuTypeId": property.New("NVIDIA RTX A4000"),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "test-pod", response.Properties.Get("name").AsString())
	require.Equal(t, "NVIDIA RTX A4000", response.Properties.Get("gpuTypeId").AsString())
}

func TestPodWithComputeTypeDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("Pod"),
		Properties: property.NewMap(map[string]property.Value{
			"name":          property.New("test-cpu-pod"),
			"gpuTypeId":     property.New("NVIDIA RTX A4000"),
			"computeType":   property.New("CPU"),
			"globalNetwork": property.New(true),
			"countryCode":   property.New("US"),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "test-cpu-pod", response.Properties.Get("name").AsString())
	require.Equal(t, "CPU", response.Properties.Get("computeType").AsString())
	require.Equal(t, true, response.Properties.Get("globalNetwork").AsBool())
	require.Equal(t, "US", response.Properties.Get("countryCode").AsString())
}

func TestNetworkVolumeDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("NetworkVolume"),
		Properties: property.NewMap(map[string]property.Value{
			"name":         property.New("test-volume"),
			"size":         property.New(20.0),
			"dataCenterId": property.New("US-TX-3"),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "test-volume", response.Properties.Get("name").AsString())
}

func TestNetworkVolumeWithNextGenDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("NetworkVolume"),
		Properties: property.NewMap(map[string]property.Value{
			"name":             property.New("test-nextgen-volume"),
			"size":             property.New(50.0),
			"dataCenterId":     property.New("US-TX-3"),
			"isNextGenStorage": property.New(true),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "test-nextgen-volume", response.Properties.Get("name").AsString())
	require.Equal(t, true, response.Properties.Get("isNextGenStorage").AsBool())
}

func TestEndpointDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("Endpoint"),
		Properties: property.NewMap(map[string]property.Value{
			"name":       property.New("test-endpoint"),
			"templateId": property.New("abc123"),
			"gpuIds":     property.New("AMPERE_16"),
			"workersMin": property.New(0.0),
			"workersMax": property.New(3.0),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "test-endpoint", response.Properties.Get("name").AsString())
}

func TestEndpointWithNewFieldsDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("Endpoint"),
		Properties: property.NewMap(map[string]property.Value{
			"name":               property.New("test-flashboot-endpoint"),
			"templateId":         property.New("abc123"),
			"gpuIds":             property.New("AMPERE_16"),
			"workersMin":         property.New(0.0),
			"workersMax":         property.New(3.0),
			"flashBootType":      property.New("FLASHBOOT"),
			"executionTimeoutMs": property.New(30000.0),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "test-flashboot-endpoint", response.Properties.Get("name").AsString())
	require.Equal(t, "FLASHBOOT", response.Properties.Get("flashBootType").AsString())
	require.Equal(t, 30000.0, response.Properties.Get("executionTimeoutMs").AsNumber())
}

func TestSecretDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("Secret"),
		Properties: property.NewMap(map[string]property.Value{
			"name":        property.New("MY_SECRET"),
			"value":       property.New("super-secret-value"),
			"description": property.New("A test secret"),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "MY_SECRET", response.Properties.Get("name").AsString())
	require.Equal(t, "A test secret", response.Properties.Get("description").AsString())
}

func TestContainerRegistryAuthDryRun(t *testing.T) {
	t.Parallel()

	prov := newProvider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("ContainerRegistryAuth"),
		Properties: property.NewMap(map[string]property.Value{
			"name":     property.New("dockerhub"),
			"username": property.New("myuser"),
			"password": property.New("mypassword"),
		}),
		DryRun: true,
	})

	require.NoError(t, err)
	require.Equal(t, "dockerhub", response.Properties.Get("name").AsString())
	require.Equal(t, "myuser", response.Properties.Get("username").AsString())
}

// urn is a helper function to build an urn for running integration tests.
func urn(typ string) resource.URN {
	return resource.NewURN("stack", "proj", "",
		tokens.Type("test:index:"+typ), "name")
}

// newProvider creates a test server.
func newProvider(t *testing.T) integration.Server {
	s, err := integration.NewServer(
		context.Background(),
		xyz.Name,
		semver.MustParse("0.1.0"),
		integration.WithProvider(xyz.Provider()),
	)
	require.NoError(t, err)
	return s
}
