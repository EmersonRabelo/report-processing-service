VERSION=$(shell cat VERSION)
GOFLAGS=-ldflags "-X main.Version=$(VERSION)"

build:
	@CGO_ENABLED=0 go build $(GOFLAGS) -o ./build/app ./cmd/report-processing-service

clean:
	@go clean ./..
	@rm -rf build/

help:
	@echo "Makefile commands:"
	@echo "  build    - Build the report processing service binary"
	@echo "  clean    - Clean up build artifacts"
	@echo "  help     - Show this help message"