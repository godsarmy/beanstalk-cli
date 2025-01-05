GO ?= go

ifndef GOOS
  GOOS := $(shell uname | tr '[:upper:]' '[:lower:]')
endif

ifndef GOARCH
  MACHINE=$(shell uname -m)
  ifeq ($(MACHINE), x86_64)
	GOARCH := amd64
    GOEXT :=
  else ifeq ($(MACHINE), i386)
	GOARCH := i386
    GOEXT :=
  else ifeq ($(MACHINE), arm)
	GOARCH := arm64
    GOEXT := .exe
  endif
endif
# Set build targets based on OS
VERSION ?= $(shell cat ./VERSION)

LDFLAGS_COMMON = \
	-X main.version=$(VERSION) 

GO_BUILD := $(GO) build $(EXTRA_FLAGS) -ldflags "$(LDFLAGS_COMMON)"

.DEFAULT: all

.PHONY: all
build: beanstalk-cli
	mkdir -p bin
	cp -f $< bin/$<.$(GOARCH).$(GOOS)$(GOEXT)

beanstalk-cli: cmd/main.go cmd/functions.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) -o $@ $^

format:
	gofmt -w */*.go
	golines -w */*.go

.PHONY: all
clean:
	rm -rf beanstalk-cli

superclean: clean
	rm -rf bin/
