package gt_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	require := require.New(t)

	// create a pool and add a single byte
	p := gt.NewPool(func() *bytes.Buffer {
		return bytes.NewBuffer([]byte{})
	})

	p.Put(bytes.NewBuffer([]byte{'a'}))
	buf1 := p.Get()
	require.NotNil(buf1)

	// this will trigger the race detector reasonably reliably
	// if there is a race condition
	wg := &sync.WaitGroup{}
	for range 1000 {
		wg.Add(1)
		go concurrentPoolOperations(wg, p)
	}
	wg.Wait()
}

func concurrentPoolOperations(wg *sync.WaitGroup, p *gt.Pool[*bytes.Buffer]) {
	defer wg.Done()

	for ndx := range 100 {
		p.Put(bytes.NewBuffer([]byte{byte(ndx)}))
		_ = p.Get()
	}
}

func FuzzPool(f *testing.F) {
	f.Add([]byte("hello"), uint8(1))
	f.Add([]byte{}, uint8(4))
	f.Add([]byte{0, 1, 2, 3}, uint8(8))
	f.Add([]byte("test data for pool"), uint8(16))

	f.Fuzz(func(t *testing.T, data []byte, numGoroutines uint8) {
		// limit goroutines to a reasonable number
		goroutines := int(numGoroutines%32) + 1

		// create a pool that returns new byte slices
		p := gt.NewPool(func() []byte {
			return make([]byte, len(data))
		})

		// verify Get returns a valid slice from the New function
		got := p.Get()
		if got == nil {
			t.Error("Get should return non-nil slice")
		}
		if len(got) != len(data) {
			t.Errorf("expected slice of len %d, got %d", len(data), len(got))
		}

		// run concurrent Put/Get operations
		wg := &sync.WaitGroup{}
		for range goroutines {
			wg.Add(1)
			go func(input []byte) {
				defer wg.Done()

				for range 10 {
					// Put a copy of the input data
					buf := make([]byte, len(input))
					copy(buf, input)
					p.Put(buf)

					// Get should always return a non-nil slice
					result := p.Get()
					if result == nil {
						t.Error("Get returned nil")
					}
				}
			}(data)
		}
		wg.Wait()
	})
}
