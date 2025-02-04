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
	require.NotNil(t, l)

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
	for ndx := 0; ndx < 1000; ndx++ {
		wg.Add(1)
		go concurrentLockOperations(wg, &l)
	}
	wg.Wait()
}

func concurrentLockOperations(wg *sync.WaitGroup, l *gt.Locked[int]) {
	defer wg.Done()

	for ndx := 0; ndx < 100; ndx++ {
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
