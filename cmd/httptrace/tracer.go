package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type contextKey int

const (
	keyTracer contextKey = iota
)

type Tracer struct {
	tr      opentracing.Tracer
	root    opentracing.Span
	sp      opentracing.Span
	dns     opentracing.Span
	connect opentracing.Span
}

func TraceRequest(tr opentracing.Tracer, req *http.Request) (*http.Request, *Tracer) {
	ht := &Tracer{tr: tr}
	ctx := req.Context()
	req = req.WithContext(context.WithValue(ctx, keyTracer, ht))
	return req, ht
}

func TracerFromRequest(req *http.Request) *Tracer {
	tr, ok := req.Context().Value(keyTracer).(*Tracer)
	if !ok {
		return nil
	}
	return tr
}

type TracingTransport struct {
	http.RoundTripper
}

type closeTracker struct {
	io.ReadCloser
	sp opentracing.Span
}

func (c closeTracker) Close() error {
	err := c.ReadCloser.Close()
	c.sp.LogFields(log.String("event", "ClosedBody"))
	c.sp.Finish()
	return err
}

func (t *TracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.RoundTripper
	if rt == nil {
		rt = http.DefaultTransport
	}
	tracer := TracerFromRequest(req)
	if tracer == nil {
		panic("no tracer")
		return rt.RoundTrip(req)
	}

	tracer.start(req)

	ext.HTTPMethod.Set(tracer.sp, req.Method)
	ext.HTTPUrl.Set(tracer.sp, req.URL.String())

	resp, err := rt.RoundTrip(req)
	if err != nil {
		tracer.sp.Finish()
		return resp, err
	}
	ext.HTTPStatusCode.Set(tracer.sp, uint16(resp.StatusCode))
	if resp.StatusCode >= http.StatusInternalServerError {
		ext.Error.Set(tracer.sp, true)
	}
	if req.Method == "HEAD" {
		tracer.sp.Finish()
	} else {
		resp.Body = closeTracker{resp.Body, tracer.sp}
	}
	return resp, nil
}

func (h *Tracer) start(req *http.Request) opentracing.Span {
	if h.root == nil {
		root := h.tr.StartSpan("HTTP Client")
		h.root = root
	}

	ctx := h.root.Context()
	h.sp = h.tr.StartSpan("HTTP "+req.Method, opentracing.ChildOf(ctx))
	ext.SpanKindRPCClient.Set(h.sp)
	ext.Component.Set(h.sp, "net/http")
	return h.sp
}

func (h *Tracer) startDNS(t time.Time) {
	ctx := h.root.Context()
	h.dns = h.tr.StartSpan("DNS", opentracing.ChildOf(ctx), opentracing.StartTime(t))
}

func (h *Tracer) doneDNS(t time.Time) {
	h.dns.FinishWithOptions(opentracing.FinishOptions{FinishTime: t})
	h.dns = nil
}

func (h *Tracer) startConnect(t time.Time) {
	ctx := h.root.Context()
	h.connect = h.tr.StartSpan("Connect", opentracing.ChildOf(ctx), opentracing.StartTime(t))
}

func (h *Tracer) doneConnect(t time.Time) {
	h.connect.FinishWithOptions(opentracing.FinishOptions{FinishTime: t})
	h.connect = nil
}

// Finish finishes the span of the traced request.
func (h *Tracer) Finish() {
	if h.root != nil {
		h.root.Finish()
	}
}
