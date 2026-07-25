---
title: Runpod Installation & Configuration
meta_desc: Information on how to install the Runpod provider for Pulumi.
layout: installation
---

## Installation

The Runpod provider is available as a package in all Pulumi languages:

* JavaScript/TypeScript: [`@runpod/pulumi`](https://www.npmjs.com/package/@runpod/pulumi)
* Python: [`pulumi_runpod`](https://pypi.org/project/pulumi_runpod/)
* Go: [`github.com/runpod/pulumi-runpod/sdk/go/runpod`](https://pkg.go.dev/github.com/runpod/pulumi-runpod/sdk/go/runpod)
* .NET: [`Pulumi.Runpod`](https://www.nuget.org/packages/Pulumi.Runpod)

## Configuration

The following configuration options are available:

* `runpod:apiKey` (required, secret) - The Runpod API key for authentication. Can also be set via the `RUNPOD_API_KEY` environment variable.
* `runpod:apiUrl` (optional) - The Runpod API URL. Defaults to `https://api.runpod.io/graphql`. Can also be set via the `RUNPOD_API_URL` environment variable.

### Setting your API key

Your Runpod API key can be created from the [Runpod Console](https://www.runpod.io/console/user/settings).

Set it as a Pulumi secret:

```bash
pulumi config set --secret runpod:apiKey YOUR_API_KEY
```

Or use the environment variable:

```bash
export RUNPOD_API_KEY=YOUR_API_KEY
```
