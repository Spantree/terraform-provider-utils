.PHONY: build install test testacc fmt lint clean

# Default target
default: build

# Build the provider
build:
	go build -o terraform-provider-template

# Install the provider locally for development
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/spantree/template/1.0.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-template ~/.terraform.d/plugins/registry.terraform.io/spantree/template/1.0.0/$$(go env GOOS)_$$(go env GOARCH)/

# Run unit tests
test:
	go test -v ./...

# Run acceptance tests
testacc:
	TF_ACC=1 go test -v ./internal/provider/ -timeout 120m

# Format code
fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

# Lint code
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -f terraform-provider-template
	rm -rf dist/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Generate documentation
docs:
	go generate ./...

# Run all checks before committing
check: fmt test
	@echo "All checks passed!"

