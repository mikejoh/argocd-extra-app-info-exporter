APPNAME := $(notdir $(CURDIR))

GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")
GIT_SHA := $(shell git rev-parse --short HEAD)
VERSION := $(GIT_TAG)$(shell git diff --quiet || echo "-$(GIT_SHA)-dirty")

CMDPATH := ./cmd/$(APPNAME)
BUILDPATH := ./build

# Go parameters
GOVERSION=1.21.4
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

LDFLAGS := -ldflags '-s -w \
						-X=github.com/mikejoh/$(APPNAME)/internal/buildinfo.Version=$(VERSION) \
						-X=github.com/mikejoh/$(APPNAME)/internal/buildinfo.Name=$(APPNAME) \
						-X=github.com/mikejoh/$(APPNAME)/internal/buildinfo.GitSHA=$(GIT_SHA)'

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -v -o $(BUILDPATH)/$(APPNAME) $(CMDPATH)

test: 
	$(GOTEST) -v ./...

testcov:
	$(GOTEST) ./... -coverprofile=coverage.out

dep:
	$(GOCMD) mod download

vet:
	$(GOCMD) vet ./...

lint:
	golangci-lint run -v --timeout=15m ./...

clean: 
	$(GOCLEAN)
	rm -f $(BUILDPATH)/$(APPNAME)

run:
	$(GOBUILD) -v -o $(BUILDPATH)/$(APPNAME) $(CMDPATH)
	$(BUILDPATH)/$(APPNAME)

install:
	cp $(BUILDPATH)/$(APPNAME) ~/.local/bin

.PHONY: all build test testcov clean run install dep vet lint