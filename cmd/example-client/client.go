package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/graphaelli/otmux"
	"github.com/graphaelli/otmux/cmd/example-common"
)

func main() {
	remote := flag.String("remote", "", "remote URL")
	flag.Parse()

	logger := log.New(os.Stderr, "", common.LogFmt)

	elasticOpenTracer, elasticCloser := common.ElasticTracer()
	defer elasticCloser.Close()
	jaegerOpenTracer, jaegerCloser := common.JaegerTracer()
	defer jaegerCloser.Close()
	zipkinOpenTracer, zipkinCloser := common.ZipkinTracer()
	defer zipkinCloser.Close()
	haystackOpenTracer, haystackCloser := common.HaystackTracer()

	// give haystack time to flush
	defer func() {
		// give haystack time to flush before closing, otherwise: Fail to dispatch to haystack-agent with error rpc error: code = Canceled desc = grpc: the client connection is closing
		time.Sleep(1 * time.Second)
		haystackCloser.Close()
	}()
	// Opentracing tracer
	tracer := otmux.NewTracer(elasticOpenTracer, jaegerOpenTracer, zipkinOpenTracer, haystackOpenTracer)
	opentracing.SetGlobalTracer(tracer)

	// Start an HTTP server
	if *remote == "" {
		addr, serverCloser := common.StartServer(tracer, logger)
		defer serverCloser.Close()
		*remote = "http://" + addr + "/"
	}

	// Make a request
	span := opentracing.StartSpan("root")
	defer span.Finish()
	time.Sleep(10 * time.Millisecond)
	child := opentracing.StartSpan("request", opentracing.ChildOf(span.Context()))
	req, _ := http.NewRequest(http.MethodGet, *remote, nil)
	/*
		if err := tracer.Inject(
			child.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
			logger.Printf("failed to inject: %s", err)
		}
	*/
	// do something concurrently with the request
	done := make(chan struct{})
	go func() {
		sibling := opentracing.StartSpan("background", opentracing.ChildOf(span.Context()))
		defer sibling.Finish()
		time.Sleep(13 * time.Millisecond)
		done <- struct{}{}
	}()

	if rsp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req); err != nil {
		child.LogKV("event", "error")
		child.LogKV("error", err.Error())
	} else {
		child.LogKV("status_code", rsp.StatusCode)
		rsp.Body.Close()
	}
	child.Finish()
	<-done

	// Do something afterwards
	post := opentracing.StartSpan("post", opentracing.ChildOf(span.Context()))
	defer post.Finish()
	time.Sleep(7 * time.Millisecond)
}
