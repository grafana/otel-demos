package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func newMeterProvider() (*metric.MeterProvider, error) {
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("otel-test"),
			semconv.ServiceVersion("0.1.0"),
		))
	if err != nil {
		return nil, err
	}

	me, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint("localhost:8000"),
		otlpmetrichttp.WithURLPath("/otlp/v1/metrics"),
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
