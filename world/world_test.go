package world

import (
	"math"
	"sync"
	"testing"

	"github.com/danieltmartin/ray-tracer/float"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/danieltmartin/ray-tracer/light"
	"github.com/danieltmartin/ray-tracer/material"
	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/ray"
	"github.com/danieltmartin/ray-tracer/transform"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyWorld(t *testing.T) {
	w := New()

	assert.Empty(t, w.Primitives())
	assert.Empty(t, w.Lights())
}

func TestRayIntersectWorld(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))

	xs := w.intersect(r)

	assert.Len(t, xs, 4)
	assert.Equal(t, 4.0, xs[0].Distance())
	assert.Equal(t, 4.5, xs[1].Distance())
	assert.Equal(t, 5.5, xs[2].Distance())
	assert.Equal(t, 6.0, xs[3].Distance())
}

func TestPrecomputingIntersectionState(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := primitive.NewSphere()
	i := primitive.NewIntersection(4, &s)

	hc := prepareHitComputations(i, r)

	assert.Equal(t, &s, hc.object)
	assert.Equal(t, tuple.NewPoint(0, 0, -1), hc.hitPoint)
	assert.Equal(t, tuple.NewVector(0, 0, -1), hc.eyev)
	assert.Equal(t, tuple.NewVector(0, 0, -1), hc.normalv)
}

func TestPrecomputingReflectionVector(t *testing.T) {
	s := primitive.NewPlane()
	r := ray.New(tuple.NewPoint(0, 1, -1), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	i := primitive.NewIntersection(math.Sqrt2, &s)

	hc := prepareHitComputations(i, r)

	assert.Equal(t, tuple.NewVector(0, math.Sqrt2/2, math.Sqrt2/2), hc.reflectv)
}

func TestPrecomputingOverpoint(t *testing.T) {
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	s := primitive.NewSphere()
	s.SetTransform(transform.Translation(0, 0, 1))
	i := primitive.NewIntersection(5, &s)

	hc := prepareHitComputations(i, r)

	assert.Less(t, hc.overPoint.Z, -float.Epsilon/2)
	assert.Greater(t, hc.hitPoint.Z, hc.overPoint.Z)
}

func TestHitOutside(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	shape := w.Primitives()[0]
	i := primitive.NewIntersection(4, shape)

	hc := prepareHitComputations(i, r)

	assert.False(t, hc.inside)
}

func TestHitInside(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	shape := w.Primitives()[0]
	i := primitive.NewIntersection(1, shape)

	hc := prepareHitComputations(i, r)

	assert.Equal(t, tuple.NewPoint(0, 0, 1), hc.hitPoint)
	assert.Equal(t, tuple.NewVector(0, 0, -1), hc.eyev)
	assert.True(t, hc.inside)
	assert.Equal(t, tuple.NewVector(0, 0, -1), hc.normalv)
}

func TestShadingAnIntersection(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	shape := w.Primitives()[0]
	i := primitive.NewIntersection(4, shape)

	hc := prepareHitComputations(i, r)
	c := w.shadeHit(hc, 1)

	assert.True(t, floatcolor.New(0.38066, 0.47583, 0.2855).Equals(c))
}

func TestShadingAnIntersectionFromTheInside(t *testing.T) {
	w := testWorld()
	l := light.NewPointLight(tuple.NewPoint(0, 0.25, 0), floatcolor.White)
	w.lights = []*light.PointLight{&l}
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	shape := w.Primitives()[1]
	i := primitive.NewIntersection(0.5, shape)

	hc := prepareHitComputations(i, r)
	c := w.shadeHit(hc, 1)

	assert.True(t, floatcolor.New(0.90498, 0.90498, 0.90498).Equals(c))
}

func TestShadingAnIntersectionInAShadow(t *testing.T) {
	w := New()
	light := light.NewPointLight(tuple.NewPoint(0, 0, -10), floatcolor.White)
	w.AddLights(&light)
	s1 := primitive.NewSphere()
	s2 := primitive.NewSphere()
	s2.SetTransform(transform.Translation(0, 0, 10))
	w.AddPrimitives(&s1, &s2)
	r := ray.New(tuple.NewPoint(0, 0, 5), tuple.NewVector(0, 0, 1))
	i := primitive.NewIntersection(4, &s2)

	hc := prepareHitComputations(i, r)
	c := w.shadeHit(hc, 1)

	assert.True(t, floatcolor.New(0.1, 0.1, 0.1).Equals(c))
}

func TestColorWhenRayMisses(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 1, 0))

	c := w.ColorAt(r, 1)

	assert.Equal(t, floatcolor.Black, c)
}

func TestColorWhenRayHits(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))

	c := w.ColorAt(r, 1)

	assert.True(t, floatcolor.New(0.38066, 0.47583, 0.2855).Equals(c))
}

func TestColorWithIntersectionBehindRay(t *testing.T) {
	w := testWorld()
	outerSphere := w.Primitives()[0]
	outerSphere.SetMaterial(outerSphere.Material().WithAmbient(1))
	innerSphere := w.Primitives()[1]
	innerSphere.SetMaterial(innerSphere.Material().WithAmbient(1))
	r := ray.New(tuple.NewPoint(0, 0, 0.75), tuple.NewVector(0, 0, -1))

	c := w.ColorAt(r, 1)

	assert.True(t, floatcolor.White.Equals(c))
}

func TestNoShadowWhenNothingBetweenIntersectionAndLight(t *testing.T) {
	w := testWorld()
	p := tuple.NewPoint(0, 10, 0)

	assert.False(t, w.isShadowed(p, w.Lights()[0].Position()))
}

func TestShadowWhenObjectIsBetweenIntersectionAndLight(t *testing.T) {
	w := testWorld()
	p := tuple.NewPoint(10, -10, 10)

	assert.True(t, w.isShadowed(p, w.Lights()[0].Position()))
}

func TestNoShadowWhenObjectIsBehindLight(t *testing.T) {
	w := testWorld()
	p := tuple.NewPoint(-20, 20, -20)

	assert.False(t, w.isShadowed(p, w.Lights()[0].Position()))
}

func TestNoShadowWhenObjectIsBehindPoint(t *testing.T) {
	w := testWorld()
	p := tuple.NewPoint(-2, 2, -2)

	assert.False(t, w.isShadowed(p, w.Lights()[0].Position()))
}

func TestReflectedColorForNonReflectiveMaterial(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	shape := w.primitives[1]
	shape.SetMaterial(shape.Material().WithAmbient(1))
	i := primitive.NewIntersection(1, shape)

	hc := prepareHitComputations(i, r)

	assert.Equal(t, floatcolor.Black, w.reflectedColor(hc, 1))
}

func TestReflectedColorForReflectiveMaterial(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	shape := primitive.NewPlane()
	shape.SetMaterial(material.Default.WithReflective(0.5))
	shape.SetTransform(transform.Translation(0, -1, 0))
	w.AddPrimitives(&shape)
	i := primitive.NewIntersection(math.Sqrt2, &shape)

	hc := prepareHitComputations(i, r)

	assertAlmost(t, floatcolor.New(0.19032, 0.2379, 0.14274), w.reflectedColor(hc, 1))
}

func TestShadeHitForReflectiveMaterial(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	shape := primitive.NewPlane()
	shape.SetMaterial(material.Default.WithReflective(0.5))
	shape.SetTransform(transform.Translation(0, -1, 0))
	w.AddPrimitives(&shape)
	i := primitive.NewIntersection(math.Sqrt2, &shape)

	hc := prepareHitComputations(i, r)
	color := w.shadeHit(hc, 1)

	assertAlmost(t, floatcolor.New(0.87677, 0.92436, 0.82918), color)
}

func TestReflectedColorAtMaximumRecursionDepth(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	shape := primitive.NewPlane()
	shape.SetMaterial(material.Default.WithReflective(0.5))
	shape.SetTransform(transform.Translation(0, -1, 0))
	w.AddPrimitives(&shape)
	i := primitive.NewIntersection(math.Sqrt2, &shape)

	hc := prepareHitComputations(i, r)

	assertAlmost(t, floatcolor.Black, w.reflectedColor(hc, 0))
}

func TestColorAtWithMutuallyReflectiveSurfaces(t *testing.T) {
	w := New()
	l := light.NewPointLight(tuple.NewPoint(0, 0, 0), floatcolor.White)
	w.AddLights(&l)
	lower := primitive.NewPlane()
	lower.SetMaterial(material.Default.WithReflective(1))
	lower.SetTransform(transform.Translation(0, -1, 0))
	w.AddPrimitives(&lower)

	upper := primitive.NewPlane()
	upper.SetMaterial(material.Default.WithReflective(1))
	upper.SetTransform(transform.Translation(0, 1, 0))
	w.AddPrimitives(&upper)

	ray := ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 1, 0))
	w.ColorAt(ray, 1)
	// Should not overflow stack from infinite recursion
}

func TestGenerateID(t *testing.T) {
	w := New()
	numIDs := 10000
	ch := make(chan ID, numIDs)
	var wg sync.WaitGroup
	wg.Add(numIDs)

	for i := 0; i < numIDs; i++ {
		go func() {
			ch <- w.NextID()
			wg.Done()
		}()
	}

	wg.Wait()

	last := <-ch
	for i := 1; i < numIDs; i++ {
		next := <-ch
		t.Logf("checking %v and %v\n", last, next)
		require.NotEqual(t, last, next)
		last = next
	}
}

func assertAlmost(t *testing.T, c1 floatcolor.Float64Color, c2 floatcolor.Float64Color) {
	r1, g1, b1 := c1.RGB()
	r2, g2, b2 := c2.RGB()
	assert.True(t, float.AlmostEqual(r1, r2, 0.001), "R values differ: c1.R=%v, c2.R=%v", r1, r2)
	assert.True(t, float.AlmostEqual(g1, g2, 0.001), "G values differ: c1.G=%v, c2.G=%v", g1, g2)
	assert.True(t, float.AlmostEqual(b1, b2, 0.001), "B values differ: c1.B=%v, c2.B=%v", b1, b2)
}

func testWorld() *World {
	w := New()

	light := light.NewPointLight(tuple.NewPoint(-10, 10, -10), floatcolor.White)
	s1 := primitive.NewSphere()
	s1.SetMaterial(material.Default.
		WithColor(floatcolor.New(0.8, 1.0, 0.6)).
		WithDiffuse(0.7).
		WithSpecular(0.2),
	)

	s2 := primitive.NewSphere()
	s2.SetTransform(transform.Scaling(0.5, 0.5, 0.5))

	w.AddLights(&light)
	w.AddPrimitives(&s1, &s2)

	return w
}
