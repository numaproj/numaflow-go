# Sink UDF with OpenTelemetry tracing

A log sink that emits a `user.persist` span per message under the Numaflow
platform's per-message `numaflow.{topology}.sink.write` span.

## What it shows

The Numaflow data plane injects the platform's per-stage trace context into the
message's `sys_metadata["tracing_udf"]` field (W3C traceparent + optional
tracestate) before invoking the sink UDF. This example:

1. Sets up an OTLP gRPC tracer in the sink process (`initTracer`).
2. Extracts the platform parent context per incoming message
   (`extractTraceContext`).
3. Starts a child span (`user.persist`) so user-defined persistence work — a
   DB write, HTTP POST, etc. — appears nested under the platform's sink span
   in the same trace.


## Required environment variables

These should be set on the **Numaflow vertex container** (via
`containerTemplate.env`) AND on the UDF container — they govern both the
platform spans and the UDF's own exporter:

| Variable | Purpose |
| --- | --- |
| `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` (or `OTEL_EXPORTER_OTLP_ENDPOINT`) | OTLP gRPC endpoint, e.g. `http://otel-collector.observability.svc:4317`. The UDF strips the `http://` scheme automatically. |
| `OTEL_SERVICE_NAME` | Optional. Defaults to `numaflow-udf` in this example. |

When neither endpoint variable is set, `initTracer` returns a no-op and the
example continues to function as a plain log sink.

## Sampling

The example registers `ParentBased(AlwaysSample)`. That means the sink honours
whatever sampling decision the upstream platform made:

- If the platform sampled this trace in, the sink's `user.persist` span is
  exported.
- If the platform sampled this trace out, the sink's span is also dropped.

This is the right default — it prevents orphaned sink spans when the platform
has decided not to record the trace, and keeps cross-process traces consistent.

## Expected trace tree

For a MonoVertex `source → map → sink (this UDF)`:

```
numaflow.vertex.process
├── numaflow.monovertex.source.dispatch
├── numaflow.monovertex.map
└── numaflow.monovertex.sink.write
    └── user.persist                       ← emitted by this example
```

For a Pipeline, the `user.persist` span appears under the sink pod's
`numaflow.pipeline.sink.write` span, which is itself a child of the sink pod's
`numaflow.vertex.process` (parented across the ISB to the upstream pod's
`vertex.process`).

See the platform-side tracing guide for the full hierarchy:
<https://numaflow.numaproj.io/user-guide/reference/tracing/>

## Quick start

```bash
make image
```

Then reference the image in your Pipeline / MonoVertex spec:

```yaml
sink:
  udsink:
    container:
      image: quay.io/numaio/numaflow-go/sink-tracing:stable
      imagePullPolicy: IfNotPresent
containerTemplate:
  env:
    - name: OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
      value: "http://simplest-jaeger-collector.default.svc:4317"
    - name: OTEL_TRACES_SAMPLER
      value: "parentbased_always_on"   # for debugging only
    - name: OTEL_SERVICE_NAME
      value: "numaflow-sink-tracing"
```