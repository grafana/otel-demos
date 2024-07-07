package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
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
			attribute.String("cloud.availability_zone", "A"),
			attribute.String("cloud.region", "Outer space"),
			attribute.String("container.name", "Docker"),
			attribute.String("deployment.environment", "Gaia"),
			attribute.String("k8s.cluster.name", "Cluster"),
			attribute.String("k8s.container.name", "Docker"),
			attribute.String("k8s.cronjob.name", "The job"),
			attribute.String("k8s.daemonset.name", "No daemon"),
			attribute.String("k8s.deployment.name", "The deployment"),
			attribute.String("k8s.job.name", "The job"),
			attribute.String("k8s.namespace.name", "The namespace"),
			attribute.String("k8s.pod.name", "The pod"),
			attribute.String("k8s.replicaset.name", "The RS"),
			attribute.String("k8s.statefulset.name", "The StatefulSet"),
		))
	if err != nil {
		return nil, err
	}

	me, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpointURL("http://localhost:9090"),
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
