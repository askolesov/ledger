# Run unit tests
.PHONY: test
test:
	go test -v ./pkg/...

# Run integration tests
.PHONY: test-integration
test-integration: build
	go test -v ./test/integration/...

# Run all tests
.PHONY: test-all
test-all: test test-integration

# Build the binary
.PHONY: build
build:
	go build -o build/ -v ./...

# Install the binary to GOPATH/bin
.PHONY: install
install:
	go install -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf build/

# Clean integration test outputs
.PHONY: test-clean
test-clean:
	rm -f test/integration/testdata/migration/output.yaml
