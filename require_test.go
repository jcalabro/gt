package gt_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// require provides test assertion functions that mirror testify/require's API
// since we don't want to import it and drag in its transitive dependencies.
var require = struct {
	Equal      func(t testing.TB, expected, actual any)
	NoError    func(t testing.TB, err error)
	Error      func(t testing.TB, err error)
	True       func(t testing.TB, val bool)
	False      func(t testing.TB, val bool)
	ErrorIs    func(t testing.TB, err, target error)
	NotErrorIs func(t testing.TB, err, target error)
	Contains   func(t testing.TB, s, substr string)
	Panics     func(t testing.TB, fn func())
	NotNil     func(t testing.TB, val any)
	Nil        func(t testing.TB, val any)
	Empty      func(t testing.TB, val any)
	Len        func(t testing.TB, val any, length int)
	FailNow    func(t testing.TB, msg string)
	New        func(t testing.TB) *Require
}{
	Equal: func(t testing.TB, expected, actual any) {
		t.Helper()
		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("expected %v (%T), got %v (%T)", expected, expected, actual, actual)
		}
	},
	NoError: func(t testing.TB, err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	},
	Error: func(t testing.TB, err error) {
		t.Helper()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	},
	True: func(t testing.TB, val bool) {
		t.Helper()
		if !val {
			t.Fatal("expected true, got false")
		}
	},
	False: func(t testing.TB, val bool) {
		t.Helper()
		if val {
			t.Fatal("expected false, got true")
		}
	},
	ErrorIs: func(t testing.TB, err, target error) {
		t.Helper()
		if !errors.Is(err, target) {
			t.Fatalf("expected error %v to be %v", err, target)
		}
	},
	NotErrorIs: func(t testing.TB, err, target error) {
		t.Helper()
		if errors.Is(err, target) {
			t.Fatalf("expected error %v to not be %v", err, target)
		}
	},
	Contains: func(t testing.TB, s, substr string) {
		t.Helper()
		if !strings.Contains(s, substr) {
			t.Fatalf("expected %q to contain %q", s, substr)
		}
	},
	Panics: func(t testing.TB, fn func()) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, but did not panic")
			}
		}()
		fn()
	},
	NotNil: func(t testing.TB, val any) {
		t.Helper()
		if isNil(val) {
			t.Fatal("expected non-nil, got nil")
		}
	},
	Nil: func(t testing.TB, val any) {
		t.Helper()
		if !isNil(val) {
			t.Fatalf("expected nil, got %v", val)
		}
	},
	Empty: func(t testing.TB, val any) {
		t.Helper()
		v := reflect.ValueOf(val)
		if !v.IsValid() {
			return
		}
		if v.Len() != 0 {
			t.Fatalf("expected empty, got length %d", v.Len())
		}
	},
	Len: func(t testing.TB, val any, length int) {
		t.Helper()
		v := reflect.ValueOf(val)
		if v.Len() != length {
			t.Fatalf("expected length %d, got %d", length, v.Len())
		}
	},
	FailNow: func(t testing.TB, msg string) {
		t.Helper()
		t.Fatal(msg)
	},
	New: func(t testing.TB) *Require {
		return &Require{t: t}
	},
}

// Require is an assertion object for use with require.New(t).
type Require struct {
	t testing.TB
}

func (r *Require) Equal(expected, actual any) {
	r.t.Helper()
	require.Equal(r.t, expected, actual)
}

func (r *Require) NoError(err error) {
	r.t.Helper()
	require.NoError(r.t, err)
}

func (r *Require) Error(err error) {
	r.t.Helper()
	require.Error(r.t, err)
}

func (r *Require) True(val bool) {
	r.t.Helper()
	require.True(r.t, val)
}

func (r *Require) False(val bool) {
	r.t.Helper()
	require.False(r.t, val)
}

func (r *Require) NotNil(val any) {
	r.t.Helper()
	require.NotNil(r.t, val)
}

func (r *Require) Nil(val any) {
	r.t.Helper()
	require.Nil(r.t, val)
}

func (r *Require) Contains(s, substr string) {
	r.t.Helper()
	require.Contains(r.t, s, substr)
}

func isNil(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	}

	return false
}
