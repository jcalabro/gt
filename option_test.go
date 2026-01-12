package gt_test

import (
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestOptionType(t *testing.T) {
	{
		val := 123
		opt := gt.Some(val)
		require.True(t, opt.HasVal())
		require.False(t, opt.IsNone())
		require.Equal(t, val, opt.Val())
	}

	{
		opt := gt.None[any]()
		require.False(t, opt.HasVal())
		require.True(t, opt.IsNone())

		recoverCalled := false
		defer func() { require.True(t, recoverCalled) }()
		defer func() {
			if r := recover(); r != nil {
				recoverCalled = true
			}
		}()
		require.Equal(t, 0, opt.Val())
	}
}
