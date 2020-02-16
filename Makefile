# # Some interesting links on Makefiles:
# https://danishpraka.sh/2019/12/07/using-makefiles-for-go.html
# https://tech.davis-hansson.com/p/make/

MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
SHELL := bash

## build: install crystal dependencies and compile a release
build:
	@if ! [ -x "$$(command -v shards)" ]; then echo "Follow https://crystal-lang.org/install/ to install crystal"; exit 1; fi
	@shards build --release
.PHONY: build

## clean: clean the generated files and directories
clean:
	@rm -rf bin docs lib tmp
.PHONY: clean

## compile: compile with the full error trace
compile: lib
	@crystal build --error-trace ./src/cli.cr -o bin/cozy-desktop-ng
.PHONY: compile

## docs: build the documentation
docs: lib
	@crystal docs
.PHONY: docs

lib: shard.lock shard.yml
	@shards install

## lint: enforces a consistent code style and detect code smells
lint: lib
	@crystal tool format --check
	@bin/ameba
.PHONY: lint

## pretty: make the assets more prettier
pretty:
	@if ! [ -x "$$(command -v prettier)" ]; then echo "Install prettier with 'npm install -g prettier'"; exit 1; fi
	@prettier --write --no-semi src/web/public/*.{css,js}
	@prettier --write --parser html src/web/views/*
.PHONY: pretty

## update-deps: update the shards dependencies
update-deps:
	@shards update
.PHONY: update-deps

## tests: run the spec/tests
tests: lib
	@crystal spec spec/unit
.PHONY: tests

## web: start a web server for development
web: lib
	@crystal run src/cli.cr -- web
.PHONY: web

## help: print this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help
