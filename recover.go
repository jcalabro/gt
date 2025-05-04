package gt

import (
	"fmt"
)

// Recover is a simple way to catch and return errors due to panics. It
// coalesces the given error `e` and the given recover `rec` down to a
// single error. It prefers the recover's error to the passed error, if
// a panic was caught.
func Recover(e error, rec any) error {
	if rec != nil {
		if err, ok := rec.(error); ok {
			return err
		} else {
			return fmt.Errorf("caught panic: %w", e)
		}
	}

	return e
}
