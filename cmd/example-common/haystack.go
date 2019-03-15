package common

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ExpediaDotCom/haystack-client-go"
	"github.com/opentracing/opentracing-go"
)

type consoleLogger struct{}

func (logger *consoleLogger) Error(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (logger *consoleLogger) Info(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (logger *consoleLogger) Debug(format string, v ...interface{}) {
	log.Printf(format, v...)
	log.Print("\n")
}

func HaystackTracer() (opentracing.Tracer, io.Closer) {
	tracer, closer := haystack.NewTracer(filepath.Base(os.Args[0]),
		haystack.NewAgentDispatcher("localhost", 34000, 3*time.Second, 1000),
		haystack.TracerOptionsFactory.Logger(&consoleLogger{}))
	return tracer, closer
}
