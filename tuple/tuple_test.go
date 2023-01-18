package tuple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTupleIsPoint(t *testing.T) {
	p := New(4.3, -4.2, 3.1, 1.0)

	assert.Equal(t, 4.3, p.X)
	assert.Equal(t, -4.2, p.Y)
	assert.Equal(t, 3.1, p.Z)
	assert.Equal(t, 1.0, p.W)

	assert.True(t, p.IsPoint())
	assert.False(t, p.IsVector())
}

func TestTupleIsVector(t *testing.T) {
	v := New(4.3, -4.2, 3.1, 0.0)

	assert.Equal(t, 4.3, v.X)
	assert.Equal(t, -4.2, v.Y)
	assert.Equal(t, 3.1, v.Z)
	assert.Equal(t, 0.0, v.W)

	assert.False(t, v.IsPoint())
	assert.True(t, v.IsVector())
}

func TestCreatePoint(t *testing.T) {
	p := NewPoint(4, -4, 3)

	assert.Equal(t, 4.0, p.X)
	assert.Equal(t, -4.0, p.Y)
	assert.Equal(t, 3.0, p.Z)
	assert.Equal(t, 1.0, p.W)
}

func TestCreateVector(t *testing.T) {
	p := NewVector(4, -4, 3)

	assert.Equal(t, 4.0, p.X)
	assert.Equal(t, -4.0, p.Y)
	assert.Equal(t, 3.0, p.Z)
	assert.Equal(t, 0.0, p.W)
}

func TestEqualsExactly(t *testing.T) {
	p := NewPoint(1, 1, 1)
	p2 := NewPoint(1, 1, 1)

	assert.True(t, p.Equals(p2))
}

func TestNotEquals(t *testing.T) {
	p := NewPoint(1, 1, 1)
	p2 := NewPoint(1, 0, 1)

	assert.False(t, p.Equals(p2))
}

func TestEqualsWithinEpsilon(t *testing.T) {
	p := NewPoint(1, 1, 1)
	p2 := NewPoint(1.0000001, 1, 1)

	assert.True(t, p.Equals(p2))
}
