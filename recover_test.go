package gt_test

import (
	"fmt"
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestRecover(t *testing.T) {
	testErr := fmt.Errorf("test error")

	{
		// test panic+recover with an Error
		err := func() (err error) {
			defer func() { err = gt.Recover(nil, recover()) }()
			panic(testErr)
		}()
		require.ErrorIs(t, err, testErr)
	}

	{
		// test panic+recover with a string
		err := func() (err error) {
			defer func() { err = gt.Recover(nil, recover()) }()
			panic("panic happened")
		}()
		require.Error(t, err)
		require.NotErrorIs(t, err, testErr)
	}

	myFunc := func(input error) (err error) {
		defer func() { err = gt.Recover(input, recover()) }()
		return nil
	}

	{
		// test no panic with an error returned
		err := myFunc(testErr)
		require.Error(t, err)
		require.ErrorIs(t, err, testErr)
	}

	{
		// test no panic with an error returned
		err := myFunc(nil)
		require.NoError(t, err)
		require.NotErrorIs(t, err, testErr)
	}
}
