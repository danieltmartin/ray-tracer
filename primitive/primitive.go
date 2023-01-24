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
	SetMaterial(m material.Material)
	SetTransform(t matrix.Matrix)

	NormalAt(worldPoint tuple.Tuple) tuple.Tuple
	Intersects(worldRay ray.Ray) Intersections
}

type Data struct {
	material  material.Material
	transform matrix.Matrix
}

func newData() Data {
	return Data{
		material.Default,
		matrix.Identity4(),
	}
}

func (d *Data) Material() material.Material {
	return d.material
}

func (d *Data) Transform() matrix.Matrix {
	return d.transform
}

func (d *Data) SetTransform(m matrix.Matrix) {
	d.transform = m
}

func (d *Data) SetMaterial(m material.Material) {
	d.material = m
}

func (d *Data) worldPointToLocal(world tuple.Tuple) tuple.Tuple {
	return d.transform.Inverse().MulTuple(world)
}

func (d *Data) localNormalToWorld(localNormal tuple.Tuple) tuple.Tuple {
	worldNormal := d.transform.Inverse().Transpose().MulTuple(localNormal)
	return tuple.New(worldNormal.X, worldNormal.Y, worldNormal.Z, 0).Norm()
}

func (d *Data) worldRayToLocal(r ray.Ray) ray.Ray {
	return r.Transform(d.transform.Inverse())
}

type localIntersecter interface {
	localIntersects(localRay ray.Ray) Intersections
}

func (d *Data) worldIntersects(worldRay ray.Ray, localIntersecter localIntersecter) Intersections {
	localRay := d.worldRayToLocal(worldRay)
	return localIntersecter.localIntersects(localRay)
}

type localNormalizer interface {
	localNormalAt(localPoint tuple.Tuple) tuple.Tuple
}

func (d *Data) worldNormalAt(worldPoint tuple.Tuple, localNormalizer localNormalizer) tuple.Tuple {
	localPoint := d.worldPointToLocal(worldPoint)
	localNormal := localNormalizer.localNormalAt(localPoint)
	return d.localNormalToWorld(localNormal)
}
