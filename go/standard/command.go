package main

import (
	"context"
	"os/exec"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func main() {
	tracer, closer, err := config.Configuration{
		ServiceName: "test",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  "127.0.0.1:5775",
		},
	}.NewTracer()
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	// span, ctx := opentracing.StartSpanFromContext(context.Background(), "test")
	// defer span.Finish()
	ctx := context.Background()
	exec.CommandContext(ctx, "sleep", "5").Run()
}
