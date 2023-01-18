package tuple

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/float"
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

func TestAddPointAndVector(t *testing.T) {
	p := NewPoint(3, -2, 5)
	v := NewVector(-2, 3, 1)

	assert.Equal(t, p.Add(v), NewPoint(1, 1, 6))
}

func TestSubtractTwoPoints(t *testing.T) {
	p1 := NewPoint(3, 2, 1)
	p2 := NewPoint(5, 6, 7)

	assert.Equal(t, p1.Sub(p2), NewVector(-2, -4, -6))
}

func TestSubtractVectorFromPoint(t *testing.T) {
	p := NewPoint(3, 2, 1)
	v := NewVector(5, 6, 7)

	assert.Equal(t, p.Sub(v), NewPoint(-2, -4, -6))
}

func TestSubtractTwoVectors(t *testing.T) {
	v1 := NewVector(3, 2, 1)
	v2 := NewVector(5, 6, 7)

	assert.Equal(t, v1.Sub(v2), NewVector(-2, -4, -6))
}

func TestSubtractFromZeroVector(t *testing.T) {
	v1 := NewVector(0, 0, 0)
	v2 := NewVector(1, -2, 3)

	assert.Equal(t, v1.Sub(v2), NewVector(-1, 2, -3))
}

func TestNegate(t *testing.T) {
	tu := New(1, -2, 3, -4)

	assert.Equal(t, tu.Neg(), New(-1, 2, -3, 4))
}

func TestMultiplyByScalar(t *testing.T) {
	tu := New(1, -2, 3, -4)

	assert.Equal(t, tu.Mul(3.5), New(3.5, -7, 10.5, -14))
}

func TestMultiplyByFraction(t *testing.T) {
	tu := New(1, -2, 3, -4)

	assert.Equal(t, tu.Mul(0.5), New(0.5, -1, 1.5, -2))
}

func TestDivisionByScalar(t *testing.T) {
	tu := New(1, -2, 3, -4)

	assert.Equal(t, tu.Div(2), New(0.5, -1, 1.5, -2))
}

func TestMagnitude(t *testing.T) {
	assert.Equal(t, NewVector(1, 0, 0).Mag(), 1.0)
	assert.Equal(t, NewVector(0, 1, 0).Mag(), 1.0)
	assert.Equal(t, NewVector(0, 0, 1).Mag(), 1.0)
	assert.Equal(t, NewVector(1, 2, 3).Mag(), math.Sqrt(14))
	assert.Equal(t, NewVector(-1, -2, -3).Mag(), math.Sqrt(14))
}

func TestNormalize(t *testing.T) {
	assert.Equal(t, NewVector(4, 0, 0).Norm(), NewVector(1, 0, 0))
	assertAlmost(t, NewVector(1, 2, 3).Norm(), NewVector(0.26726, 0.53452, 0.80178))
}

func TestMagnitudeOfNormalizedVector(t *testing.T) {
	v := NewVector(1, 2, 3)

	assert.Equal(t, v.Norm().Mag(), 1.0)
}

func TestDotProduct(t *testing.T) {
	v1 := NewVector(1, 2, 3)
	v2 := NewVector(2, 3, 4)

	assert.Equal(t, v1.Dot(v2), 20.0)
}

func TestCrossProduct(t *testing.T) {
	v1 := NewVector(1, 2, 3)
	v2 := NewVector(2, 3, 4)

	assert.Equal(t, v1.Cross(v2), NewVector(-1, 2, -1))
	assert.Equal(t, v2.Cross(v1), NewVector(1, -2, 1))
}

func assertAlmost(t *testing.T, t1 Tuple, t2 Tuple) {
	assert.True(t, float.Equal(t1.X, t2.X), "X values differ: t1.X=%v, t2.X=%v", t1.X, t2.X)
	assert.True(t, float.Equal(t1.Y, t2.Y), "Y values differ: t1.Y=%v, t2.Y=%v", t1.Y, t2.Y)
	assert.True(t, float.Equal(t1.Z, t2.Z), "Z values differ: t1.Z=%v, t2.Z=%v", t1.Z, t2.Z)
	assert.True(t, float.Equal(t1.W, t2.W), "W values differ: t1.W=%v, t2.W=%v", t1.W, t2.W)
}
