package common

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

const LogFmt = log.Ldate | log.Ltime | log.Lshortfile

func StartServer(tracer opentracing.Tracer, logger *log.Logger) (string, io.Closer) {
	lis, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		panic(err)
	}
	addr := lis.Addr().(*net.TCPAddr).String()
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var opts []opentracing.StartSpanOption
			carrier := opentracing.HTTPHeadersCarrier(r.Header)
			if dtContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier); err == nil {
				opts = append(opts, opentracing.ChildOf(dtContext))
			}
			span := opentracing.StartSpan(fmt.Sprintf("%s %s", r.Method, r.URL.String()), opts...)
			for h, v := range r.Header {
				span.LogFields(otlog.String("header="+h, strings.Join(v, ",")))
			}
			logger.Println(r.URL.String(), r.Header)
			time.Sleep(6 * time.Millisecond)
			defer span.Finish()
		}),
	}
	go func() {
		server.Serve(lis)
	}()

	return addr, server
}

func NewTracers() ([]opentracing.Tracer, []io.Closer) {
	elasticOpenTracer, elasticCloser := ElasticTracer()
	jaegerOpenTracer, jaegerCloser := JaegerTracer()
	zipkinOpenTracer, zipkinCloser := ZipkinTracer()
	haystackOpenTracer, haystackCloser := HaystackTracer()

	return []opentracing.Tracer{
		elasticOpenTracer,
		haystackOpenTracer,
		jaegerOpenTracer,
		zipkinOpenTracer,
	}, []io.Closer{
		elasticCloser,
		haystackCloser,
		jaegerCloser,
		zipkinCloser,
	}
}
