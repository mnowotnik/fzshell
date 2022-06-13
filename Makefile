SHELL := bash

VERSION := $(shell git describe --tags --abbrev=0 2> /dev/null| sed 's/v//g' )
ifeq ($(VERSION),)
$(error Not on git repository; cannot determine $$VERSION)
endif

test:
	go test -v ./...

lint:
	go vet ./...
	staticcheck ./...

test-cover:
	go test -coverprofile=cover.out ./...

cover.out: test-cover

cover-report: cover.out
	go tool cover -html=cover.out

precommit: test
	grep -qF $(VERSION) scripts/install.sh

.PHONY: test test-cover cover-report
