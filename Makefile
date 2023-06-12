GO ?= go
GO_SRC = $(shell find . -name '*.go')

gee: $(GO_SRC)
	$(GO) build -o $@

test: gee
	$(GO) test ./...

clean:
	@rm -f gee

.PHONY: test clean
