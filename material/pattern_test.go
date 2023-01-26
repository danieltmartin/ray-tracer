package material

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
)

func TestSolidPatternConstant(t *testing.T) {
	p := SolidPattern(floatcolor.Blue)

	assert.Equal(t, floatcolor.Blue, p.colorAtObject(obj, tuple.NewPoint(0, 0, 0)))
	assert.Equal(t, floatcolor.Blue, p.colorAtObject(obj, tuple.NewPoint(1, 1, 1)))
	assert.Equal(t, floatcolor.Blue, p.colorAtObject(obj, tuple.NewPoint(-1, -1, -1)))
	assert.Equal(t, floatcolor.Blue, p.colorAtObject(obj, tuple.NewPoint(500, -300, 800)))
}

func TestStripePatternConstantInY(t *testing.T) {
	p := NewStripePattern(floatcolor.White, floatcolor.Black)

	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 1, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 2, 0)))
}

func TestStripePatternIsConstantInZ(t *testing.T) {
	p := NewStripePattern(floatcolor.White, floatcolor.Black)

	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 1)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 2)))
}

func TestStripePatternAlternatesInX(t *testing.T) {
	p := NewStripePattern(floatcolor.White, floatcolor.Black)

	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0.9, 0, 0)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(1, 0, 0)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(-0.1, 0, 0)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(-1, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(-1.1, 0, 0)))
}

func TestStripePatternWithObjectTransformation(t *testing.T) {
	o := dummyObject(transform.Scaling(2, 2, 2))
	p := NewStripePattern(floatcolor.White, floatcolor.Black)

	c := p.colorAtObject(&o, tuple.NewPoint(1.5, 0, 0))

	assert.Equal(t, floatcolor.White, c)
}

func TestStripePatternWithPatternTransformation(t *testing.T) {
	o := dummyObject(matrix.Identity4())
	p := NewStripePattern(floatcolor.White, floatcolor.Black).
		WithTransform(transform.Scaling(2, 2, 2))

	c := p.colorAtObject(&o, tuple.NewPoint(1.5, 0, 0))

	assert.Equal(t, floatcolor.White, c)
}

func TestStripePatternWithObjectAndPatternTransformation(t *testing.T) {
	o := dummyObject(transform.Scaling(2, 2, 2))
	p := NewStripePattern(floatcolor.White, floatcolor.Black).
		WithTransform(transform.Translation(0.5, 0, 0))

	c := p.colorAtObject(&o, tuple.NewPoint(2.5, 0, 0))

	assert.Equal(t, floatcolor.White, c)
}

func TestGradientPatternLinearlyInterpolatesBetweenColors(t *testing.T) {
	p := NewGradientPattern(floatcolor.White, floatcolor.Black)

	// Interpolation is implemented in HCL color space so it's not exactly linear
	// in terms of the numerical values, therefore we're using a very liberal epsilon value.
	assertColorsEqual(t, floatcolor.New(1, 1, 1), p.colorAt(tuple.NewPoint(0, 0, 0)))
	assertColorsEqual(t, floatcolor.New(0.75, 0.75, 0.75), p.colorAt(tuple.NewPoint(0.25, 0, 0)))
	assertColorsEqual(t, floatcolor.New(0.5, 0.5, 0.5), p.colorAt(tuple.NewPoint(0.5, 0, 0)))
	assertColorsEqual(t, floatcolor.New(0.25, 0.25, 0.25), p.colorAt(tuple.NewPoint(0.75, 0, 0)))
}

func TestRingPatternExtendsBothInXAndZ(t *testing.T) {
	p := NewRingPattern(floatcolor.White, floatcolor.Black)

	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 0)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(1, 0, 0)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(0, 0, 1)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(0.708, 0, 0.708)))
}

func TestCheckerPatternShouldRepeatInX(t *testing.T) {
	p := NewCheckerPattern(floatcolor.White, floatcolor.Black)

	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0.99, 0, 0)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(1.01, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(2.01, 0, 0)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(-0.99, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(-1.01, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(-1.99, 0, 0)))
}

func TestCheckerPatternShouldRepeatInY(t *testing.T) {
	p := NewCheckerPattern(floatcolor.White, floatcolor.Black)

	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0.99, 0)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(0, 1.01, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 2.01, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, -1.01, 0)))
}

func TestCheckerPatternShouldRepeatInZ(t *testing.T) {
	p := NewCheckerPattern(floatcolor.White, floatcolor.Black)

	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 0)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 0.99)))
	assert.Equal(t, floatcolor.Black, p.colorAt(tuple.NewPoint(0, 0, 1.01)))
	assert.Equal(t, floatcolor.White, p.colorAt(tuple.NewPoint(0, 0, 2.01)))
}

func assertColorsEqual(t *testing.T, expected, actual floatcolor.Float64Color) {
	assert.Truef(t, expected.AlmostEquals(actual, 0.1), "expected %v to equal %v", expected, actual)
}

type dummyObject matrix.Matrix

func (d dummyObject) Transform() matrix.Matrix {
	return matrix.Matrix(d)
}
