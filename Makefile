BINARY  := yummycli
PREFIX  ?= /usr/local
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
DATE    := $(shell date +%Y-%m-%d)
LDFLAGS := -X github.com/yummysource/yummycli/internal/build.Version=$(VERSION) \
           -X github.com/yummysource/yummycli/internal/build.Date=$(DATE)

.PHONY: build install uninstall clean test

## build: compile the binary with version metadata
build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) .

## install: build and install to PREFIX/bin (default /usr/local/bin)
install: build
	install -d $(PREFIX)/bin
	install -m755 $(BINARY) $(PREFIX)/bin/$(BINARY)

## uninstall: remove the installed binary
uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)

## clean: remove the local build artifact
clean:
	rm -f $(BINARY)

## test: run all tests
test:
	go test ./...
