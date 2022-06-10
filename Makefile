SHELL := bash

test:
	go test -v ./...

test-cover:
	go test -coverprofile=cover.out ./...

cover-report:
	go tool cover -html=cover.out

.PHONY: test test-cover cover-report
