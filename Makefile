.PHONY: help lint fmt install-dev-dependencies test
# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

install-dev-dependencies: ## Install development dependencies
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest

lint: ## Run linter
	@for dir in */; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Linting $$dir..."; \
			(cd "$$dir" && golangci-lint run -c ../.golangci.yml ./...); \
		fi \
	done

fmt: ## Format code
	@for dir in */; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Formatting $$dir..."; \
			(cd "$$dir" && go fmt ./...); \
		fi \
	done

test: ## Run tests with formatted output and coverage
	@for dir in */; do \
		if [ -f "$$dir/go.mod" ]; then \
			(cd "$$dir" && go test $$(go list ./... | grep -v /examples/) -json -cover 2>&1 | gotestfmt -hide successful-tests,empty-packages); \
		fi \
	done
