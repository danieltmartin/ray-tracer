package primitive

import "github.com/danieltmartin/ray-tracer/tuple"

type Bounds struct {
	min, max tuple.Tuple
}

func (b Bounds) Min() tuple.Tuple {
	return b.min
}

func (b Bounds) Max() tuple.Tuple {
	return b.max
}

func newBounds(min, max tuple.Tuple) Bounds {
	return Bounds{min, max}
}
