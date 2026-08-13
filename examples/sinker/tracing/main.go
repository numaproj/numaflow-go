// Tracing-aware sink UDF example.
//
// The Numaflow data plane injects the platform's per-stage span context into
// `sys_metadata["tracing_udf"]` (W3C traceparent + optional tracestate) before
// calling the sink UDF. This example shows how to:
//
//   1. Initialise an OTLP gRPC tracer in the sink process so it can export spans.
//   2. Extract the platform parent context per incoming message.
//   3. Create a child span (`user.persist`) under the platform's
//      `numaflow.{topology}.sink.write` span — typical place to span an
//      external DB write, HTTP POST, or other persistence call.
//
// Required environment variables (set by the Pipeline/MonoVertex containerTemplate):
//
//   OTEL_EXPORTER_OTLP_TRACES_ENDPOINT  (or the generic OTEL_EXPORTER_OTLP_ENDPOINT)
//   OTEL_SERVICE_NAME                   (optional; defaults to "numaflow-udf")
//
// When neither endpoint variable is set the tracer init is a no-op and the
// example continues to function as a plain log sink.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
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
// creation in the sink body is essentially free.
func initTracer(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if endpoint == "" {
		endpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}
	if endpoint == "" {
		log.Println("[tracing] OTLP endpoint not set; sink spans will be no-ops")
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
		// orphaned sink spans when the platform sampled the trace out.
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
// the platform-side `numaflow.{topology}.sink.write` span. Spans started with
// this context become children of the platform sink span in the trace tree.
//
// Safe to call when tracing is disabled or when no parent context is present:
// the returned context is then the input context unchanged.
func extractTraceContext(ctx context.Context, d sinksdk.Datum) context.Context {
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

// TracedSink is a log sink that emits a `user.persist` span per message under
// the platform's per-message sink.write span. Replace the body of the loop
// with your real persistence work (DB write, HTTP POST, etc.) and add nested
// spans there for finer attribution.
type TracedSink struct{}

func (s *TracedSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	tracer := otel.Tracer("numaflow-go-example/sinker-tracing")

	for d := range datumStreamCh {
		msgCtx := extractTraceContext(ctx, d)
		_, span := tracer.Start(msgCtx, "user.persist")

		// Real persistence work would go here. We just log the value so the
		// example stays focused on the tracing wiring.
		fmt.Println("Traced sink:", string(d.Value()))

		span.End()
		result = result.Append(sinksdk.ResponseOK(d.ID()))
	}
	return result
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

	if err := sinksdk.NewServer(&TracedSink{}).Start(ctx); err != nil {
		log.Panic("Failed to start sink server: ", err)
	}
}
