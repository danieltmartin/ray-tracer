package transform

import (
	"math"

	"github.com/danieltmartin/ray-tracer/matrix"
)

type Transform struct {
	m matrix.Matrix
}

func Identity() Transform {
	return Transform{matrix.Identity4()}
}

func (t Transform) Translation(x, y, z float64) Transform {
	t.m = Translation(x, y, z).Mul(t.m)
	return t
}

func (t Transform) Scaling(x, y, z float64) Transform {
	t.m = Scaling(x, y, z).Mul(t.m)
	return t
}

func (t Transform) RotationX(radians float64) Transform {
	t.m = RotationX(radians).Mul(t.m)
	return t
}

func (t Transform) RotationY(radians float64) Transform {
	t.m = RotationY(radians).Mul(t.m)
	return t
}

func (t Transform) RotationZ(radians float64) Transform {
	t.m = RotationZ(radians).Mul(t.m)
	return t
}

func (t Transform) Shearing(xy, xz, yx, yz, zx, zy float64) Transform {
	t.m = Shearing(xy, xz, yx, yz, zx, zy).Mul(t.m)
	return t
}

func (t Transform) Matrix() matrix.Matrix {
	return t.m.Copy()
}

func Translation(x, y, z float64) matrix.Matrix {
	return matrix.NewFromSlice([][]float64{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	})
}

func Scaling(x, y, z float64) matrix.Matrix {
	if x == 0 || y == 0 || z == 0 {
		panic("cannot scale to 0")
	}
	return matrix.NewFromSlice([][]float64{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	})
}

func RotationX(radians float64) matrix.Matrix {
	return matrix.NewFromSlice([][]float64{
		{1, 0, 0, 0},
		{0, math.Cos(radians), -math.Sin(radians), 0},
		{0, math.Sin(radians), math.Cos(radians), 0},
		{0, 0, 0, 1},
	})
}

func RotationY(radians float64) matrix.Matrix {
	return matrix.NewFromSlice([][]float64{
		{math.Cos(radians), 0, math.Sin(radians), 0},
		{0, 1, 0, 0},
		{-math.Sin(radians), 0, math.Cos(radians), 0},
		{0, 0, 0, 1},
	})
}

func RotationZ(radians float64) matrix.Matrix {
	return matrix.NewFromSlice([][]float64{
		{math.Cos(radians), -math.Sin(radians), 0, 0},
		{math.Sin(radians), math.Cos(radians), 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	})
}

func Shearing(xy, xz, yx, yz, zx, zy float64) matrix.Matrix {
	return matrix.NewFromSlice([][]float64{
		{1, xy, xz, 0},
		{yx, 1, yz, 0},
		{zx, zy, 1, 0},
		{0, 0, 0, 1},
	})
}
