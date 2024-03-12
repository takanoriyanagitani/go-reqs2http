#!/bin/sh

cargo \
	watch \
	--watch src \
	--watch Cargo.toml \
	--exec 'check --target wasm32-unknown-unknown --profile release-wasm'
