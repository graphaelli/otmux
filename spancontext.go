package otmux

import "github.com/opentracing/opentracing-go"

type spancontext struct {
	spancontexts []opentracing.SpanContext
}

func (sc *spancontext) ForeachBaggageItem(handler func(k, v string) bool) {
	for _, spancontext := range sc.spancontexts {
		spancontext.ForeachBaggageItem(handler)
	}
}
