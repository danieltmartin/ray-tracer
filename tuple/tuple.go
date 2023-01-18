package tuple

import (
	"math"

	"github.com/danieltmartin/ray-tracer/float"
)

type Tuple struct {
	X, Y, Z, W float64
}

func New(x, y, z, w float64) Tuple {
	return Tuple{x, y, z, w}
}

func NewPoint(x, y, z float64) Tuple {
	return Tuple{x, y, z, 1.0}
}

func NewVector(x, y, z float64) Tuple {
	return Tuple{x, y, z, 0.0}
}

func (t Tuple) Equals(t2 Tuple) bool {
	return float.Equal(t.X, t2.X) &&
		float.Equal(t.Y, t2.Y) &&
		float.Equal(t.Z, t2.Z) &&
		t.W == t2.W
}

func (t Tuple) IsPoint() bool {
	return t.W == 1.0
}

func (t Tuple) IsVector() bool {
	return !t.IsPoint()
}

func (t Tuple) Add(t2 Tuple) Tuple {
	return Tuple{
		t.X + t2.X,
		t.Y + t2.Y,
		t.Z + t2.Z,
		t.W + t2.W,
	}
}

func (t Tuple) Sub(t2 Tuple) Tuple {
	return Tuple{
		t.X - t2.X,
		t.Y - t2.Y,
		t.Z - t2.Z,
		t.W - t2.W,
	}
}

func (t Tuple) Mul(v float64) Tuple {
	return Tuple{
		t.X * v,
		t.Y * v,
		t.Z * v,
		t.W * v,
	}
}

func (t Tuple) Div(v float64) Tuple {
	return Tuple{
		t.X / v,
		t.Y / v,
		t.Z / v,
		t.W / v,
	}
}

func (t Tuple) Mag() float64 {
	return math.Sqrt(t.X*t.X + t.Y*t.Y + t.Z*t.Z)
}

func (t Tuple) Neg() Tuple {
	return Tuple{-t.X, -t.Y, -t.Z, -t.W}
}

func (t Tuple) Norm() Tuple {
	mag := t.Mag()
	return Tuple{
		t.X / mag,
		t.Y / mag,
		t.Z / mag,
		t.W / mag,
	}
}

func (t Tuple) Dot(t2 Tuple) float64 {
	return t.X*t2.X + t.Y*t2.Y + t.Z*t2.Z + t.W*t2.W
}

func (t Tuple) Cross(t2 Tuple) Tuple {
	return Tuple{
		t.Y*t2.Z - t.Z*t2.Y,
		t.Z*t2.X - t.X*t2.Z,
		t.X*t2.Y - t.Y*t2.X,
		0.0,
	}
}

func (t Tuple) Hadamard(t2 Tuple) Tuple {
	return Tuple{
		t.X * t2.X,
		t.Y * t2.Y,
		t.Z * t2.Z,
		t.W * t2.W,
	}
}
