package primitive

import (
	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Primitive interface {
	Material() material.Material
	NormalAt(worldPoint tuple.Tuple) tuple.Tuple
	Intersects(r ray.Ray) Intersections
	SetMaterial(m material.Material)
}
