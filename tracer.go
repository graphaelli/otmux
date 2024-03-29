package otmux

import (
	"errors"

	"github.com/opentracing/opentracing-go"
)

type tracer struct {
	tracers []opentracing.Tracer
}

func NewTracer(tracers ...opentracing.Tracer) opentracing.Tracer {
	return &tracer{
		tracers: tracers,
	}
}

func (tr *tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	spans := make([]opentracing.Span, len(tr.tracers))
	for i, t := range tr.tracers {
		spanOpts := make([]opentracing.StartSpanOption, len(opts))
		for j, o := range opts {
			switch v := o.(type) {
			case opentracing.SpanReference:
				if sc, ok := v.ReferencedContext.(*spancontext); ok {
					v.ReferencedContext = sc.spancontexts[i]
					spanOpts[j] = v
				}
			default:
				spanOpts[j] = o
			}
		}
		spans[i] = t.StartSpan(operationName, spanOpts...)
	}
	return &span{
		tracer: tr,
		spans: spans,
	}
}

func (tr *tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	var err error
	sc, ok := sm.(*spancontext)
	if !ok {
		return errors.New("bad span context")
	}
	for i, t := range tr.tracers {
		er := t.Inject(sc.spancontexts[i], format, carrier)
		if er != nil {
			err = er
		}
	}
	return err
}

func (tr *tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	var err error
	spancontexts := make([]opentracing.SpanContext, len(tr.tracers))
	for i, t := range tr.tracers {
		sc, er := t.Extract(format, carrier)
		if er != nil {
			err = er
		}
		spancontexts[i] = sc
	}
	return &spancontext{
		spancontexts: spancontexts,
	}, err
}
