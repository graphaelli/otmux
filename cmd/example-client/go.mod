module github.com/graphaelli/otmux/cmd/example-client

require (
	github.com/graphaelli/otmux v0.0.0-00010101000000-000000000000
	github.com/graphaelli/otmux/cmd/example-common v0.0.0-00010101000000-000000000000
	github.com/opentracing/opentracing-go v1.0.2
)

replace github.com/graphaelli/otmux => ../..

replace github.com/graphaelli/otmux/cmd/example-common => ../example-common