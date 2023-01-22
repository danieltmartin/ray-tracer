package transform

import (
	"math"
	"testing"

	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func TestMultiplyByTranslationMatrix(t *testing.T) {
	transform := Translation(5, -3, 2)
	point := tuple.NewPoint(-3, 4, 5)

	assert.Equal(t, tuple.NewPoint(2, 1, 7), transform.MulTuple(point))
}

func TestMultiplyByInverseOfTranslationMatrix(t *testing.T) {
	transform := Translation(5, -3, 2)
	inverse := transform.Inverse()
	point := tuple.NewPoint(-3, 4, 5)

	assert.Equal(t, tuple.NewPoint(-8, 7, 3), inverse.MulTuple(point))
}

func TestTranslationDoesNotAffectVectors(t *testing.T) {
	transform := Translation(5, -3, 2)
	vector := tuple.NewVector(-3, 4, 5)

	assert.Equal(t, vector, transform.MulTuple(vector))
}

func TestScalePoint(t *testing.T) {
	transform := Scaling(2, 3, 4)
	point := tuple.NewPoint(-4, 6, 8)

	assert.Equal(t, tuple.NewPoint(-8, 18, 32), transform.MulTuple(point))
}

func TestScaleVector(t *testing.T) {
	transform := Scaling(2, 3, 4)
	vector := tuple.NewVector(-4, 6, 8)

	assert.Equal(t, tuple.NewVector(-8, 18, 32), transform.MulTuple(vector))
}

func TestScaleByInverse(t *testing.T) {
	transform := Scaling(2, 3, 4)
	inverse := transform.Inverse()
	vector := tuple.NewVector(-4, 6, 8)

	assert.Equal(t, tuple.NewVector(-2, 2, 2), inverse.MulTuple(vector))
}

func TestReflection(t *testing.T) {
	transform := Scaling(-1, 1, 1)
	point := tuple.NewPoint(2, 3, 4)

	assert.Equal(t, tuple.NewPoint(-2, 3, 4), transform.MulTuple(point))
}

func TestRotationAroundXAxis(t *testing.T) {
	point := tuple.NewPoint(0, 1, 0)
	halfQuarter := RotationX(math.Pi / 4)
	fullQuarter := RotationX(math.Pi / 2)

	expected := tuple.NewPoint(0, math.Sqrt2/2, math.Sqrt2/2)
	actual := halfQuarter.MulTuple(point)
	assert.True(t, expected.Equals(actual))

	expected = tuple.NewPoint(0, 0, 1)
	actual = fullQuarter.MulTuple(point)
	assert.True(t, expected.Equals(actual))
}

func TestInverseOfRotationRotatesInOppositeDirection(t *testing.T) {
	point := tuple.NewPoint(0, 1, 0)
	halfQuarter := RotationX(math.Pi / 4)

	expected := tuple.NewPoint(0, math.Sqrt2/2, -math.Sqrt2/2)
	actual := halfQuarter.Inverse().MulTuple(point)

	assert.True(t, expected.Equals(actual))
}

func TestRotationAroundYAxis(t *testing.T) {
	point := tuple.NewPoint(0, 0, 1)
	halfQuarter := RotationY(math.Pi / 4)
	fullQuarter := RotationY(math.Pi / 2)

	expected := tuple.NewPoint(math.Sqrt2/2, 0, math.Sqrt2/2)
	actual := halfQuarter.MulTuple(point)
	assert.True(t, expected.Equals(actual))

	expected = tuple.NewPoint(1, 0, 0)
	actual = fullQuarter.MulTuple(point)
	assert.True(t, expected.Equals(actual))
}

func TestRotationAroundZAxis(t *testing.T) {
	point := tuple.NewPoint(0, 1, 0)
	halfQuarter := RotationZ(math.Pi / 4)
	fullQuarter := RotationZ(math.Pi / 2)

	expected := tuple.NewPoint(-math.Sqrt2/2, math.Sqrt2/2, 0)
	actual := halfQuarter.MulTuple(point)
	assert.True(t, expected.Equals(actual))

	expected = tuple.NewPoint(-1, 0, 0)
	actual = fullQuarter.MulTuple(point)
	assert.True(t, expected.Equals(actual))
}

func TestShearingXInProportionToY(t *testing.T) {
	transform := Shearing(1, 0, 0, 0, 0, 0)
	point := tuple.NewPoint(2, 3, 4)

	assert.Equal(t, tuple.NewPoint(5, 3, 4), transform.MulTuple(point))
}
func TestShearingXInProportionToZ(t *testing.T) {
	transform := Shearing(0, 1, 0, 0, 0, 0)
	point := tuple.NewPoint(2, 3, 4)

	assert.Equal(t, tuple.NewPoint(6, 3, 4), transform.MulTuple(point))
}

func TestShearingYInProportionToX(t *testing.T) {
	transform := Shearing(0, 0, 1, 0, 0, 0)
	point := tuple.NewPoint(2, 3, 4)

	assert.Equal(t, tuple.NewPoint(2, 5, 4), transform.MulTuple(point))
}

func TestShearingYInProportionToZ(t *testing.T) {
	transform := Shearing(0, 0, 0, 1, 0, 0)
	point := tuple.NewPoint(2, 3, 4)

	assert.Equal(t, tuple.NewPoint(2, 7, 4), transform.MulTuple(point))
}

func TestShearingZInProportionToX(t *testing.T) {
	transform := Shearing(0, 0, 0, 0, 1, 0)
	point := tuple.NewPoint(2, 3, 4)

	assert.Equal(t, tuple.NewPoint(2, 3, 6), transform.MulTuple(point))
}

func TestShearingZInProportionToY(t *testing.T) {
	transform := Shearing(0, 0, 0, 0, 0, 1)
	point := tuple.NewPoint(2, 3, 4)

	assert.Equal(t, tuple.NewPoint(2, 3, 7), transform.MulTuple(point))
}

func TestTranformationsInSequence(t *testing.T) {
	p := tuple.NewPoint(1, 0, 1)
	a := RotationX(math.Pi / 2)
	b := Scaling(5, 5, 5)
	c := Translation(10, 5, 7)

	p2 := a.MulTuple(p)
	assert.True(t, tuple.NewPoint(1, -1, 0).Equals(p2))

	p3 := b.MulTuple(p2)
	assert.True(t, tuple.NewPoint(5, -5, 0).Equals(p3))

	p4 := c.MulTuple(p3)
	assert.True(t, tuple.NewPoint(15, 0, 7).Equals(p4))
}

func TestChainedTransformations(t *testing.T) {
	p := tuple.NewPoint(1, 0, 1)
	a := RotationX(math.Pi / 2)
	b := Scaling(5, 5, 5)
	c := Translation(10, 5, 7)

	transform := c.Mul(b).Mul(a)

	assert.True(t, tuple.NewPoint(15, 0, 7).Equals(transform.MulTuple(p)))
}

func TestFluentAPI(t *testing.T) {
	p := tuple.NewPoint(1, 0, 1)
	transform := Identity().
		RotationX(math.Pi/2).
		Scaling(5, 5, 5).
		Translation(10, 5, 7).
		Matrix()

	assert.True(t, tuple.NewPoint(15, 0, 7).Equals(transform.MulTuple(p)))
}

func TestViewTransformDefaultOrientation(t *testing.T) {
	from := tuple.NewPoint(0, 0, 0)
	to := tuple.NewPoint(0, 0, -1)
	up := tuple.NewVector(0, 1, 0)

	tr := ViewTransform(from, to, up)

	assert.Equal(t, matrix.Identity4(), tr)
}

func TestViewTransformInPositiveZDirection(t *testing.T) {
	from := tuple.NewPoint(0, 0, 0)
	to := tuple.NewPoint(0, 0, 1)
	up := tuple.NewVector(0, 1, 0)

	tr := ViewTransform(from, to, up)

	assert.Equal(t, Scaling(-1, 1, -1), tr)
}

func TestViewTransformMovesTheWorld(t *testing.T) {
	from := tuple.NewPoint(0, 0, 8)
	to := tuple.NewPoint(0, 0, 0)
	up := tuple.NewVector(0, 1, 0)

	tr := ViewTransform(from, to, up)

	assert.Equal(t, Translation(0, 0, -8), tr)
}

func TestViewTransformArbitrary(t *testing.T) {
	from := tuple.NewPoint(1, 3, 2)
	to := tuple.NewPoint(4, -2, 8)
	up := tuple.NewVector(1, 1, 0)

	tr := ViewTransform(from, to, up)

	expected := matrix.NewFromSlice([][]float64{
		{-0.50709, 0.50709, 0.67612, -2.36643},
		{0.76772, 0.60609, 0.12122, -2.82843},
		{-0.35857, 0.59761, -0.71714, 0.00000},
		{0.00000, 0.00000, 0.00000, 1.00000},
	})

	assert.True(t, expected.Equals(tr))
}
