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
	Parent() *Group
	NormalAt(worldPoint tuple.Tuple) tuple.Tuple
	Intersects(worldRay ray.Ray) Intersections
	WorldPointToLocal(worldPoint tuple.Tuple) tuple.Tuple

	Bounds() Bounds
	setParent(g *Group)
}

type data struct {
	material         material.Material
	transform        matrix.Matrix
	inverseTransform matrix.Matrix
	parent           *Group
}

func newData() data {
	ident := matrix.Identity4()
	return data{
		material.Default,
		ident,
		ident,
		nil,
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

func (d *data) WorldPointToLocal(world tuple.Tuple) tuple.Tuple {
	if d.parent != nil {
		world = d.parent.WorldPointToLocal(world)
	}
	return d.inverseTransform.MulTuple(world)
}

func (d *data) localNormalToWorld(localNormal tuple.Tuple) tuple.Tuple {
	transformed := d.inverseTransform.Transpose().MulTuple(localNormal)
	worldNormal := tuple.New(transformed.X, transformed.Y, transformed.Z, 0).Norm()

	if d.parent != nil {
		worldNormal = d.parent.localNormalToWorld(worldNormal)
	}
	return worldNormal
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

func (d *data) Parent() *Group {
	return d.parent
}

func (d *data) setParent(g *Group) {
	d.parent = g
}

type localNormalizer interface {
	localNormalAt(localPoint tuple.Tuple) tuple.Tuple
}

func (d *data) worldNormalAt(worldPoint tuple.Tuple, localNormalizer localNormalizer) tuple.Tuple {
	localPoint := d.WorldPointToLocal(worldPoint)
	localNormal := localNormalizer.localNormalAt(localPoint)
	return d.localNormalToWorld(localNormal)
}
