package primitive

import (
	"sort"

	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Group struct {
	children []Primitive
	data
}

func NewGroup() Group {
	return Group{nil, newData()}
}

func (g *Group) Add(p ...Primitive) {
	for _, p := range p {
		if p == g {
			panic("can't add group to itself")
		}
		p.setParent(g)
	}
	g.children = append(g.children, p...)
}

func (g *Group) Intersects(worldRay ray.Ray) Intersections {
	return g.worldIntersects(worldRay, g)
}

func (g *Group) NormalAt(worldPoint tuple.Tuple) tuple.Tuple {
	return g.worldNormalAt(worldPoint, g)
}

func (g *Group) localIntersects(localRay ray.Ray) Intersections {
	var xs Intersections
	for _, c := range g.children {
		xs = append(xs, c.Intersects(localRay)...)
	}

	sort.Slice(xs, func(i, j int) bool { return xs[i].distance < xs[j].distance })

	return xs
}

func (g *Group) localNormalAt(localPoint tuple.Tuple) tuple.Tuple {
	panic("can't compute local normal on a group")
}
