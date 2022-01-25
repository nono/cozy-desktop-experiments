# # Some interesting links on Makefiles:
# https://danishpraka.sh/2019/12/07/using-makefiles-for-go.html
# https://tech.davis-hansson.com/p/make/

MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
SHELL := bash

## install: compile the code and installs in binary in $GOPATH/bin
install: web/assets/cozy-bs.min.css
	@go install
.PHONY: install

## run: start the client for local development
run: web/assets/cozy-bs.min.css
	@go run .
.PHONY: run

web/assets/cozy-bs.min.css:
	@curl https://unpkg.com/cozy-bootstrap@1.11.3/dist/cozy-bs.min.css -o $@

## lint: enforce a consistent code style and detect code smells
lint: scripts/golangci-lint
	@scripts/golangci-lint run -E gofmt -E unconvert -E misspell -E whitespace -E exportloopref -E bodyclose -E exhaustive -E nilnil -E bidichk --max-same-issues 10
.PHONY: lint

scripts/golangci-lint: Makefile
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ./scripts v1.44.0

## unit-tests: run the tests
unit-tests:
	@go test -timeout 2m -short ./...
.PHONY: unit-tests

## integration-tests: run the tests
integration-tests:
	@go test ./client -rapid.checks=10 -rapid.steps=10 -rapid.shrinktime=5s -rapid.v -rapid.nofailfile
	@go test ./localfs -rapid.checks=10 -rapid.steps=100 -rapid.shrinktime=5s -rapid.v -rapid.nofailfile
.PHONY: integration-tests

## clean: clean the generated files and directories
clean:
	@go clean
.PHONY: clean

## help: print this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help
