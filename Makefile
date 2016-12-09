PACKAGES = $(shell go list ./...)

.PHONY: build fmt test

export GO15VENDOREXPERIMENT=1

default: build

build: fmt
	go build -v -o bin/swan-resolver *.go

clean:
	rm -rf bin/*

fmt:
	go fmt ./src/...

test:
	go test -cover=true .

