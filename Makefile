GO ?= go

VERSION ?= $(shell cat ./VERSION)

LDFLAGS_COMMON = \
	-X main.version=$(VERSION) 
	
GO_BUILD := $(GO) build $(EXTRA_FLAGS) -ldflags "$(LDFLAGS_COMMON)"

.DEFAULT: all

.PHONY: all
beanstalk-cli: cmd/main.go
	$(GO_BUILD) -o $@ $<

format:
	gofmt -w */*.go
	golines -w */*.go

.PHONY: all
clean:
	rm -rf beanstalk-cli

