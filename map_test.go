package gt_test

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestSafeMap(t *testing.T) {
	key := 123
	val := 456

	m := gt.NewSafeMap[int, int]()
	require.Equal(t, 0, m.Size())
	require.False(t, m.Get(key).HasVal())

	{
		// insert an item
		m.Put(key, val)
		require.Equal(t, 1, m.Size())

		item := m.Get(key)
		require.True(t, item.HasVal())
		require.Equal(t, val, item.Val())
	}

	{
		// overwrite the item
		val2 := val + 1
		m.Put(key, val2)
		require.Equal(t, 1, m.Size())

		item := m.Get(key)
		require.True(t, item.HasVal())
		require.Equal(t, val2, item.Val())
	}

	// this will trigger the race detector reasonably reliably
	// if there is a race condition
	wg := &sync.WaitGroup{}
	for ndx := 0; ndx < 1000; ndx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for ndx := 0; ndx < 100; ndx++ {
				val := 100
				m.Put(rand.Intn(val), rand.Intn(val))
				_ = m.Get(rand.Intn(val))
				_ = m.Size()
			}
		}()
	}
	wg.Wait()
}
