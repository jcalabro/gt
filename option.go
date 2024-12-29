package gt

// Carries either a value of type T or nothing
type Option[T any] struct {
	hasValue bool
	item     T
}

// Returns true if the Option type has a value set, false otherwise
func (o Option[T]) HasValue() bool {
	return o.hasValue
}

// Returns the value of the Option. Panics if the Option has no value.
func (o Option[T]) Get() T {
	if !o.hasValue {
		panic("option has no value set, but Get() was called")
	}
	return o.item
}

// Sets the Option with the given value
func OptionSome[T any](item T) Option[T] {
	return Option[T]{hasValue: true, item: item}
}

// Sets the Option with no value
func OptionNone[T any]() Option[T] {
	return Option[T]{hasValue: false}
}
