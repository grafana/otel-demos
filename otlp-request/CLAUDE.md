# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an OpenTelemetry (OTel) demo application that generates and sends OTLP (OpenTelemetry Protocol) metric requests to test metric ingestion pipelines. The application is part of the larger `grafana/otel-demos` repository containing various OTel demonstration applications.

## Architecture

- **Single-file Go application** (`main.go`) that constructs OTLP JSON payloads
- **Self-contained**: No external dependencies beyond Go standard library
- **HTTP client**: Sends POST requests to `/otlp/v1/metrics` endpoint with proper headers
- **Metric structure**: Implements complete OTLP metric data model including:
  - Resource attributes (service.name, service.version)
  - Scope metrics with instrumentation metadata
  - Sum metrics with data points, attributes, and exemplars
  - Exemplar data with trace correlation (spanId, traceId)

## Key Components

### Data Structures
- `OTLPRequest` - Root structure containing resource metrics
- `ResourceMetrics` - Service-level metadata and scope metrics
- `ScopeMetrics` - Instrumentation library metrics
- `Metric` - Individual metric definition with Sum data
- `NumberDataPoint` - Time-series data with exemplars
- `Exemplar` - Individual exemplar with trace context

### Request Flow
The `sendRequest()` function creates two test scenarios:
1. Normal timing: exemplar timestamp matches data point timestamp  
2. Out-of-order timing: exemplar timestamp is 1 second before data point timestamp

Target endpoint: `http://localhost:8080/otlp/v1/metrics` with headers:
- `X-Scope-OrgID: 1000` (Grafana multi-tenancy)
- `Content-Type: application/json`

## Development Commands

### Build and Run
```bash
go run main.go        # Compile and execute
go build             # Create binary executable
go build -o otlp-request  # Create named binary
```

### Development
```bash
go fmt               # Format code
go vet               # Static analysis
go mod tidy          # Clean up dependencies
```

### Testing
```bash
go test              # Run tests (none currently exist)
```

## Usage Context

This application is designed to test OTLP metric ingestion endpoints, particularly:
- Mimir/Grafana metric storage systems
- OTLP gateway implementations  
- Metric pipeline validation with exemplars and trace correlation
- Out-of-order data handling scenarios
