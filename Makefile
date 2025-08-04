# Run tests
.PHONY: test
test:
	go test -v ./...

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