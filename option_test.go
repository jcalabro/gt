package gt_test

import (
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestOptionType(t *testing.T) {
	{
		val := 123
		opt := gt.OptionSome(val)
		require.True(t, opt.HasValue())
		require.Equal(t, val, opt.Get())
	}

	{
		opt := gt.OptionNone[any]()
		require.False(t, opt.HasValue())

		recoverCalled := false
		defer func() { require.True(t, recoverCalled) }()
		defer func() {
			if r := recover(); r != nil {
				recoverCalled = true
			}
		}()
		require.Equal(t, 0, opt.Get())
	}
}
