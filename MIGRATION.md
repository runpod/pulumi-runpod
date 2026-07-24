# Migration Guide

## Upgrading to v2.0.0

v2.0.0 brings the provider in line with the current Runpod GraphQL API and fixes
drift that had accumulated against it. Most changes are additive; the breaking
changes are small and listed first.

### Breaking changes

#### Endpoint: `modelName` removed
The Runpod API's `EndpointInput` no longer accepts `modelName` — it was a dead
field that the provider was still sending. Use `modelReferences` instead.

```diff
  new runpod.Endpoint("ep", {
    name: "my-endpoint",
-   modelName: "my-model",
+   modelReferences: ["my-model"],
  });
```

#### Pod: `desiredStatus` is now a settable input
Previously `desiredStatus` was a read-only output. It is now an optional input
that declaratively controls the pod's run state:

- Set `desiredStatus: "EXITED"` to stop (pause) the pod **in place**.
- Set `desiredStatus: "RUNNING"` to resume it.
- Leave it unset to not manage the run state (previous default behavior).

If you previously *read* `pod.desiredStatus` as an output, it is now only
populated when you set it as an input. No action is needed if you didn't use it.

### New features (non-breaking)

| Resource / Function | New field(s) |
|---|---|
| `Endpoint` | `workersPFBTarget`, `requestTTL`, `dataCenterIds` (settable inputs) |
| `Template` | `minVram`, `minRam` |
| `getDataCenters` | `cpuAvailability` (per–data-center CPU flavor availability) |
| `Pod` | `desiredStatus` (stop/resume — see above) |

`dataCenterIds` is the structured replacement for the legacy comma-separated
`locations` string on `Endpoint`; `locations` still works.

---

## Coming from the legacy `pulumi-runpod-native` provider

If you used the older provider published from `runpod/pulumi-runpod-native`
(npm `@runpod/pulumi-runpod`, PyPI `runpodinfra`), note these resource changes:

| Legacy | Now |
|---|---|
| `runpod:NetworkStorage` | `runpod:NetworkVolume` |
| `runpod:Gpu` (resource) | `getGpuTypes` (function) |
| `runpod:DataCenter` (resource) | `getDataCenters` (function) |

A detailed legacy-import guide will accompany the registry cutover.
