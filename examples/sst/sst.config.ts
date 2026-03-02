/// <reference path="./.sst/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "runpod-sst-test",
      removal: "remove",
      home: "local",
    };
  },

  async run() {
    const runpod = await import("@runpod/pulumi");

    // Set RUNPOD_API_KEY env var before deploying.
    const provider = new runpod.Provider("runpod", {
      apiKey: process.env.RUNPOD_API_KEY,
    });

    const opts = { provider };

    // --- Free resources ---

    // 1. Template (serverless, free)
    const template = new runpod.Template(
      "sst-test-template",
      {
        name: "sst-test-template",
        imageName: "runpod/serverless-hello-world:latest",
        containerDiskInGb: 5,
        volumeInGb: 0,
        isServerless: true,
        env: {},
      },
      opts,
    );

    // 2. Secret (free)
    const secret = new runpod.Secret(
      "sst-test-secret",
      {
        name: "sst-test-secret",
        value: process.env.SECRET_VALUE ?? "REPLACE_ME",
        description: "SST integration test secret",
      },
      opts,
    );

    // 3. Container registry auth (free)
    const registryAuth = new runpod.ContainerRegistryAuth(
      "sst-test-registry",
      {
        name: "sst-test-registry",
        username: "testuser",
        password: process.env.REGISTRY_PASSWORD ?? "REPLACE_ME",
      },
      opts,
    );

    // 4. Endpoint with 0 workers (free — no workers = no cost)
    const endpoint = new runpod.Endpoint(
      "sst-test-endpoint",
      {
        name: "sst-test-endpoint",
        templateId: template.templateId,
        gpuIds: "AMPERE_16",
        workersMin: 0,
        workersMax: 0,
        idleTimeout: 5,
      },
      opts,
    );

    // 5. Network volume (10GB min in MFS-1 data center)
    const networkVolume = new runpod.NetworkVolume(
      "sst-test-volume",
      {
        name: "sst-test-volume",
        size: 10,
        dataCenterId: "MFS-1",
      },
      opts,
    );

    // 6. Pod (RTX 4090 secure cloud in dev)
    const pod = new runpod.Pod(
      "sst-test-pod",
      {
        name: "sst-test-pod",
        gpuTypeId: "NVIDIA GeForce RTX 4090",
        imageName:
          "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04",
        gpuCount: 1,
        volumeInGb: 0,
        containerDiskInGb: 10,
        cloudType: "SECURE",
      },
      opts,
    );

    return {
      templateId: template.templateId,
      secretId: secret.secretId,
      registryAuthId: registryAuth.registryAuthId,
      endpointId: endpoint.endpointId,
      networkVolumeId: networkVolume.networkVolumeId,
      podId: pod.podId,
    };
  },
});
