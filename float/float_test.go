package float

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExactlyEqual(t *testing.T) {
	assert.True(t, Equal(1.0, 1.0))
}

func TestAlmostEqual(t *testing.T) {
	assert.True(t, Equal(1.0, 1.000001))
}

func TestNotEqual(t *testing.T) {
	assert.False(t, Equal(1.0, 2.0))
}
