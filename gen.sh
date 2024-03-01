#!/bin/sh

which protoc | fgrep -q protoc || exec sh -c 'echo protoc missing.; exit 1'
which protoc-gen-go | fgrep -q protoc-gen-go || exec sh -c 'echo protoc-gen-go missing.; exit 1'

protodir="./reqs2http-proto"

find \
	"${protodir}" \
	-type f \
	-name '*.proto' |
	xargs \
		protoc \
		--proto_path=reqs2http-proto \
		--go_out=. \
		--go_opt=paths=source_relative
