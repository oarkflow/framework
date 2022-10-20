package middleware

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/sujit-baniya/framework/contracts/http"
)

const (
	OpentracingTracer = "opentracing_tracer"
	OpentracingCtx    = "opentracing_ctx"
)

func Opentracing(tracer opentracing.Tracer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var parentSpan opentracing.Span

		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Headers()))
		if err != nil {
			parentSpan = tracer.StartSpan(ctx.Path())
			defer parentSpan.Finish()
		} else {
			parentSpan = opentracing.StartSpan(
				ctx.Path(),
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
			defer parentSpan.Finish()
		}

		ctx.WithValue(OpentracingTracer, tracer)
		ctx.WithValue(OpentracingCtx, opentracing.ContextWithSpan(context.Background(), parentSpan))
		return ctx.Next()
	}
}
