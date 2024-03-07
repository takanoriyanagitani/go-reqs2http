package ch

import (
	"context"
	"errors"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
)

func Fold[T, U any](
	ctx context.Context,
	init U,
	ch <-chan T,
	reducer func(state U, next T) U,
) U {
	var state U = init
	for {
		select {
		case <-ctx.Done():
			return state
		case next, ok := <-ch:
			if !ok {
				return state
			}
			state = reducer(state, next)
		}
	}
}

func TryFold[T, U any](
	ctx context.Context,
	init U,
	ch <-chan pair.Pair[error, T],
	reducer func(state U, next T) pair.Pair[error, U],
) pair.Pair[error, U] {
	var state pair.Pair[error, U] = pair.Right[error](init)
	e := TryForEach(
		ctx,
		ch,
		func(t T) error {
			if nil != state.Left {
				return state.Left
			}
			state = reducer(state.Right, t)
			return nil
		},
	)
	return pair.Pair[error, U]{
		Left:  errors.Join(state.Left, e),
		Right: state.Right,
	}
}
