# otel-exporter

This is a demonstration of [manual OTel instrumentation](https://opentelemetry.io/docs/instrumentation/go/manual/).
The demo creates one metric, a counter named "test.counter". The counter is given a unit: "By".
Every second, the counter is incremented, and every 3 seconds metrics get written over OTLP to 
http://localhost:8000/otlp/v1/metrics.
