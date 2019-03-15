package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/opentracing/opentracing-go"

	"github.com/graphaelli/otmux"
	"github.com/graphaelli/otmux/cmd/example-common"
)

func main() {
	logger := log.New(os.Stderr, "", common.LogFmt)

	tracers, closers := common.NewTracers()
	for _, c := range closers {
		// ok
		defer c.Close()
	}
	// Opentracing tracer
	tracer := otmux.NewTracer(tracers...)
	opentracing.SetGlobalTracer(tracer)

	addr, serverCloser := common.StartServer(tracer, logger)
	defer serverCloser.Close()
	fmt.Println("listening on", "http://"+addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
