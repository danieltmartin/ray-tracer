package world

import (
	"sync"
	"testing"

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

	xs := w.Intersect(r)

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
	c := w.shadeHit(hc)

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
	c := w.shadeHit(hc)

	assert.True(t, floatcolor.New(0.90498, 0.90498, 0.90498).Equals(c))
}

func TestColorWhenRayMisses(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 1, 0))

	c := w.ColorAt(r)

	assert.Equal(t, floatcolor.Black, c)
}

func TestColorWhenRayHits(t *testing.T) {
	w := testWorld()
	r := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))

	c := w.ColorAt(r)

	assert.True(t, floatcolor.New(0.38066, 0.47583, 0.2855).Equals(c))
}

func TestColorWithIntersectionBehindRay(t *testing.T) {
	w := testWorld()
	outerSphere := w.Primitives()[0]
	outerSphere.SetMaterial(outerSphere.Material().WithAmbient(1))
	innerSphere := w.Primitives()[1]
	innerSphere.SetMaterial(innerSphere.Material().WithAmbient(1))
	r := ray.New(tuple.NewPoint(0, 0, 0.75), tuple.NewVector(0, 0, -1))

	c := w.ColorAt(r)

	assert.True(t, innerSphere.Material().Color().Equals(c))
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
