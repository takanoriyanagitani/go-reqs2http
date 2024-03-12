#!/bin/sh

cargo \
	check \
	--target wasm32-unknown-unknown \
	--profile release-wasm
