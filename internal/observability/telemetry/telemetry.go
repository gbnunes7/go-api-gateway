package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const serviceName = "api-gateway"

var TracerProvider *sdktrace.TracerProvider

func otlpEndpoint() string {
	if e := os.Getenv("OTLP_ENDPOINT"); e != "" {
		return e
	}
	return "localhost:4318"
}

func Init(ctx context.Context) error {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint()),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("telemetry: create otlp exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		return fmt.Errorf("telemetry: create resource: %w", err)
	}

	TracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(TracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

func Shutdown(ctx context.Context) error {
	if TracerProvider == nil {
		return nil
	}
	return TracerProvider.Shutdown(ctx)
}
