package gt_test

import (
	"errors"
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestResultTypes(t *testing.T) {
	{
		val := 123
		res := gt.ResultOK(val)
		require.NoError(t, res.Err())
		require.Equal(t, val, res.OK())

		switch r := res.Match().(type) {
		case error:
			require.FailNow(t, "match arm should not be an error")
		default:
			require.Equal(t, val, r)
		}
	}

	{
		err := errors.New("test error")
		res := gt.ResultErr[int](err)
		require.ErrorIs(t, res.Err(), err)

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
