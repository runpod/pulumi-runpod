# Changelog

## [0.1.5](https://github.com/runpod/pulumi-runpod/compare/v0.1.4...v0.1.5) (2026-03-13)


### Bug Fixes

* populate go.sum with transitive dependency checksums ([c50a236](https://github.com/runpod/pulumi-runpod/commit/c50a236535ef4a336999edc86006136cdebc21d0))
* remove provider:secret from apiKey (SST compatibility) ([83349f2](https://github.com/runpod/pulumi-runpod/commit/83349f2fee9ad4e7ec8d633d51ec0618cb6a2226))
* remove provider:secret from apiKey to fix SST provider replacement ([8cefe8d](https://github.com/runpod/pulumi-runpod/commit/8cefe8d1c4fe49daaa3111293c7280808110a3c8))
* run go mod tidy in Go example test after replace directive change ([767d30d](https://github.com/runpod/pulumi-runpod/commit/767d30dec17eba14d2e9fefb5098c2c7574e7c09))
* update SDK and example go.mod to match provider dependency versions ([2e17ec9](https://github.com/runpod/pulumi-runpod/commit/2e17ec98ed03edebb9ad52a99c67083fc0902947))
* upgrade pulumi-go-provider to v1.3.0 ([3ca3b0e](https://github.com/runpod/pulumi-runpod/commit/3ca3b0e0ea0471d1375f4273f781123a35186d9c))
* use language-specific unique template names in examples ([05dcc0d](https://github.com/runpod/pulumi-runpod/commit/05dcc0d26537c67e2db41270e831ac05585dd7c6))
* use unique template names in examples to prevent CI test collisions ([cd927ed](https://github.com/runpod/pulumi-runpod/commit/cd927ed6c1a3082d16c2cae4298b4c2564467eff))

## [0.1.4](https://github.com/runpod/pulumi-runpod/compare/v0.1.3...v0.1.4) (2026-03-13)


### Bug Fixes

* add NPM_TOKEN fallback for npm publish ([f42d464](https://github.com/runpod/pulumi-runpod/commit/f42d46410016f9b3137d9cb18e49c9b78ce885fb))
* address review feedback ([cdffa0c](https://github.com/runpod/pulumi-runpod/commit/cdffa0c9986ba95d9499d31c6f66bb53952d53b5))
* follow standard provider pattern for apiKey config ([434f828](https://github.com/runpod/pulumi-runpod/commit/434f828bee6e0ba59e134e3029d054826b1f83fd))
* initialize shared client in Configure, validate API key, use pointer receivers ([4655ab0](https://github.com/runpod/pulumi-runpod/commit/4655ab0e2da4fd4776b7617497a486d0544da898))
* make Delete idempotent across all resources ([03e5780](https://github.com/runpod/pulumi-runpod/commit/03e5780eecb748cb7186559079ac3030ceeea9fe))
* remove provider:secret from apiKey + use RELEASE_PAT ([b21d5ba](https://github.com/runpod/pulumi-runpod/commit/b21d5ba6c56a842bbb2cf29f86d66b66827671af))
* remove provider:secret from apiKey and use RELEASE_PAT for release-please ([cf2b681](https://github.com/runpod/pulumi-runpod/commit/cf2b681dcf93f045605094165632467203f9a7bd))

## [0.1.3](https://github.com/runpod/pulumi-runpod/compare/v0.1.2...v0.1.3) (2026-03-12)


### Bug Fixes

* accept version config key to prevent provider replacement ([42fb55b](https://github.com/runpod/pulumi-runpod/commit/42fb55b65c2e1b02b731697bf2c3d33388fce8e7))
* bump pulumi-go-provider to v1.1.1 and fix schema issues ([e58b8b7](https://github.com/runpod/pulumi-runpod/commit/e58b8b7e6390f9d76326d55cc26c22ffb3f29165))
* update circl transitive dep to v1.6.1 in SDK and example go.mods ([dbfd8c6](https://github.com/runpod/pulumi-runpod/commit/dbfd8c63cac738de5d8134532d118ee3d642f813))

## [0.1.2](https://github.com/runpod/pulumi-runpod/compare/v0.1.1...v0.1.2) (2026-03-12)


### Bug Fixes

* use npm OIDC trusted publishing ([4b81758](https://github.com/runpod/pulumi-runpod/commit/4b81758c0f54b7164e3230559f334cc437ecbcc6))
* use npm OIDC trusted publishing instead of NPM_TOKEN ([936ed88](https://github.com/runpod/pulumi-runpod/commit/936ed88f7a753064315e1f740aa7572c2e06c4ca))

## [0.1.1](https://github.com/runpod/pulumi-runpod/compare/v0.1.0...v0.1.1) (2026-03-12)


### Bug Fixes

* make gpuTypeId optional for CPU pod support ([5ee962d](https://github.com/runpod/pulumi-runpod/commit/5ee962d3194d4bd63db810df9e5fa290bf2de972))
* make gpuTypeId optional for CPU pod support ([e95d820](https://github.com/runpod/pulumi-runpod/commit/e95d820a3dad8017d0c6eec7a941e1681838582d))
* rename npm package from @runpod/pulumi to pulumi-runpod ([5a35623](https://github.com/runpod/pulumi-runpod/commit/5a356237ca32812dad5be2568943a14cb3294607))
* rename npm package to pulumi-runpod ([e1d961e](https://github.com/runpod/pulumi-runpod/commit/e1d961e1558e179733f6d62529787a0d6bda4255))
* sync release-please manifest and config ([e5bd197](https://github.com/runpod/pulumi-runpod/commit/e5bd1977ec0bba97ef5736ba39a12875a612b133))
* sync release-please manifest to v0.1.0 and add packages config ([90f75ba](https://github.com/runpod/pulumi-runpod/commit/90f75ba795eb3ce71bfd0fb8f0ab95fd2dbbf309))

## v0.1.0-alpha.1

Initial pre-release of the RunPod Pulumi provider.

### Resources

- `runpod:index:Pod` — Create and manage GPU pods
- `runpod:index:Template` — Create and manage pod templates
- `runpod:index:Endpoint` — Create and manage serverless endpoints
- `runpod:index:NetworkVolume` — Create and manage network storage volumes
- `runpod:index:Secret` — Create and manage secrets
- `runpod:index:ContainerRegistryAuth` — Manage container registry authentication

### Functions

- `runpod:index:getGpuTypes` — Query available GPU types with pricing and availability
- `runpod:index:getCPUFlavors` — Query available CPU configurations
- `runpod:index:getDataCenters` — Query data centers with GPU availability

### SDK Packages

- npm: `@runpod/pulumi`
- PyPI: `pulumi_runpod`
- NuGet: `Pulumi.Runpod`
- Go: `github.com/runpod/pulumi-runpod/sdk/go/runpod`
