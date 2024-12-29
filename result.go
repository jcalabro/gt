package gt

// Carries a payload or an error, but never both (similar to Rust's Result type).
type Result[T any] struct {
	item T
	err  error
}

// Returns the value of the Result. Panics if an error is set.
func (res Result[T]) OK() T {
	if res.err != nil {
		panic("result error set, but Ok() called")
	}
	return res.item
}

// Returns the error of the Result, or nil if no error is set.
func (res Result[T]) Err() error {
	return res.err
}

// Returns either the error or the value.
func (res Result[T]) Match() any {
	if res.err != nil {
		return res.err
	}

	return res.item
}

// Sets the result payload to a successful result
func ResultOK[T any](item T) Result[T] {
	return Result[T]{item: item}
}

// Sets the result payload to an unsuccessful result
func ResultErr[T any](err error) Result[T] {
	return Result[T]{err: err}
}
