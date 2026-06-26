#!/usr/bin/env make
# e-backend-cli

LDFLAGS=-s -w

all: build

build:
	go build -ldflags "$(LDFLAGS)"

test:
	go test ./generator/module ./generator/app
