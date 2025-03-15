package core

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewTracer initializes an OpenTelemetry trace provider configured with AWS X-Ray ID generator,
// sets it as the global trace provider, and configures X-Ray propagation.
//
// It returns a shutdown function to clean up resources, which is a closure
// that should be called when the application terminates.
//
// Parameters:
//   - ctx: Context for tracer initialization
//   - appEnv: Application environment (e.g., "production", "staging")
//   - otlpEndpoint: Endpoint for the OpenTelemetry Collector (e.g., "collector:4317")
//   - service: Name of the service
//   - version: Version of the service being traced
//
// Returns:
//   - func(context.Context) error: Shutdown function to clean up the tracer
//   - error: Any error encountered during setup
func NewTracer(ctx context.Context, appEnv, otlpEndpoint, service, version string) (func(context.Context) error, error) {
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %v", err)
	}

	// Create and register trace provider with X-Ray ID generator
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(service),
				semconv.ServiceVersionKey.String(version),
				semconv.DeploymentEnvironmentKey.String(appEnv),
			),
		),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Default to sampling all traces
	)

	// Set global trace provider and X-Ray propagator
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(xray.Propagator{})

	return traceProvider.Shutdown, nil
}
