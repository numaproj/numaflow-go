# Map UDF with OpenTelemetry tracing

A pass-through map UDF that emits a `user.work` span under the Numaflow
platform's per-message `numaflow.{topology}.map` span.

## What it shows

The Numaflow data plane injects the platform's per-stage trace context into the
message's `sys_metadata["tracing_udf"]` field (W3C traceparent + optional
tracestate) before invoking the UDF. This example:

1. Sets up an OTLP gRPC tracer in the UDF process (`initTracer`).
2. Extracts the platform parent context from each incoming message
   (`extractTraceContext`).
3. Starts a child span (`user.work`) so user-defined work appears nested under
   the platform's stage span in the same trace.

## Required environment variables

These should be set on the **Numaflow vertex container** (via
`containerTemplate.env`) AND on the UDF container — they govern both the
platform spans and the UDF's own exporter:

| Variable | Purpose |
| --- | --- |
| `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` (or `OTEL_EXPORTER_OTLP_ENDPOINT`) | OTLP gRPC endpoint, e.g. `http://otel-collector.observability.svc:4317`. The UDF strips the `http://` scheme automatically. |
| `OTEL_SERVICE_NAME` | Optional. Defaults to `numaflow-udf` in this example. |

When neither endpoint variable is set, `initTracer` returns a no-op and the
example continues to function as a plain pass-through map.

## Sampling

The example registers `ParentBased(AlwaysSample)`. That means the UDF honours
whatever sampling decision the upstream platform made:

- If the platform sampled this trace in, the UDF's `user.work` span is exported.
- If the platform sampled this trace out, the UDF's span is also dropped.

This is the right default — it prevents orphaned UDF spans when the platform
has decided not to record the trace, and keeps cross-process traces consistent.

## Expected trace tree

For a MonoVertex `source → flatmap (this UDF) → sink`:

```
numaflow.vertex.process
├── numaflow.monovertex.source.dispatch
├── numaflow.monovertex.map
│   └── user.work                          ← emitted by this example
└── numaflow.monovertex.sink.write
```

For a Pipeline, the `user.work` span appears under the map pod's
`numaflow.pipeline.map` span, which is itself a child of the map pod's
`numaflow.vertex.process` (parented across the ISB to the source pod's
`vertex.process`).

See the platform-side tracing guide for the full hierarchy:
<https://numaflow.numaproj.io/user-guide/reference/tracing/>

## Quick start

```bash
make image
```

Then reference the image in your Pipeline / MonoVertex spec:

```yaml
udf:
  container:
    image: quay.io/numaio/numaflow-go/map-tracing:stable
    imagePullPolicy: IfNotPresent
containerTemplate:
  env:
    - name: OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
      value: "http://simplest-jaeger-collector.default.svc:4317"
    - name: OTEL_TRACES_SAMPLER
      value: "parentbased_always_on"   # for debugging only
    - name: OTEL_SERVICE_NAME
      value: "numaflow-map-tracing"
```