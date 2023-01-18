package tuple

import "github.com/danieltmartin/ray-tracer/float"

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
