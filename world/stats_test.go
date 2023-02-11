package world

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentStats(t *testing.T) {
	s := Stats{}
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() { s.eyeRayCount.inc(); wg.Done() }()
	}
	wg.Wait()

	assert.EqualValues(t, 10000, s.eyeRayCount.val())
}
