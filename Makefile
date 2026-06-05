GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOFMT = $(GOCMD) fmt
GOVET = $(GOCMD) vet
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

build: darwin

all: darwin linux

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/v4poolkb.darwin ./cmd

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/v4poolkb.linux ./cmd

clean:
	$(GOCLEAN) ./cmd/...
	rm -rf bin

fmt:
	$(GOFMT) ./...

vet:
	$(GOVET) ./...

dep:
	$(GOGET) -u ./cmd/...
	$(GOMOD) tidy
	$(GOMOD) verify

check:
	go install honnef.co/go/tools/cmd/staticcheck
	$(HOME)/go/bin/staticcheck -checks all,-S1002,-ST1003 ./cmd/...
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
	$(GOVET) -vettool=$(HOME)/go/bin/shadow ./cmd/...

test:
	$(GOTEST) ./...
