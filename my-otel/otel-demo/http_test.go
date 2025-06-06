package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel"
)

var (
	serverUrl = "http://localhost:8090"
)

func TestMyHttpClient(t *testing.T) {
	traceShutdown := InitTracer()
	defer traceShutdown(context.Background())
	httpClient := http.Client{}
	req, _ := http.NewRequest("GET", serverUrl+"/my-trace", nil)
	tr := otel.Tracer(traceName)
	_, span := tr.Start(context.Background(), "my-http-client")
	fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())

	req.Header.Add("trace-id", span.SpanContext().TraceID().String())
	req.Header.Add("span-id", span.SpanContext().SpanID().String())
	r, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r)
	span.End()
}

func TestHttpClient(t *testing.T) {
	traceShutdown := InitTracer()
	defer traceShutdown(context.Background())
	httpClient := http.Client{}
	req, _ := http.NewRequest("GET", serverUrl, nil)
	tr := otel.Tracer(traceName)
	ctx, span := tr.Start(context.Background(), "my-http-client")
	fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())

	p := otel.GetTextMapPropagator()
	// 将ctx注入到包裹里(注入到了http请求中的header)
	p.Inject(ctx, propagation.HeaderCarrier(req.Header))
	r, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r)
	span.End()
}
