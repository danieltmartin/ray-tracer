package primitive

import (
	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Primitive interface {
	Material() material.Material
	Transform() matrix.Matrix
	InverseTransform() matrix.Matrix
	SetMaterial(m material.Material)
	SetTransform(t matrix.Matrix)

	NormalAt(worldPoint tuple.Tuple) tuple.Tuple
	Intersects(worldRay ray.Ray) Intersections
}

type data struct {
	material         material.Material
	transform        matrix.Matrix
	inverseTransform matrix.Matrix
}

func newData() data {
	ident := matrix.Identity4()
	return data{
		material.Default,
		ident,
		ident,
	}
}

func (d *data) Material() material.Material {
	return d.material
}

func (d *data) Transform() matrix.Matrix {
	return d.transform
}

func (d *data) InverseTransform() matrix.Matrix {
	return d.inverseTransform
}

func (d *data) SetTransform(m matrix.Matrix) {
	d.transform = m
	d.inverseTransform = m.Inverse()
}

func (d *data) SetMaterial(m material.Material) {
	d.material = m
}

func (d *data) worldPointToLocal(world tuple.Tuple) tuple.Tuple {
	return d.inverseTransform.MulTuple(world)
}

func (d *data) localNormalToWorld(localNormal tuple.Tuple) tuple.Tuple {
	worldNormal := d.inverseTransform.Transpose().MulTuple(localNormal)
	return tuple.New(worldNormal.X, worldNormal.Y, worldNormal.Z, 0).Norm()
}

func (d *data) worldRayToLocal(r ray.Ray) ray.Ray {
	return r.Transform(d.inverseTransform)
}

type localIntersecter interface {
	localIntersects(localRay ray.Ray) Intersections
}

func (d *data) worldIntersects(worldRay ray.Ray, localIntersecter localIntersecter) Intersections {
	localRay := d.worldRayToLocal(worldRay)
	return localIntersecter.localIntersects(localRay)
}

type localNormalizer interface {
	localNormalAt(localPoint tuple.Tuple) tuple.Tuple
}

func (d *data) worldNormalAt(worldPoint tuple.Tuple, localNormalizer localNormalizer) tuple.Tuple {
	localPoint := d.worldPointToLocal(worldPoint)
	localNormal := localNormalizer.localNormalAt(localPoint)
	return d.localNormalToWorld(localNormal)
}
