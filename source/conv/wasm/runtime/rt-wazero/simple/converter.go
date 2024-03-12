package sconv

import (
	"context"
	"errors"
	"fmt"

	wz "github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"

	util "github.com/takanoriyanagitani/go-reqs2http/util"

	wcnv "github.com/takanoriyanagitani/go-reqs2http/source/conv/wasm"
)

var (
	ErrUnableToWrite = errors.New("unable to write")
	ErrUnableToRead  = errors.New("unable to read")
	ErrNoReturnValue = errors.New("no return value")
)

type Converter struct {
	wa.Module

	ResizeI func(ctx context.Context, size int32) (capacity int32, e error)
	ResetO  func(ctx context.Context, size int32) (capacity int32, e error)

	OffsetI func(context.Context) (int32, error)
	OffsetO func(context.Context) (int32, error)

	Conv func(context.Context) (int32, error)
}

func (c Converter) Convert(
	ctx context.Context,
	input []byte,
) (output []byte, e error) {
	var sz int32 = int32(len(input))
	_, e = c.ResizeI(ctx, sz)
	if nil != e {
		return nil, e
	}

	_, e = c.ResetO(ctx, sz)
	if nil != e {
		return nil, e
	}

	var mem wa.Memory = c.Module.Memory()

	offi, ei := c.OffsetI(ctx)
	ei = errors.Join(ei, c.input2mem(mem, uint32(offi), input))

	osz, ec := c.Conv(ctx)

	offo, eo := c.OffsetO(ctx)
	output, ok := mem.Read(
		uint32(offo),
		uint32(osz),
	)
	if !ok {
		return nil, errors.Join(
			ei,
			eo,
			ec,
			ErrUnableToRead,
		)
	}

	return output, errors.Join(ei, eo, ec)
}

func (c Converter) Close(ctx context.Context) error {
	return c.Module.Close(ctx)
}

func (c Converter) input2mem(m wa.Memory, offset uint32, input []byte) error {
	return util.Select(
		ErrUnableToWrite,
		nil,
		nil != m && m.Write(offset, input),
	)
}

func (c Converter) AsIf() wcnv.Converter { return c }

type ConverterFactory struct {
	wa.Module

	ResizeI wa.Function
	ResetO  wa.Function

	OffsetI wa.Function
	OffsetO wa.Function

	Conv wa.Function
}

type UnaryFunc32i func(context.Context, int32) (int32, error)

func UnaryFunc32iNew(f wa.Function) UnaryFunc32i {
	return func(ctx context.Context, i int32) (o int32, e error) {
		ret, e := f.Call(ctx, wa.EncodeI32(i))
		if nil != e {
			return 0, e
		}
		if len(ret) < 1 {
			return 0, ErrNoReturnValue
		}
		return wa.DecodeI32(ret[0]), nil
	}
}

func UnitFn32iNew(f wa.Function) func(context.Context) (int32, error) {
	return func(ctx context.Context) (o int32, e error) {
		ret, e := f.Call(ctx)
		if nil != e {
			return 0, e
		}
		if len(ret) < 1 {
			return 0, ErrNoReturnValue
		}
		return wa.DecodeI32(ret[0]), nil
	}
}

func (f ConverterFactory) CreateConverter() (Converter, error) {
	var cnv Converter

	cnv.Module = f.Module

	cnv.ResizeI = UnaryFunc32iNew(f.ResizeI)
	cnv.ResetO = UnaryFunc32iNew(f.ResetO)

	cnv.OffsetI = UnitFn32iNew(f.OffsetI)
	cnv.OffsetO = UnitFn32iNew(f.OffsetO)
	cnv.Conv = UnitFn32iNew(f.Conv)

	return cnv, nil
}

type RawConverterFactory struct {
	wa.Module

	ResizeI string
	ResetO  string

	OffsetI string
	OffsetO string

	Conv string
}

func module2func(m wa.Module, name string) (wa.Function, error) {
	var f wa.Function = m.ExportedFunction(name)
	switch f {
	case nil:
		return nil, fmt.Errorf("no such function: %s", name)
	default:
		return f, nil
	}
}

func (r RawConverterFactory) ToFactory() (ConverterFactory, error) {
	ResizeI, eri := module2func(r.Module, r.ResizeI)
	ResetO, ero := module2func(r.Module, r.ResetO)
	OffsetI, eoi := module2func(r.Module, r.OffsetI)
	OffsetO, eoo := module2func(r.Module, r.OffsetO)
	Conv, ecn := module2func(r.Module, r.Conv)

	return ConverterFactory{
			Module:  r.Module,
			ResizeI: ResizeI,
			ResetO:  ResetO,
			OffsetI: OffsetI,
			OffsetO: OffsetO,
			Conv:    Conv,
		}, errors.Join(
			eri,
			ero,
			eoi,
			eoo,
			ecn,
		)
}

type Compiled struct {
	wz.CompiledModule
}
