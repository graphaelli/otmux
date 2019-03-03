package common

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmot"
)

type flusher interface {
	Flush(abort <-chan struct{})
}

type flushCloser struct {
	flusher
}

func (f flushCloser) Close() error {
	f.Flush(nil)
	return nil
}

func ElasticTracer() (opentracing.Tracer, io.Closer) {
	elasticTracer := apm.DefaultTracer
	tracer := apmot.New(apmot.WithTracer(elasticTracer))
	return tracer, flushCloser{flusher: elasticTracer}
}
