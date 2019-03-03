package otmux

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type span struct {
	tracer opentracing.Tracer
	spans []opentracing.Span
}

func (sp *span) BaggageItem(restrictedKey string) string {
	var previous string
	for _, s := range sp.spans {
		if p := s.BaggageItem(restrictedKey); p != "" {
			previous = p
		}
	}
	return previous
}

func (sp *span) Context() opentracing.SpanContext {
	spancontexts := make([]opentracing.SpanContext, len(sp.spans))
	for i, s := range sp.spans {
		spancontexts[i] = s.Context()
	}
	return &spancontext{
		spancontexts: spancontexts,
	}
}

func (sp *span) Finish() {
	for _, s := range sp.spans {
		s.Finish()
	}
}

func (sp *span) FinishWithOptions(opts opentracing.FinishOptions) {
	for _, s := range sp.spans {
		s.FinishWithOptions(opts)
	}
}

func (sp *span) SetOperationName(operationName string) opentracing.Span {
	for _, s := range sp.spans {
		s.SetOperationName(operationName)
	}
	return sp
}

func (sp *span) SetTag(key string, value interface{}) opentracing.Span {
	for _, s := range sp.spans {
		s.SetTag(key, value)
	}
	return sp
}

func (sp *span) LogFields(fields ...log.Field) {
	for _, s := range sp.spans {
		s.LogFields(fields...)
	}

}

func (sp *span) LogKV(alternatingKeyValues ...interface{}) {
	for _, s := range sp.spans {
		s.LogKV(alternatingKeyValues...)
	}

}

func (sp *span) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	for _, s := range sp.spans {
		s.SetBaggageItem(restrictedKey, value)
	}
	return sp
}

func (sp *span) Tracer() opentracing.Tracer {
	return sp.tracer
}

func (sp *span) LogEvent(event string) {
	for _, s := range sp.spans {
		s.LogEvent(event)
	}

}

func (sp *span) LogEventWithPayload(event string, payload interface{}) {
	for _, s := range sp.spans {
		s.LogEventWithPayload(event, payload)
	}

}

func (sp *span) Log(data opentracing.LogData) {
	for _, s := range sp.spans {
		s.Log(data)
	}
}
