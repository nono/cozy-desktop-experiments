# # Some interesting links on Makefiles:
# https://danishpraka.sh/2019/12/07/using-makefiles-for-go.html
# https://tech.davis-hansson.com/p/make/

MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

## build: install crystal dependencies and compile a release
build:
	@shards build --release

## clean: clean the generated files and directories
clean:
	@rm -rf bin docs lib tmp

## compile: compile with the full error trace
compile:
	@crystal build --error-trace ./src/cli.cr -o bin/cozy-desktop-ng

## docs: build the documentation
docs:
	@crystal docs

## lint: enforces a consistent code style and detect code smells
lint:
	@crystal tool format --check
	@bin/ameba

## pretty: make the assets more prettier
pretty:
	@prettier --write --no-semi src/web/public/*.{css,js}
	@prettier --write --parser html src/web/views/*

## update-deps: update the shards dependencies
update-deps:
	@shards update

## tests: run the spec/tests
tests:
	@crystal spec spec/unit

## web: start a web server for development
web:
	@crystal run src/cli.cr -- web

## help: print this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: build check clean compile docs lint pretty update-deps tests web help
