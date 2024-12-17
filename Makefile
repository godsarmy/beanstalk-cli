GO ?= go

UNAME := $(shell uname)

# Set build targets based on OS
ifeq ($(UNAME), Linux)
	TARGET := linux
	GOOS := linux
	GOARCH := amd64
else ifeq ($(UNAME), Darwin)
	TARGET := darwin
	GOOS := darwin
	GOARCH := amd64
else ifeq ($(UNAME), Windows_NT)
	TARGET := windows
	GOOS := windows
	GOARCH := amd64
endif

VERSION ?= $(shell cat ./VERSION)

LDFLAGS_COMMON = \
	-X main.version=$(VERSION) 

GO_BUILD := $(GO) build $(EXTRA_FLAGS) -ldflags "$(LDFLAGS_COMMON)"

.DEFAULT: all

.PHONY: all
beanstalk-cli-$(TARGET): cmd/main.go cmd/functions.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) -o $@ $^

format:
	gofmt -w */*.go
	golines -w */*.go

.PHONY: all
clean:
	rm -rf beanstalk-cli-*

