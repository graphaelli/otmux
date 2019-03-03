package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/opentracing/opentracing-go"

	"github.com/graphaelli/otmux"
	common "github.com/graphaelli/otmux/cmd/example-common"
)

func main() {
	logger := log.New(os.Stderr, "", common.LogFmt)

	elasticOpenTracer, elasticCloser := common.ElasticTracer()
	defer elasticCloser.Close()
	jaegerOpenTracer, jaegerCloser := common.JaegerTracer()
	defer jaegerCloser.Close()

	// Opentracing tracer
	tracer := otmux.NewTracer(elasticOpenTracer, jaegerOpenTracer)
	opentracing.SetGlobalTracer(tracer)

	addr, serverCloser := common.StartServer(tracer, logger)
	defer serverCloser.Close()
	fmt.Println("listening on", "http://"+addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
