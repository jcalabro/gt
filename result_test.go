package gt_test

import (
	"errors"
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestResultType(t *testing.T) {
	{
		val := 123
		res := gt.OK(val)
		require.NoError(t, res.Err())
		require.Equal(t, val, res.OK())
		require.True(t, res.IsOK())
		require.False(t, res.IsErr())

		switch r := res.Match().(type) {
		case error:
			require.FailNow(t, "match arm should not be an error")
		default:
			require.Equal(t, val, r)
		}
	}

	{
		err := errors.New("test error")
		res := gt.Err[int](err)
		require.ErrorIs(t, res.Err(), err)
		require.False(t, res.IsOK())
		require.True(t, res.IsErr())

		switch r := res.Match().(type) {
		case error:
			require.ErrorIs(t, err, r)
		default:
			require.FailNow(t, "match arm should not be a value")
		}

		recoverCalled := false
		defer func() { require.True(t, recoverCalled) }()
		defer func() {
			if r := recover(); r != nil {
				recoverCalled = true
			}
		}()
		require.Equal(t, 0, res.OK())
	}
}

func TestResultUnwrap(t *testing.T) {
	{
		// Unwrap returns value and nil error for OK result
		val := 42
		res := gt.OK(val)
		v, err := res.Unwrap()
		require.Equal(t, val, v)
		require.NoError(t, err)
	}

	{
		// Unwrap returns zero value and error for Err result
		testErr := errors.New("test error")
		res := gt.Err[int](testErr)
		v, err := res.Unwrap()
		require.Equal(t, 0, v)
		require.ErrorIs(t, err, testErr)
	}
}
