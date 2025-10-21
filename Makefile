.PHONY: build test clean run-challenge run-verifier run-gateway run-client help

help: ## Display this help message
	@echo "Proof of Work - Makefile commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build all services
	bazel build //...

test: ## Run all tests
	bazel test //...

clean: ## Clean build artifacts
	bazel clean

run-challenge: ## Run the challenge service
	bazel run //services/challenge

run-verifier: ## Run the verifier service
	bazel run //services/verifier

run-gateway: ## Run the gateway service
	bazel run //services/gateway

run-client: ## Run the example client
	bazel run //examples/client

update-deps: ## Update Go dependencies
	bazel run //:gazelle-update-repos

format: ## Format and update BUILD files
	bazel run //:gazelle
