# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an OpenTelemetry (OTel) metrics demonstration project written in Go that showcases manual OTel instrumentation. The application creates a single counter metric named "test.counter" and periodically sends metrics data via OTLP HTTP to a local endpoint.

## Architecture

The codebase consists of two main files:
- `main.go`: Contains the main application logic that creates and increments a counter metric every second
- `provider.go`: Contains the `newMeterProvider()` function that sets up the OTLP HTTP exporter and configures resource attributes including service metadata, cloud, container, and Kubernetes attributes

The application sends metrics every 3 seconds to `http://localhost:9090/api/v1/otlp/v1/metrics` using the OTLP HTTP exporter.

## Development Commands

### Build
```bash
go build
```

### Run
```bash
go run .
```

### Format Code
```bash
go fmt ./...
```

### Test
```bash
go test ./...
```
Note: Currently no test files exist in the project.

### Dependency Management
```bash
# Update dependencies
go mod tidy

# Verify dependencies
go mod verify

# Download dependencies
go mod download
```

## Key Dependencies

- `go.opentelemetry.io/otel` - Core OpenTelemetry Go SDK
- `go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp` - OTLP HTTP metrics exporter
- `go.opentelemetry.io/otel/sdk/metric` - Metrics SDK
- `github.com/google/uuid` - UUID generation for service instance ID

## Configuration

The application is configured to:
- Send metrics to localhost:9090 via HTTP (insecure)
- Use OTLP metrics endpoint: `/api/v1/otlp/v1/metrics`
- Export metrics every 3 seconds
- Include comprehensive resource attributes for cloud, container, and Kubernetes environments