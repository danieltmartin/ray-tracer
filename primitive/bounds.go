package primitive

import "github.com/danieltmartin/ray-tracer/tuple"

type bounds struct {
	min, max tuple.Tuple
}

func newBounds(min, max tuple.Tuple) bounds {
	return bounds{min, max}
}
