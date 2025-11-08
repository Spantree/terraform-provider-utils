.PHONY: build install test testacc fmt lint clean

# Default target
default: build

# Build the provider
build:
	go build -o terraform-provider-utils

# Install the provider locally for development
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/spantree/utils/1.0.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-utils ~/.terraform.d/plugins/registry.terraform.io/spantree/utils/1.0.0/$$(go env GOOS)_$$(go env GOARCH)/

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
	rm -f terraform-provider-utils
	rm -rf dist/
	rm -rf docs/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Generate Terraform Registry documentation (in docs/ directory)
docs:
	go generate ./...

# Run all checks before committing
check: fmt lint test
	@echo "All checks passed!"

# Help target
help:
	@echo "Available targets:"
	@echo "  build      - Build the provider binary"
	@echo "  install    - Install provider locally for development"
	@echo "  test       - Run unit tests"
	@echo "  testacc    - Run acceptance tests (requires TF_ACC=1)"
	@echo "  fmt        - Format code and examples"
	@echo "  lint       - Run golangci-lint"
	@echo "  docs       - Generate Terraform Registry documentation"
	@echo "  clean      - Remove build artifacts and generated docs"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  check      - Run fmt, lint, and test"
