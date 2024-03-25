package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func main() {
	meterProvider, err := newMeterProvider()
	if err != nil {
		panic(err)
	}

	meter := meterProvider.Meter("test-meter")

	counter, err := meter.Int64Counter("test.counter", metric.WithDescription("Number of iterations"), metric.WithUnit("By"))
	if err != nil {
		panic(err)
	}

	for {
		counter.Add(context.Background(), 1,
			metric.WithAttributes(semconv.HTTPStatusCode(200)))
		time.Sleep(time.Second)
	}
}
