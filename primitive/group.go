package primitive

import (
	"math"
	"sort"
	"sync"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Group struct {
	children []Primitive
	data
	b           Bounds
	boundsDirty bool
	mut         sync.Mutex
}

func NewGroup() *Group {
	var g Group
	g.data = newData()
	return &g
}

func (g *Group) Add(p ...Primitive) {
	if len(p) == 0 {
		return
	}
	for _, p := range p {
		if p == g {
			panic("can't add group to itself")
		}
		if math.IsInf(p.Bounds().min.X, 0) ||
			math.IsInf(p.Bounds().min.Y, 0) ||
			math.IsInf(p.Bounds().min.Z, 0) ||
			math.IsInf(p.Bounds().max.X, 0) ||
			math.IsInf(p.Bounds().max.Y, 0) ||
			math.IsInf(p.Bounds().max.Z, 0) {
			panic("primitives with infinite bounds should not be added to a group")
		}
		p.setParent(g)
	}
	g.children = append(g.children, p...)
	g.boundsDirty = true
}

func (g *Group) Children() []Primitive {
	return g.children
}

func (g *Group) Intersects(worldRay ray.Ray) Intersections {
	return g.worldIntersects(worldRay, g)
}

func (g *Group) NormalAt(worldPoint tuple.Tuple, xn Intersection) tuple.Tuple {
	return g.worldNormalAt(worldPoint, xn, g)
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
	bounds := g.Bounds()
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

func (g *Group) localNormalAt(localPoint tuple.Tuple, _ Intersection) tuple.Tuple {
	panic("can't compute local normal on a group")
}

func (g *Group) Bounds() Bounds {
	g.mut.Lock()
	defer g.mut.Unlock()
	if g.boundsDirty {
		g.b = g.calcBounds()
		g.boundsDirty = false
	}
	return g.b
}

func (g *Group) calcBounds() Bounds {
	if len(g.children) == 0 {
		return g.b
	}
	var minX, minY, minZ = math.Inf(1), math.Inf(1), math.Inf(1)
	var maxX, maxY, maxZ = math.Inf(-1), math.Inf(-1), math.Inf(-1)
	for _, c := range g.children {
		b := c.Bounds()
		t := c.Transform()
		corners := []tuple.Tuple{
			t.MulTuple(b.min),
			t.MulTuple(b.max),
			t.MulTuple(tuple.NewPoint(b.min.X, b.min.Y, b.max.Z)),
			t.MulTuple(tuple.NewPoint(b.min.X, b.max.Y, b.min.Z)),
			t.MulTuple(tuple.NewPoint(b.min.X, b.max.Y, b.max.Z)),
			t.MulTuple(tuple.NewPoint(b.max.X, b.min.Y, b.min.Z)),
			t.MulTuple(tuple.NewPoint(b.max.X, b.min.Y, b.max.Z)),
			t.MulTuple(tuple.NewPoint(b.max.X, b.max.Y, b.min.Z)),
		}
		for _, p := range corners {
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

	return Bounds{min, max}
}
