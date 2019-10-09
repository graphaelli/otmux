module github.com/graphaelli/otmux/cmd/example-server

require (
	github.com/graphaelli/otmux v0.0.0-00010101000000-000000000000
	github.com/graphaelli/otmux/cmd/example-common v0.0.0-00010101000000-000000000000
	github.com/opentracing/opentracing-go v1.1.0
)

replace github.com/graphaelli/otmux => ../..

replace github.com/graphaelli/otmux/cmd/example-common => ../example-common
