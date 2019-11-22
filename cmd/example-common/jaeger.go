package common

import (
	"io"
	"os"
	"path/filepath"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func JaegerTracer() (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Headers: &jaeger.HeadersConfig{
			JaegerDebugHeader:        jaeger.JaegerDebugHeader,
			JaegerBaggageHeader:      jaeger.JaegerBaggageHeader,
			TraceBaggageHeaderPrefix: "jaeger-ctx-",
			TraceContextHeaderName:   "jaeger-trace-id",
		},
		ServiceName: filepath.Base(os.Args[0]),
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		// unreachable?
		panic(err)
	}
	return tracer, closer
}

func init() {
	RegisterTracer("jaeger", JaegerTracer)
}