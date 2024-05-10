package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func newMeterProvider() (*metric.MeterProvider, error) {
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("otel-test"),
			semconv.ServiceVersion("0.1.0"),
			semconv.ServiceInstanceID(uuid.NewString()),
		))
	if err != nil {
		return nil, err
	}

	me, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint("localhost:9090"),
		otlpmetrichttp.WithURLPath("/api/v1/otlp/v1/metrics"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(me,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	), nil
}
