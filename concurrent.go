package gt

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
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

	results := make([]V, len(items))
	errs := make([]error, len(items))

	var g errgroup.Group
	g.SetLimit(workers)

	for i, item := range items {
		if err := ctx.Err(); err != nil {
			// Stop dispatching new work on cancellation, but let anything
			// already running finish so we don't leak goroutines.
			_ = g.Wait()
			return nil, err
		}

		g.Go(func() (retErr error) {
			defer func() { errs[i] = Recover(retErr, recover()) }()

			res, err := fn(item)
			if err != nil {
				return err
			}

			results[i] = res
			return nil
		})
	}

	// g.Wait returns the first non-nil error from g.Go closures, but we want
	// every error joined. Swallow Wait's return and build the joined error
	// ourselves from the per-index slice.
	_ = g.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return results, nil
}
