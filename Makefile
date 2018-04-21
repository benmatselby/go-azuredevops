.PHONY: explain
explain:
	### Welcome
	#
	# Makefile for go-vsts
	#

.PHONY: clean
clean:
	rm -fr vendor

.PHONY: install
install:
	dep ensure

.PHONY: vet
vet:
	go vet -v ./...

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -v -race ./... -coverprofile=coverage.out

.PHONY: test-cov
test-cov: test
	go tool cover -html=coverage.out

.PHONY: all
all: clean install vet build test

.PHONY: doc
doc:
	godoc -http=:6060
