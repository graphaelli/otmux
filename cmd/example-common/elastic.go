package common

import (
	"io"
	"log"
	"os"

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

type elasticLogger struct {
	*log.Logger
}

func (l *elasticLogger) Debugf(format string, args ...interface{}) {
	l.Printf("[debug] "+format, args...)
}

func (l *elasticLogger) Errorf(format string, args ...interface{}) {
	l.Printf("[error] "+format, args...)
}

func ElasticTracer() (opentracing.Tracer, io.Closer) {
	elasticTracer := apm.DefaultTracer
	logger := log.New(os.Stderr, "elastic ", log.Ldate|log.Ltime|log.Lshortfile)
	apm.DefaultTracer.SetLogger(&elasticLogger{Logger: logger})
	tracer := apmot.New(apmot.WithTracer(elasticTracer))
	return tracer, flushCloser{flusher: elasticTracer}
}
