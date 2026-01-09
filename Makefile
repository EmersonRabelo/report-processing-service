VERSION=$(shell cat VERSION)
GOFLAGS=-ldflags "-X main.Version=$(VERSION)"

build:
	@CGO_ENABLED=0 go build $(GOFLAGS) -o ./build/app ./cmd/report-processing-service

clean:
	@go clean ./..
	@rm -rf build/