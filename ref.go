package gt

// Ref is an optional heap-allocated value, built on top of Option[*T].
// It exists to allow recursive type definitions that would otherwise
// create infinite-size structs with the value-based Option type.
//
// Use Ref for fields where types may reference themselves.
type Ref[T any] = Option[*T]

// SomeRef creates a Ref containing the given value.
func SomeRef[T any](item T) Ref[T] {
	return Some(&item)
}

// NoneRef creates an empty Ref.
func NoneRef[T any]() Ref[T] {
	return None[*T]()
}
