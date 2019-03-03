# Opentracing Multiplexer for Go

This is an experiment for duplicating all trace reporting through multiple opentracing tracers within a single go application.
It is not intended for production use.

## Quickstart

```go
import (
	"github.com/graphaelli/otmux"
)

tracer := otmux.NewTracer(elasticOpenTracer, jaegerOpenTracer)
opentracing.SetGlobalTracer(tracer)
```

## Example

[cmd/](cmd) contains an example client and server wired up to [Jaeger](http://jaegertracing.io) and [Elastic APM](https://www.elastic.co/solutions/apm).

![Jaeger Tracing and Elastic APM tracing the same activity](opentracing-jaeger-elastic.png)
