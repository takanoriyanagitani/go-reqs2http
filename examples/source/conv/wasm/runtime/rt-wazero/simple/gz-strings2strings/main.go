package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"

	"strings"

	wz "github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
	ua "github.com/takanoriyanagitani/go-reqs2http/util/arr"
	uch "github.com/takanoriyanagitani/go-reqs2http/util/ch"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"

	cser "github.com/takanoriyanagitani/go-reqs2http/source/conv/ser"

	wcnv "github.com/takanoriyanagitani/go-reqs2http/source/conv/wasm"
	wser "github.com/takanoriyanagitani/go-reqs2http/source/conv/wasm/ser"

	//revive:disable:line-length-limit
	w0s "github.com/takanoriyanagitani/go-reqs2http/source/conv/wasm/runtime/rt-wazero/simple"
	//revive:enable:line-length-limit
)

func must[T any](t T, e error) T {
	if nil != e {
		panic(e)
	}
	return t
}

func mustNil(e error) {
	if nil != e {
		panic(e)
	}
}

func mustEq[T comparable](a, b T) {
	if a != b {
		fmt.Printf("left: %v\n", a)
		fmt.Printf("right: %v\n", b)
		panic("unexpected value got")
	}
}

var dfs fs.FS = os.DirFS(".")

const wasmPath = "rs_gstr2str.wasm"

var wasmBytes []byte = must(fs.ReadFile(dfs, wasmPath))

var rtime wz.Runtime = wz.NewRuntime(context.Background())
var compiled wz.CompiledModule = must(rtime.CompileModule(
	context.Background(),
	wasmBytes,
))

var instance wa.Module = must(rtime.InstantiateModule(
	context.Background(),
	compiled,
	wz.NewModuleConfig(),
))

var raw w0s.RawConverterFactory = w0s.RawConverterFactory{
	Module: instance,

	ResizeI: "resize_i",
	ResetO:  "reset_o",

	OffsetI: "offset_i",
	OffsetO: "offset_o",

	Conv: "convert",
}

var fact w0s.ConverterFactory = must(raw.ToFactory())
var conv w0s.Converter = must(fact.CreateConverter())
var cnv wcnv.Converter = conv.AsIf()

type Msg struct {
	strs []string
}

var cfn cser.ConvertFn[*Msg] = func(serialized []byte, buf *Msg) error {
	var s string = string(serialized)
	var splited []string = strings.Split(s, ",")
	buf.strs = splited
	return nil
}

const base string = "https://localhost/"

var m2r cser.Message2requests[*Msg] = func(
	ctx context.Context,
	msg *Msg,
	dst chan<- pair.Pair[error, *rhp.Request],
) error {
	var urls []string = msg.strs
	return ua.TryForEach(
		urls,
		func(u string) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			dst <- pair.Right[error](&rhp.Request{Url: base + u})
			return nil
		},
	)
}

var i2r wser.Input2requests[*Msg] = wser.Input2requests[*Msg]{
	In2out:    cnv,
	Bytes2msg: cfn,
	Msg2reqs:  m2r,
}

const inputGzipBytes64 string = "H4sIAAAAAAAEAzPUMdIx1jEBAGeLu+wHAAAA"

var inputGzipBytes []byte = must(
	base64.StdEncoding.DecodeString(inputGzipBytes64),
)

func main() {
	defer func() {
		mustNil(compiled.Close(context.Background()))
		mustNil(rtime.Close(context.Background()))
	}()

	reqs := make(chan pair.Pair[error, *rhp.Request])

	go func() {
		defer close(reqs)

		var buf Msg

		mustNil(i2r.Input2chan(
			context.Background(),
			inputGzipBytes,
			&buf,
			reqs,
		))
	}()

	var requests pair.Pair[error, []*rhp.Request] = uch.TryFold(
		context.Background(),
		nil,
		reqs,
		func(
			state []*rhp.Request,
			next *rhp.Request,
		) pair.Pair[error, []*rhp.Request] {
			return pair.Right[error](append(state, next))
		},
	)

	mustNil(requests.Left)
	var q []*rhp.Request = requests.Right
	mustEq(len(q), 4)
	for _, req := range q {
		var url string = req.GetUrl()
		fmt.Printf("url: %s\n", url)
	}
}
