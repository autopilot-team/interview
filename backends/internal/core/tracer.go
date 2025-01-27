package core

import (
	"context"
	"crypto/tls"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Package core provides core functionality for the application.

// NewTracer initializes and configures an OpenTelemetry tracer with Axiom export capabilities.
// It sets up tracing with the provided configuration parameters and returns a shutdown function
// that should be called when the application terminates.
//
// Parameters:
//   - ctx: Context for tracer initialization
//   - appEnv: Application environment (e.g., "production", "staging")
//   - apiToken: Axiom API token for authentication
//   - service: Name of the service/dataset in Axiom
//   - version: Version of the service being traced
//
// Returns:
//   - func(context.Context) error: Shutdown function to clean up the tracer
//   - error: Any error encountered during setup
func NewTracer(ctx context.Context, appEnv, apiToken, service, version string) (func(context.Context) error, error) {
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("api.axiom.co"),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization":   "Bearer " + apiToken,
			"X-AXIOM-DATASET": "autopilot",
		}),
		otlptracehttp.WithTLSClientConfig(&tls.Config{}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %v", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(service),
				semconv.ServiceVersionKey.String(version),
				semconv.DeploymentEnvironmentKey.String(appEnv),
			),
		),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return provider.Shutdown, nil
}
