package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type OTLPRequest struct {
	ResourceMetrics []ResourceMetrics `json:"resourceMetrics"`
}

type ResourceMetrics struct {
	Resource     Resource       `json:"resource"`
	ScopeMetrics []ScopeMetrics `json:"scopeMetrics"`
}

type Resource struct {
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string         `json:"key"`
	Value AttributeValue `json:"value"`
}

type AttributeValue struct {
	StringValue string `json:"stringValue,omitempty"`
}

type ScopeMetrics struct {
	Scope   Scope    `json:"scope"`
	Metrics []Metric `json:"metrics"`
}

type Scope struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Metric struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
	Sum         *Sum   `json:"sum,omitempty"`
}

type Sum struct {
	DataPoints             []NumberDataPoint `json:"dataPoints"`
	AggregationTemporality int               `json:"aggregationTemporality"`
	IsMonotonic            bool              `json:"isMonotonic"`
}

type NumberDataPoint struct {
	Attributes        []Attribute `json:"attributes"`
	StartTimeUnixNano uint64      `json:"startTimeUnixNano,string"`
	TimeUnixNano      uint64      `json:"timeUnixNano,string"`
	AsInt             int64       `json:"asInt"`
	Exemplars         []Exemplar  `json:"exemplars,omitempty"`
}

type Exemplar struct {
	FilteredAttributes []Attribute `json:"filteredAttributes,omitempty"`
	TimeUnixNano       uint64      `json:"timeUnixNano,string"`
	AsInt              int64       `json:"asInt"`
	SpanID             string      `json:"spanId,omitempty"`
	TraceID            string      `json:"traceId,omitempty"`
}

func main() {
	now := time.Now()
	startTime := now.Add(-1 * time.Minute)

	if err := sendRequest(startTime, now, now); err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	// Send a request with an out-of-order exemplar, to test this scenario.
	if err := sendRequest(startTime, now.Add(time.Second), now.Add(-time.Second)); err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
}

func sendRequest(startTime, dataPointTime, exemplarTime time.Time) error {
	const (
		serviceName = "demo-service"
		version     = "1.0.0"
	)
	attrs := []Attribute{
		{
			Key: "method",
			Value: AttributeValue{
				StringValue: "GET",
			},
		},
		{
			Key: "status",
			Value: AttributeValue{
				StringValue: "200",
			},
		},
	}

	request := OTLPRequest{
		ResourceMetrics: []ResourceMetrics{
			{
				Resource: Resource{
					Attributes: []Attribute{
						{
							Key: "service.name",
							Value: AttributeValue{
								StringValue: serviceName,
							},
						},
						{
							Key: "service.version",
							Value: AttributeValue{
								StringValue: version,
							},
						},
					},
				},
				ScopeMetrics: []ScopeMetrics{
					{
						Scope: Scope{
							Name:    "demo-instrumentation",
							Version: "1.0.0",
						},
						Metrics: []Metric{
							{
								Name:        "http_requests_total",
								Description: "Total number of HTTP requests",
								Unit:        "1",
								Sum: &Sum{
									DataPoints: []NumberDataPoint{
										{
											Attributes:        attrs,
											StartTimeUnixNano: uint64(startTime.UnixNano()),
											TimeUnixNano:      uint64(dataPointTime.UnixNano()),
											AsInt:             150,
											Exemplars: []Exemplar{
												{
													FilteredAttributes: []Attribute{
														{
															Key: "user_id",
															Value: AttributeValue{
																StringValue: "user123",
															},
														},
													},
													TimeUnixNano: uint64(exemplarTime.UnixNano()),
													AsInt:        1,
													SpanID:       "00f067aa0ba902b7",
													TraceID:      "4bf92f3577b34da6a3ce929d0e0e4736",
												},
											},
										},
									},
									AggregationTemporality: 2, // AGGREGATION_TEMPORALITY_CUMULATIVE
									IsMonotonic:            true,
								},
							},
						},
					},
				},
			},
		},
	}

	jsonData, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %s", err)
	}

	fmt.Println(string(jsonData))

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/otlp/v1/metrics", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %s", err)
	}

	req.Header.Set("X-Scope-OrgID", "1000")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(jsonData)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading HTTP response: %s", err)
	}

	if resp.StatusCode/100 == 2 {
		fmt.Printf("Request was successfully handled with code %d: %s\n", resp.StatusCode, string(body))
	} else {
		fmt.Printf("Request failed with code %d: %s\n", resp.StatusCode, string(body))
	}

	return nil
}
