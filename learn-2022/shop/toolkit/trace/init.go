package trace

import (
	"context"

	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// 设置全局trace
func InitTracer(serviceName, url string) func(ctx context.Context) error {
	// 创建Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}
	tp := trace.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(1.0))),
		// 始终确保再生成中批量处理
		trace.WithBatcher(exp),
		// 在资源中记录有关此应用程序的信息
		trace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("exporter", "jaeger"),
			attribute.Float64("float", 312.23),
		)),
	)
	otel.SetTracerProvider(tp)
	// 设置全局传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return exp.Shutdown
}
