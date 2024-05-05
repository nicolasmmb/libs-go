package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func AddTraceIdHeader(ctx *gin.Context) {
	ctx.Header("X-Trace-ID", GetTraceID(ctx.Request.Context()))
	ctx.Next()
}
func GetTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID()
		return traceID.String()
	}
	println("TraceID not found")
	return ""
}
