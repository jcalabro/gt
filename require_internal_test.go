package gt

import (
	"reflect"
	"testing"
)

// require provides test assertion functions for the internal (package gt) test files.
var require = struct {
	Equal func(t testing.TB, expected, actual any)
}{
	Equal: func(t testing.TB, expected, actual any) {
		t.Helper()
		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("expected %v (%T), got %v (%T)", expected, expected, actual, actual)
		}
	},
}
