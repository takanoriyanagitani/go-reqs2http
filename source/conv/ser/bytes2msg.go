package cser

type ConvertFn[M any] func(serialized []byte, buf M) error
