package primitive

import (
	"math"

	"github.com/danieltmartin/ray-tracer/matrix"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

var emptyBox = *NewEmptyBoundingBox()

type BoundingBox struct {
	min, max tuple.Tuple
}

func NewEmptyBoundingBox() *BoundingBox {
	negInf := math.Inf(-1)
	posInf := math.Inf(1)
	return &BoundingBox{
		tuple.NewPoint(posInf, posInf, posInf),
		tuple.NewPoint(negInf, negInf, negInf),
	}
}

func NewBoundingBox(min, max tuple.Tuple) *BoundingBox {
	return &BoundingBox{min, max}
}

func (b *BoundingBox) AddPoint(p tuple.Tuple) {
	b.min = tuple.NewPoint(
		math.Min(b.min.X, p.X),
		math.Min(b.min.Y, p.Y),
		math.Min(b.min.Z, p.Z),
	)
	b.max = tuple.NewPoint(
		math.Max(b.max.X, p.X),
		math.Max(b.max.Y, p.Y),
		math.Max(b.max.Z, p.Z),
	)
}

func (b *BoundingBox) AddBox(b2 *BoundingBox) {
	if b2 == nil || *b2 == emptyBox || b2.hasNaN() {
		return
	}
	b.AddPoint(b2.min)
	b.AddPoint(b2.max)
}

func (b *BoundingBox) hasNaN() bool {
	return math.IsNaN(b.min.X) ||
		math.IsNaN(b.min.Y) ||
		math.IsNaN(b.min.Z) ||
		math.IsNaN(b.max.X) ||
		math.IsNaN(b.max.Y) ||
		math.IsNaN(b.max.Z)
}

func (b *BoundingBox) ContainsPoint(point tuple.Tuple) bool {
	return b.min.X <= point.X && point.X <= b.max.X &&
		b.min.Y <= point.Y && point.Y <= b.max.Y &&
		b.min.Z <= point.Z && point.Z <= b.max.Z
}

func (b *BoundingBox) ContainsBox(b2 *BoundingBox) bool {
	return b.ContainsPoint(b2.min) && b.ContainsPoint(b2.max)
}

func (b *BoundingBox) intersects(ray ray.Ray) bool {
	origX, origY, origZ, _ := ray.Origin().XYZW()
	invDirX, invDirY, invDirZ, _ := ray.InvDirection().XYZW()

	t1 := (b.min.X - origX) * invDirX
	t2 := (b.max.X - origX) * invDirX
	t3 := (b.min.Y - origY) * invDirY
	t4 := (b.max.Y - origY) * invDirY
	t5 := (b.min.Z - origZ) * invDirZ
	t6 := (b.max.Z - origZ) * invDirZ
	tmin := math.Max(math.Max(math.Min(t1, t2), math.Min(t3, t4)), math.Min(t5, t6))
	tmax := math.Min(math.Min(math.Max(t1, t2), math.Max(t3, t4)), math.Max(t5, t6))

	return tmax >= math.Max(tmin, 0)
}

func (b *BoundingBox) Transform(m matrix.Matrix) *BoundingBox {
	corners := []tuple.Tuple{
		m.MulTuple(b.min),
		m.MulTuple(b.max),
		m.MulTuple(tuple.NewPoint(b.min.X, b.min.Y, b.max.Z)),
		m.MulTuple(tuple.NewPoint(b.min.X, b.max.Y, b.min.Z)),
		m.MulTuple(tuple.NewPoint(b.min.X, b.max.Y, b.max.Z)),
		m.MulTuple(tuple.NewPoint(b.max.X, b.min.Y, b.min.Z)),
		m.MulTuple(tuple.NewPoint(b.max.X, b.min.Y, b.max.Z)),
		m.MulTuple(tuple.NewPoint(b.max.X, b.max.Y, b.min.Z)),
	}

	newBox := NewEmptyBoundingBox()
	for _, p := range corners {
		newBox.AddPoint(p)
	}

	return newBox
}

func (b *BoundingBox) Min() tuple.Tuple {
	return b.min
}

func (b *BoundingBox) Max() tuple.Tuple {
	return b.max
}
