package primitive

import (
	"math"
	"sort"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Group struct {
	children []Primitive
	data
	b bounds
}

func NewGroup() Group {
	return Group{nil, newData(), bounds{}}
}

func (g *Group) Add(p ...Primitive) {
	for _, p := range p {
		if p == g {
			panic("can't add group to itself")
		}
		if math.IsInf(p.bounds().min.X, 0) ||
			math.IsInf(p.bounds().min.Y, 0) ||
			math.IsInf(p.bounds().min.Z, 0) ||
			math.IsInf(p.bounds().max.X, 0) ||
			math.IsInf(p.bounds().max.Y, 0) ||
			math.IsInf(p.bounds().max.Z, 0) {
			panic("primitives with infinite bounds should not be added to a group")
		}
		p.setParent(g)
	}
	g.children = append(g.children, p...)
	g.b = g.calcBounds()
}

func (g *Group) Intersects(worldRay ray.Ray) Intersections {
	return g.worldIntersects(worldRay, g)
}

func (g *Group) NormalAt(worldPoint tuple.Tuple) tuple.Tuple {
	return g.worldNormalAt(worldPoint, g)
}

func (g *Group) localIntersects(localRay ray.Ray) Intersections {
	if !g.intersectsBounds(localRay) {
		return nil
	}

	var xs Intersections
	for _, c := range g.children {
		xs = append(xs, c.Intersects(localRay)...)
	}

	sort.Slice(xs, func(i, j int) bool { return xs[i].distance < xs[j].distance })

	return xs
}

func (g *Group) intersectsBounds(localRay ray.Ray) bool {
	bounds := g.bounds()
	origX, origY, origZ, _ := localRay.Origin().XYZW()
	invDirX, invDirY, invDirZ, _ := localRay.InvDirection().XYZW()

	t1 := (bounds.min.X - origX) * invDirX
	t2 := (bounds.max.X - origX) * invDirX
	t3 := (bounds.min.Y - origY) * invDirY
	t4 := (bounds.max.Y - origY) * invDirY
	t5 := (bounds.min.Z - origZ) * invDirZ
	t6 := (bounds.max.Z - origZ) * invDirZ
	tmin := math.Max(math.Max(math.Min(t1, t2), math.Min(t3, t4)), math.Min(t5, t6))
	tmax := math.Min(math.Min(math.Max(t1, t2), math.Max(t3, t4)), math.Max(t5, t6))

	return tmax >= math.Max(tmin, 0)
}

func (g *Group) localNormalAt(localPoint tuple.Tuple) tuple.Tuple {
	panic("can't compute local normal on a group")
}

func (g *Group) bounds() bounds {
	return g.b
}

func (g *Group) calcBounds() bounds {
	var minX, minY, minZ = math.Inf(1), math.Inf(1), math.Inf(1)
	var maxX, maxY, maxZ = math.Inf(-1), math.Inf(-1), math.Inf(-1)
	for _, c := range g.children {
		b := c.bounds()
		bMin := c.Transform().MulTuple(b.min)
		bMax := c.Transform().MulTuple(b.max)
		for _, p := range []tuple.Tuple{bMin, bMax} {
			minX = math.Min(p.X, minX)
			minY = math.Min(p.Y, minY)
			minZ = math.Min(p.Z, minZ)
			maxX = math.Max(p.X, maxX)
			maxY = math.Max(p.Y, maxY)
			maxZ = math.Max(p.Z, maxZ)
		}
	}

	min := tuple.NewPoint(minX, minY, minZ)
	max := tuple.NewPoint(maxX, maxY, maxZ)

	return bounds{min, max}
}
