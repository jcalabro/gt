package gt

import "encoding/json"

// Carries either a value of type T or nothing
type Option[T any] struct {
	hasValue bool
	item     T
}

// Returns true if the Option type has a value set, false otherwise
func (o Option[T]) HasVal() bool {
	return o.hasValue
}

// Returns the value of the Option. Panics if the Option has no value
func (o Option[T]) Val() T {
	if !o.hasValue {
		panic("option has no value set, but Val() was called")
	}
	return o.item
}

// Sets the Option with the given value
func Some[T any](item T) Option[T] {
	return Option[T]{hasValue: true, item: item}
}

// Sets the Option with no value
func None[T any]() Option[T] {
	return Option[T]{hasValue: false}
}

// Returns true if the Option has no value set
func (o Option[T]) IsNone() bool {
	return !o.hasValue
}

// MarshalJSON encodes the Option as JSON. None marshals as null,
// Some(v) marshals as the JSON encoding of v.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.HasVal() {
		return json.Marshal(o.item)
	}
	return []byte("null"), nil
}

// UnmarshalJSON decodes JSON into the Option. null produces None,
// any other value produces Some(v).
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	var item T
	if string(data) == "null" {
		o.hasValue = false
		o.item = item
		return nil
	}

	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}

	o.hasValue = true
	o.item = item
	return nil
}
