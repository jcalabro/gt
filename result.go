package gt

// Carries a payload or an error, but never both (similar to Rust's Result type).
type Result[T any] struct {
	err  error
	item T
}

// Returns the value of the Result. Panics if an error is set.
func (r Result[T]) OK() T {
	if r.err != nil {
		panic("result error set, but Ok() was called")
	}
	return r.item
}

// Returns the error of the Result, or nil if no error is set.
func (r Result[T]) Err() error {
	return r.err
}

// Returns either the error or the value.
func (r Result[T]) Match() any {
	if r.err != nil {
		return r.err
	}

	return r.item
}

// Sets the result payload to a successful result
func ResultOK[T any](item T) Result[T] {
	return Result[T]{item: item}
}

// Sets the result payload to an unsuccessful result
func ResultErr[T any](err error) Result[T] {
	return Result[T]{err: err}
}
