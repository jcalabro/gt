package gt

import (
	"context"
	"errors"
	"sync"
)

// Same as ConcurrentN, using 50 worker goroutines.
func Concurrent[T, V any](ctx context.Context, items []T, fn func(T) (V, error)) ([]V, error) {
	return ConcurrentN(ctx, items, 50, fn)
}

// ConcurrentN executes fn for each item across up to the specified number of
// concurrent goroutines and accumulates the results. Output result ordering
// matches input ordering. If any invocation of fn returns an error, all errors
// are joined via [errors.Join] and returned. If the context is cancelled, no new
// goroutines are dispatched, but in-flight goroutines are allowed to finish
// before returning. Partial results from completed goroutines are discarded on
// error. If fn panics, the panic is recovered and surfaced as an error. Returns
// an error if workers is 0.
func ConcurrentN[T, V any](ctx context.Context, items []T, workers int, fn func(T) (V, error)) ([]V, error) {
	if workers <= 0 {
		return nil, errors.New("concurrent workers must be greater than 0")
	}

	sem := make(chan struct{}, workers)
	results := make([]V, len(items))
	errs := make([]error, len(items))

	var wg sync.WaitGroup
	for i, item := range items {
		select {
		case <-ctx.Done():
			wg.Wait()
			return nil, ctx.Err()
		case sem <- struct{}{}:
		}

		wg.Add(1)
		go func(idx int, it T) {
			defer func() {
				<-sem
				wg.Done()
			}()

			var err error
			defer func() {
				errs[idx] = Recover(err, recover())
			}()

			res, err := fn(it)
			if err != nil {
				return
			}
			results[idx] = res
		}(i, item)
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return results, nil
}
