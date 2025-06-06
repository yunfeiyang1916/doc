package main

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	traceName = "my-otel-demo"
)

func main() {
	traceShutdown := InitTracer()
	defer traceShutdown(context.Background())
	//traceTest()
	HttpServerStart()
}

func InitTracer() func(ctx context.Context) error {
	var url = "http://127.0.0.1:14268/api/traces"
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}
	tp := trace.NewTracerProvider(trace.WithBatcher(exp),
		// 在资源中记录有关此应用程序的信息
		trace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("otel-demo"),
			attribute.String("exporter", "jaeger"),
			attribute.Float64("float", 312.23),
		)))
	otel.SetTracerProvider(tp)
	// 设置全局传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp.Shutdown
}

func traceTest() {
	tr := otel.Tracer(traceName)
	ctx, span := tr.Start(context.Background(), "traceTest")
	span.SetAttributes(attribute.String("my-attribute-test-key", "这是测试attribute"))
	span.AddEvent("my-event")
	traceTestChild(ctx)
	time.Sleep(time.Second)
	fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())
	span.End()
}

func traceTestChild(ctx context.Context) {
	tr := otel.Tracer(traceName)
	_, span := tr.Start(ctx, "traceTestChild")
	span.SetAttributes(attribute.String("my-attribute-test-key", "这是测试attribute"))
	span.AddEvent("my-event")
	time.Sleep(2 * time.Second)
	fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())
	span.End()
}
