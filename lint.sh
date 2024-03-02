#!/bin/sh

verbose=${ENV_VERBOSE}

golangci-lint \
	run \
	"${verbose}"
