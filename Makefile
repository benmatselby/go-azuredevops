.PHONY: explain
explain:
	### Welcome
	#
	# Makefile for go-vsts
	#

.PHONY: test
test:
	go test -v -race ./... -coverprofile=coverage.out

.PHONY: test-cov
test-cov: test
	go tool cover -html=coverage.out
