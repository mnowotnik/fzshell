SHELL := bash

VERSION := $(shell git describe --tags --abbrev=0 2> /dev/null)
ifeq ($(VERSION),)
$(error Not on git repository; cannot determine $$VERSION)
endif

test:
	go test -v ./...

test-cover:
	go test -coverprofile=cover.out ./...

cover-report:
	go tool cover -html=cover.out

precommit: test
	grep -qF $(VERSION) scripts/install.sh

.PHONY: test test-cover cover-report
