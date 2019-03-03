# Opentracing Mulitiplexer for Go

This is an experiment for duplicating all trace reporting through multiple opentracing tracers within a single go application.
It is not intended for production use.

## Quickstart

```
import (
	"github.com/graphaelli/otmux"
)

tracer := otmux.NewTracer(elasticOpenTracer, jaegerOpenTracer)
opentracing.SetGlobalTracer(tracer)
```
