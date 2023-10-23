TEST_PATH = $(shell go list ./... | grep -v internal)

.DEFAULT_GOAL := help

.PHONY: build
build: ## Build the CLI
	go build -o gvs cmd/gvs/main.go

.PHONY: install-deps
install-deps: ## Installs the dependencies
	go mod download && go get -v -t -d ./...

.PHONY: run
run: ## Runs the CLI. You can also pass flags in run e.g. make run FLAGS="--install-latest"
	go run cmd/gvs/main.go $(FLAGS)

.PHONY: format
format: ## Validates the files' format
	gofmt -d -s .

.PHONY: vet
vet: ## Examines Go source code and reports suspicious constructs
	go vet -tests=false ./...

.PHONY: lint
lint: ## Runs linter over the codebase. The rules can be found in ./.golangci.yml
	golangci-lint run

.PHONY: test
test: ## Runs the tests
	go test $(TEST_PATH)

.PHONY: test-file
test-file: ## Runs the tests for a specific file e.g. make test-file FILE=./version_test.go
	go test -v $(FILE)

.PHONY: test-coverage
test-coverage: ## Returns the coverage for each package
	go test -cover $(TEST_PATH)

.PHONY: test-coverage-list
test-coverage-list: ## Returns the extended coverage for each function and method per package
	go test -v -coverpkg=./... -coverprofile=profile.cov $(TEST_PATH) && go tool cover -func profile.cov

.PHONY: docs
docs: ## Extracts and generates documentation for Go. Once the server started, you can visit http://localhost:6060/
	godoc -http=:6060

.PHONY: help
help: ## parse jobs and descriptions from this Makefile
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' $(MAKEFILE_LIST) \
		| grep -Ev '^(help)\b[[:space:]]*:' \
		| awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%20s\033[0m \t%s\n", $$1, $$2}'