module github.com/graphaelli/otmux/cmd/httptrace

go 1.13

require (
	github.com/fatih/color v1.7.0
	github.com/graphaelli/otmux v0.0.0-00010101000000-000000000000
	github.com/graphaelli/otmux/cmd/example-common v0.0.0-00010101000000-000000000000
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	golang.org/x/net v0.0.0-20191119073136-fc4aabc6c914
)

replace github.com/graphaelli/otmux => ../..

replace github.com/graphaelli/otmux/cmd/example-common => ../example-common
