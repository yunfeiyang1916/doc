package main

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func HttpServerStart() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		// 负责span的抽取和生成
		ctx := c.Request.Context()
		p := otel.GetTextMapPropagator()
		spanCtx := p.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
		tr := otel.Tracer(traceName)
		_, span := tr.Start(spanCtx, "http-server")
		fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())
		time.Sleep(time.Second)
		span.End()
		c.JSON(200, "ok")
	})

	r.GET("/my-trace", func(c *gin.Context) {
		// 自定义键值对来提取和生成
		ctx := c.Request.Context()
		//p := otel.GetTextMapPropagator()
		//spanCtx := p.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
		strTraceID := c.GetHeader("trace-id")
		strSpanID := c.GetHeader("span-id")
		traceID, _ := trace.TraceIDFromHex(strTraceID)
		spanID, _ := trace.SpanIDFromHex(strSpanID)
		sCtx := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled, // 这个如果不设置的化，是不会保存的
			Remote:     true,
		})
		spanCtx := trace.ContextWithSpanContext(ctx, sCtx)
		tr := otel.Tracer(traceName)
		_, span := tr.Start(spanCtx, "http-server")
		fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())
		time.Sleep(time.Second)
		span.End()
		c.JSON(200, "ok")
	})
	r.Run(":8090")
}
