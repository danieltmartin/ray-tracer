package world

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateID(t *testing.T) {
	w := New()
	numIDs := 10000
	ch := make(chan ID, numIDs)
	var wg sync.WaitGroup
	wg.Add(numIDs)

	for i := 0; i < numIDs; i++ {
		go func() {
			ch <- w.NextID()
			wg.Done()
		}()
	}

	wg.Wait()

	last := <-ch
	for i := 1; i < numIDs; i++ {
		next := <-ch
		t.Logf("checking %v and %v\n", last, next)
		require.NotEqual(t, last, next)
		last = next
	}
}
