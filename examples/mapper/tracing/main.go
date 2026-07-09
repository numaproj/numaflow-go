// Tracing-aware map UDF example.
//
// The Numaflow data plane injects the platform's per-stage span context into
// `sys_metadata["tracing_udf"]` (W3C traceparent + optional tracestate) before
// calling the UDF. This example shows how to:
//
//   1. Initialise an OTLP gRPC tracer in the UDF process so it can export spans.
//   2. Extract the platform parent context from the message's system metadata.
//   3. Create a child span (`user.work`) under the platform's `numaflow.{topology}.map`
//      span so user-defined work shows up nested in the same trace.
//
// Required environment variables (set by the Pipeline/MonoVertex containerTemplate):
//
//   OTEL_EXPORTER_OTLP_TRACES_ENDPOINT  (or the generic OTEL_EXPORTER_OTLP_ENDPOINT)
//   OTEL_SERVICE_NAME                   (optional; defaults to "numaflow-udf")
//
// When neither endpoint variable is set the tracer init is a no-op and the
// example still runs as a plain pass-through map.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/numaproj/numaflow-go/pkg/mapper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// tracingUDFGroup is the sys_metadata group the Numaflow data plane uses to
// inject the platform parent context for the current stage.
const tracingUDFGroup = "tracing_udf"

// initTracer wires an OTLP gRPC tracer provider and the W3C propagator. It
// returns a shutdown function the caller must invoke before exiting so spans
// are flushed.
//
// If neither OTEL_EXPORTER_OTLP_TRACES_ENDPOINT nor OTEL_EXPORTER_OTLP_ENDPOINT
// is set, this is a no-op: a no-op tracer remains globally registered and span
// creation in the UDF body is essentially free.
func initTracer(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if endpoint == "" {
		endpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}
	if endpoint == "" {
		log.Println("[tracing] OTLP endpoint not set; UDF spans will be no-ops")
		return func(context.Context) error { return nil }, nil
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "numaflow-udf"
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(stripScheme(endpoint)),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		// Honour the upstream platform's sampling decision so we never produce
		// orphaned UDF spans when the platform sampled the trace out.
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
		sdktrace.WithResource(resource.NewSchemaless(
			semconv.ServiceName(serviceName),
		)),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	log.Printf("[tracing] OTLP exporter configured: endpoint=%s service=%s", endpoint, serviceName)
	return provider.Shutdown, nil
}

// stripScheme drops a leading "http://" or "https://" so the user can configure
// the endpoint with the same string the Rust data plane accepts (which keeps
// the scheme). The OTel Go gRPC exporter's WithEndpoint expects host:port only.
func stripScheme(endpoint string) string {
	for _, prefix := range []string{"http://", "https://"} {
		if len(endpoint) >= len(prefix) && endpoint[:len(prefix)] == prefix {
			return endpoint[len(prefix):]
		}
	}
	return endpoint
}

// extractTraceContext reads the W3C traceparent/tracestate the platform wrote
// into sys_metadata["tracing_udf"] and returns a context whose current span is
// the platform-side `numaflow.{topology}.map` span. Spans started with this
// context become children of the platform span in the trace tree.
//
// Safe to call when tracing is disabled or when no parent context is present:
// the returned context is then the input context unchanged.
func extractTraceContext(ctx context.Context, d mapper.Datum) context.Context {
	sysMeta := d.SystemMetadata()
	if sysMeta == nil {
		return ctx
	}
	traceparent := string(sysMeta.Value(tracingUDFGroup, "traceparent"))
	if traceparent == "" {
		return ctx
	}
	carrier := propagation.MapCarrier{"traceparent": traceparent}
	if ts := string(sysMeta.Value(tracingUDFGroup, "tracestate")); ts != "" {
		carrier["tracestate"] = ts
	}
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

// TracedMapper is a pass-through map UDF that emits a `user.work` span under
// the platform's per-message map span. Replace the body of Map with your real
// work; any further span.Start calls hang off `user.work`.
type TracedMapper struct{}

func (m *TracedMapper) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	ctx = extractTraceContext(ctx, d)
	_, span := otel.Tracer("numaflow-go-example/mapper-tracing").Start(ctx, "user.work")
	defer span.End()

	// Real UDF work would go here. We just pass the value through so the
	// example stays focused on the tracing wiring.
	return mapper.MessagesBuilder().Append(mapper.NewMessage(d.Value()).WithKeys(keys))
}

func main() {
	ctx := context.Background()
	shutdown, err := initTracer(ctx)
	if err != nil {
		log.Panic("Failed to initialise tracer: ", err)
	}
	defer func() {
		sCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if shutdownErr := shutdown(sCtx); shutdownErr != nil {
			log.Printf("[tracing] tracer shutdown error: %v", shutdownErr)
		}
	}()

	if err := mapper.NewServer(&TracedMapper{}).Start(ctx); err != nil {
		log.Panic("Failed to start map server: ", err)
	}
}
