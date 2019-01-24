.DEFAULT_GOAL := explain
.PHONY: explain
explain:
	### Welcome
	#
	# Makefile for go-azuredevops
	#
	#  _____
	# |   __|___
	# |  |  | . |
	# |_____|___|
	#   |  _  |___ _ _ ___ ___
	#   |     |- _| | |  _| -_|
	#   |__|__|___|___|_|_|___|
	#    |    \ ___ _ _|     |___ ___
	#    |  |  | -_| | |  |  | . |_ -|
	#    |____/|___|\_/|_____|  _|___|
	#                        |_|
	#
	### Targets
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Clean the local environment
	rm -fr vendor

.PHONY: install
install: ## Install dependencies
	go get ./...

.PHONY: vet
vet: ## Vet the code
	go vet ./...

.PHONY: lint
lint: ## Lint the code
	golint -set_exit_status ./...

.PHONY: build
build: ## Build the application
	go build ./...

.PHONY: test
test: ## Run the unit tests
	go test -v -race ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: test-cov
test-cov: test ## Run the unit tests with coverage
	go tool cover -html=coverage.out

.PHONY: all
all: clean install lint vet build test ## Run all the tasks

.PHONY: doc
doc: ## Generate the documentation
	godoc -http=:6060
