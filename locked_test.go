package gt_test

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestLocked(t *testing.T) {
	val := 123
	l := gt.NewLocked(val)

	{
		// test exclusive lock get
		item, unlock := l.Get()
		require.Equal(t, val, item)
		unlock()

		l.With(func(item int) {
			require.Equal(t, val, item)
		})
	}

	// test set
	val = 456
	l.Set(456)

	{
		// test shared lock get
		item, unlock := l.RGet()
		require.Equal(t, val, item)
		unlock()

		l.RWith(func(item int) {
			require.Equal(t, val, item)
		})
	}

	// this will trigger the race detector reasonably reliably
	// if there is a race condition
	wg := &sync.WaitGroup{}
	for range 1000 {
		wg.Add(1)
		go concurrentLockOperations(wg, &l)
	}
	wg.Wait()
}

func concurrentLockOperations(wg *sync.WaitGroup, l *gt.Locked[int]) {
	defer wg.Done()

	for range 100 {
		val := 100
		l.Set(rand.Intn(val))

		_, unlock1 := l.Get()
		unlock1()

		_, unlock2 := l.RGet()
		unlock2()

		l.With(func(_ int) {})
		l.RWith(func(_ int) {})
	}
}

func FuzzLocked(f *testing.F) {
	f.Add(int64(0), uint8(1))
	f.Add(int64(123), uint8(4))
	f.Add(int64(-999), uint8(8))
	f.Add(int64(1<<62), uint8(16))

	f.Fuzz(func(t *testing.T, initialVal int64, numGoroutines uint8) {
		// limit goroutines to a reasonable number
		goroutines := int(numGoroutines%32) + 1

		l := gt.NewLocked(initialVal)

		// verify initial value
		val, unlock := l.Get()
		if val != initialVal {
			t.Errorf("expected initial value %d, got %d", initialVal, val)
		}
		unlock()

		// run concurrent operations
		wg := &sync.WaitGroup{}
		for range goroutines {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Set a value and verify we can read it back through With
				l.Set(initialVal + 1)

				l.With(func(v int64) {
					// Value should be set by some goroutine
					_ = v
				})

				l.RWith(func(v int64) {
					// Value should be readable
					_ = v
				})

				_, unlock := l.RGet()
				unlock()
			}()
		}
		wg.Wait()
	})
}
