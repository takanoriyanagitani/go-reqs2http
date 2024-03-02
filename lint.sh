#!/bin/sh

verbose=${ENV_VERBOSE}

gci() {
	golangci-lint \
		run \
		"${verbose}"
}

scheck() {
	staticcheck
}

gci
scheck
