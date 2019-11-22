package common

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
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

type tm func() (opentracing.Tracer, io.Closer)

// only for use during init
var commonTracers = make(map[string]tm)

func RegisterTracer(name string, traceMaker tm) {
	commonTracers[name] = traceMaker
}

func NewTracers() ([]opentracing.Tracer, []io.Closer) {
	var tracers []opentracing.Tracer
	var closers []io.Closer

	disabledTracers := make(map[string]bool)
	for _, d:= range strings.Split(os.Getenv("DISABLE_TRACERS"), ",") {
		disabledTracers[strings.ToLower(d)] = true
	}
	for name, maker := range commonTracers {
		if disabledTracers[strings.ToLower(name)] {
			continue
		}
		t, c := maker()
		tracers = append(tracers, t)
		closers = append(closers, c)
	}
	// not strictly necessary
	if len(tracers) == 0 {
		tracers = []opentracing.Tracer{opentracing.NoopTracer{}}
	}
	return tracers, closers
}
