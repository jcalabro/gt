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
