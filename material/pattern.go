package material

import (
	"math"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Pattern interface {
	colorAtObject(object Object, point tuple.Tuple) floatcolor.Float64Color
	inverseTransform() matrix.Matrix
}

type patternTransform matrix.Matrix

func (p patternTransform) inverseTransform() matrix.Matrix {
	return matrix.Matrix(p)
}

type SolidPattern floatcolor.Float64Color

func (s SolidPattern) colorAtObject(object Object, point tuple.Tuple) floatcolor.Float64Color {
	return floatcolor.Float64Color(s)
}

func (s SolidPattern) inverseTransform() matrix.Matrix {
	return matrix.Identity4()
}

type StripePattern struct {
	color1, color2 floatcolor.Float64Color
	patternTransform
}

func NewStripePattern(color1, color2 floatcolor.Float64Color) StripePattern {
	return StripePattern{color1, color2, patternTransform(matrix.Identity4())}
}

func (s StripePattern) WithTransform(transform matrix.Matrix) StripePattern {
	return StripePattern{s.color1, s.color2, patternTransform(transform.Inverse())}
}

func (p StripePattern) colorAt(point tuple.Tuple) floatcolor.Float64Color {
	if int(math.Floor(point.X))%2 == 0 {
		return p.color1
	}
	return p.color2
}

func (p StripePattern) colorAtObject(object Object, worldPoint tuple.Tuple) floatcolor.Float64Color {
	return p.colorAt(toPatternPoint(p, object, worldPoint))
}

type GradientPattern struct {
	fromColor, toColor floatcolor.Float64Color
	patternTransform
}

func NewGradientPattern(fromColor, toColor floatcolor.Float64Color) GradientPattern {
	return GradientPattern{fromColor, toColor, patternTransform(matrix.Identity4())}
}

func (g GradientPattern) WithTransform(transform matrix.Matrix) GradientPattern {
	return GradientPattern{g.fromColor, g.toColor, patternTransform(transform.Inverse())}
}

func (s GradientPattern) colorAt(point tuple.Tuple) floatcolor.Float64Color {
	absX := math.Abs(point.X)
	floorX := math.Floor(absX)
	fraction := absX - floorX
	if int(floorX)%2 == 1 {
		// Go the opposite direction on odd numbers to create a repeating effect instead of
		// a hard transition from the end color back to the beginning color.
		fraction = 1 - fraction
	}
	return floatcolor.Lerp(s.fromColor, s.toColor, fraction)
}

func (p GradientPattern) colorAtObject(object Object, worldPoint tuple.Tuple) floatcolor.Float64Color {
	return p.colorAt(toPatternPoint(p, object, worldPoint))
}

type RingPattern struct {
	fromColor, toColor floatcolor.Float64Color
	patternTransform
}

func NewRingPattern(fromColor, toColor floatcolor.Float64Color) RingPattern {
	return RingPattern{fromColor, toColor, patternTransform(matrix.Identity4())}
}

func (r RingPattern) WithTransform(transform matrix.Matrix) RingPattern {
	return RingPattern{r.fromColor, r.toColor, patternTransform(transform.Inverse())}
}

func (r RingPattern) colorAt(point tuple.Tuple) floatcolor.Float64Color {
	d := math.Sqrt(point.X*point.X + point.Z*point.Z)
	if int(math.Floor(d))%2 == 0 {
		return r.fromColor
	}
	return r.toColor
}

func (r RingPattern) colorAtObject(object Object, worldPoint tuple.Tuple) floatcolor.Float64Color {
	return r.colorAt(toPatternPoint(r, object, worldPoint))
}

type CheckerPattern struct {
	color1, color2 floatcolor.Float64Color
	patternTransform
}

func NewCheckerPattern(color1, color2 floatcolor.Float64Color) CheckerPattern {
	return CheckerPattern{color1, color2, patternTransform(matrix.Identity4())}
}

func (r CheckerPattern) WithTransform(transform matrix.Matrix) CheckerPattern {
	return CheckerPattern{r.color1, r.color2, patternTransform(transform.Inverse())}
}

func (r CheckerPattern) colorAt(point tuple.Tuple) floatcolor.Float64Color {
	x, y, z, _ := point.XYZW()
	if int(math.Floor(x)+math.Floor(y)+math.Floor(z))%2 == 0 {
		return r.color1
	}
	return r.color2
}

func (r CheckerPattern) colorAtObject(object Object, worldPoint tuple.Tuple) floatcolor.Float64Color {
	patternPoint := toPatternPoint(r, object, worldPoint)
	return r.colorAt(patternPoint)
}

// TestPattern returns colors with the RGB values set to the XYZ of the point of intersection.
// This is useful for tests to determine if a ray was refracted.
type TestPattern struct{}

func (TestPattern) colorAtObject(object Object, point tuple.Tuple) floatcolor.Float64Color {
	return floatcolor.New(point.X, point.Y, point.Z)
}

func (TestPattern) inverseTransform() matrix.Matrix {
	return matrix.Identity4()
}

type Object interface {
	InverseTransform() matrix.Matrix
}

func toPatternPoint(pattern Pattern, object Object, worldPoint tuple.Tuple) tuple.Tuple {
	objectPoint := object.InverseTransform().MulTuple(worldPoint)
	return pattern.inverseTransform().MulTuple(objectPoint)
}
