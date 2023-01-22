package ray

import (
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Ray struct {
	origin    tuple.Tuple
	direction tuple.Tuple
}

func New(origin tuple.Tuple, direction tuple.Tuple) Ray {
	return Ray{origin, direction}
}

func (r Ray) Origin() tuple.Tuple {
	return r.origin
}

func (r Ray) Direction() tuple.Tuple {
	return r.direction
}

func (r Ray) Position(t float64) tuple.Tuple {
	return r.origin.Add((r.direction).Mul(t))
}

func (r Ray) Transform(t matrix.Matrix) Ray {
	return Ray{t.MulTuple(r.origin), t.MulTuple(r.direction)}
}
