package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// InitTracer initializes Jaeger tracing
func InitTracer(serviceName, jaegerEndpoint string) (func(context.Context) error, error) {
	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, err
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	// Get tracer
	tracer = tp.Tracer(serviceName)

	log.Printf("Jaeger tracing initialized: %s", jaegerEndpoint)

	// Return shutdown function
	return tp.Shutdown, nil
}

// GetTracer returns the global tracer
func GetTracer() trace.Tracer {
	return tracer
}

// StartSpan starts a new span
func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	if tracer == nil {
		// Return no-op span if tracer not initialized
		return ctx, trace.SpanFromContext(ctx)
	}
	return tracer.Start(ctx, spanName)
}