package primitive

import (
	"sort"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Group struct {
	children []Primitive
	data
	bounds BoundingBox
}

func NewGroup() *Group {
	var g Group
	g.data = newData()
	g.bounds = *NewEmptyBoundingBox()
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
		p.setParent(g)
		g.bounds.AddBox(p.Bounds().Transform(p.Transform()))
	}
	g.children = append(g.children, p...)
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
	if !g.bounds.intersects(localRay) {
		return nil
	}

	var xs Intersections
	for _, c := range g.children {
		xs = append(xs, c.Intersects(localRay)...)
	}

	sort.Slice(xs, func(i, j int) bool { return xs[i].distance < xs[j].distance })

	return xs
}

func (g *Group) localNormalAt(localPoint tuple.Tuple, _ Intersection) tuple.Tuple {
	panic("can't compute local normal on a group")
}

func (g *Group) Bounds() *BoundingBox {
	return &g.bounds
}
