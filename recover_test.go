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

	{
		// test panic+recover with a string preserves the panic message
		panicMsg := "specific panic message"
		err := func() (err error) {
			defer func() { err = gt.Recover(nil, recover()) }()
			panic(panicMsg)
		}()
		require.Error(t, err)
		require.Contains(t, err.Error(), panicMsg)
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

func FuzzRecover(f *testing.F) {
	f.Add("panic message")
	f.Add("")
	f.Add("error: something went wrong")
	f.Add("\x00\x01\x02")

	f.Fuzz(func(t *testing.T, panicVal string) {
		// test that Recover always returns an error when given a non-nil panic value
		err := gt.Recover(nil, panicVal)
		if err == nil {
			t.Error("Recover should return non-nil error for non-nil panic value")
		}

		// test that the panic message is preserved in the error
		if !contains(err.Error(), panicVal) {
			t.Errorf("error should contain panic value: got %q, want to contain %q", err.Error(), panicVal)
		}

		// test with a pre-existing error - panic should take precedence
		preErr := fmt.Errorf("pre-existing error")
		err = gt.Recover(preErr, panicVal)
		if err == nil {
			t.Error("Recover should return non-nil error for non-nil panic value")
		}
	})
}

func contains(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
