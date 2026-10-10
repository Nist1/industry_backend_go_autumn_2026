package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(in) == 0 {
		return []Result[R]{}, nil
	}

	results := make([]Result[R], len(in))

	var next atomic.Int64

	n := min(workers, len(in))
	var wg sync.WaitGroup
	for range n {
		// wg.Go сам делает Add(1) перед запуском и Done() по завершении.
		wg.Go(func() {
			for {
				if ctx.Err() != nil {
					return
				}
				i := int(next.Add(1) - 1)
				if i >= len(in) {
					return
				}
				v, err := fn(ctx, in[i])
				if err != nil {
					results[i] = Result[R]{Err: err}
					continue
				}
				results[i] = Result[R]{Value: v}
			}
		})
	}

	// Ждём все горутины в любом случае, в том числе при отмене
	wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
