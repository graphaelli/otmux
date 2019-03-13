package common

import (
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ZipkinTracer() (opentracing.Tracer, io.Closer) {
	logger := log.New(os.Stderr, "[zipkin] ", LogFmt)
	loggerFunc := zipkintracer.LoggerFunc(func(v ...interface{}) error { logger.Print(v); return nil })
	/*
	collector, err := zipkintracer.NewHTTPCollector(
		"http://localhost:9411/",
		zipkintracer.HTTPLogger(loggerFunc),
		zipkintracer.HTTPBatchInterval(10*time.Millisecond),
		zipkintracer.HTTPBatchSize(1),
	)
	*/
	collector, err := zipkintracer.NewScribeCollector(
		"localhost:9410",
		10*time.Second,
		zipkintracer.ScribeLogger(loggerFunc),
		zipkintracer.ScribeBatchSize(1),
		zipkintracer.ScribeBatchInterval(10*time.Millisecond),
	)
	if err != nil {
		// unreachable?
		panic(err)
	}
	tracer, err := zipkintracer.NewTracer(
		zipkintracer.NewRecorder(collector, true, "http://0.0.0.0:0", filepath.Base(os.Args[0])),
		zipkintracer.WithLogger(loggerFunc),
	)
	if err != nil {
		// unreachable?
		panic(err)
	}
	return tracer, collector
}
