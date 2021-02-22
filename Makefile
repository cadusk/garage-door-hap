.PHONY: build run clean help setup vet

PROJECTNAME=garagedoor

# Go related variables.
GOBASE=$(shell pwd)
GOFILES=$(wildcard *.go)

# makes it easy so we don't need to prepend every command with @
MAKEFLAGS += --silent

## build: Compiles binary for the current platform
build:
        echo '> Building the application ...'
        go build -o bin/$(PROJECTNAME) $(GOFILES)

release:
        echo '> Building release for ARM6 ...'
        GOOS=linux GOARCH=arm GOARM=6 go build -o bin/$(PROJECTNAME)-armv6 $(GOFILES)

## run: Runs the application
run:
        echo '> Running the application ...'
        go run -race $(GOFILES)

## test: Runs all local tests
test:
        echo '> Running all tests ...'
        go test -race ./...

## clean: Removes generated binaries
clean:
        echo '> Cleaning ...'
        rm -rf bin/

vet:
        echo '> Vetting ...'
        go vet ./...

## setup: setup go modules
setup:
        echo '> Setting up the environment ...'
        go mod tidy \
                && go mod vendor

## help: Prints this help message
help:
        echo "Usage: \n"
        sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
