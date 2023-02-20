package obj

import (
	"strings"
	"testing"

	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/tuple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIgnoresUnrecognizedLines(t *testing.T) {
	garbage := `
There was a young lady named Bright
who traveled much faster than light.
She set out one day
in a relative way,
and came back the previous night.
`

	_, err := parse(strings.NewReader(garbage))

	require.NoError(t, err)
}

func TestParseIncompleteVertexRecord(t *testing.T) {
	file := `
v -1 1
`

	_, err := parse(strings.NewReader(file))

	assert.Contains(t, err.Error(), "expected at least 3 arguments")
}

func TestParseVertexRecords(t *testing.T) {
	file := `
v -1 1 0
v -1.0000 0.5000 0.0000
v 1 0 0
v 1 1 0
`

	parsed, err := parse(strings.NewReader(file))

	require.NoError(t, err)
	require.Len(t, parsed.vertices, 4)
	assert.Equal(t, tuple.NewPoint(-1, 1, 0), parsed.vertices[0])
	assert.Equal(t, tuple.NewPoint(-1, 0.5, 0), parsed.vertices[1])
	assert.Equal(t, tuple.NewPoint(1, 0, 0), parsed.vertices[2])
	assert.Equal(t, tuple.NewPoint(1, 1, 0), parsed.vertices[3])
}

func TestParseTriangleFaceRecords(t *testing.T) {
	file := `
v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

f 1 2 3
f 1 3 4
`

	parsed, err := parse(strings.NewReader(file))

	require.NoError(t, err)
	children := parsed.rootGroup.Children()
	require.Len(t, children, 2)

	v1, v2, v3 := children[0].(*primitive.Triangle).Vertices()
	assert.Equal(t, parsed.vertices[0], v1)
	assert.Equal(t, parsed.vertices[1], v2)
	assert.Equal(t, parsed.vertices[2], v3)

	v1, v2, v3 = children[1].(*primitive.Triangle).Vertices()
	assert.Equal(t, parsed.vertices[0], v1)
	assert.Equal(t, parsed.vertices[2], v2)
	assert.Equal(t, parsed.vertices[3], v3)
}

func TestTriangulatingPolygons(t *testing.T) {
	file := `
v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
v 0 2 0

f 1 2 3 4 5
	`

	parsed, err := parse(strings.NewReader(file))

	require.NoError(t, err)
	children := parsed.rootGroup.Children()
	require.Len(t, children, 3)

	v1, v2, v3 := children[0].(*primitive.Triangle).Vertices()
	assert.Equal(t, parsed.vertices[0], v1)
	assert.Equal(t, parsed.vertices[1], v2)
	assert.Equal(t, parsed.vertices[2], v3)
}

func TestTrianglesInGroups(t *testing.T) {
	file := `
v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

g FirstGroup
f 1 2 3
g SecondGroup
f 1 3 4
`

	parsed, err := parse(strings.NewReader(file))

	require.NoError(t, err)
	children := parsed.rootGroup.Children()
	require.Len(t, children, 2)
	g1 := children[0].(*primitive.Group)
	g2 := children[1].(*primitive.Group)
	t1 := g1.Children()[0].(*primitive.Triangle)
	t2 := g2.Children()[0].(*primitive.Triangle)

	v1, v2, v3 := t1.Vertices()
	assert.Equal(t, parsed.vertices[0], v1)
	assert.Equal(t, parsed.vertices[1], v2)
	assert.Equal(t, parsed.vertices[2], v3)

	v1, v2, v3 = t2.Vertices()
	assert.Equal(t, parsed.vertices[0], v1)
	assert.Equal(t, parsed.vertices[2], v2)
	assert.Equal(t, parsed.vertices[3], v3)
}

func TestVertexNormalRecords(t *testing.T) {
	file := `
vn 0 0 1
vn 0.707 0 -0.707
vn 1 2 3
`

	parsed, err := parse(strings.NewReader(file))

	require.NoError(t, err)

	require.Len(t, parsed.normals, 3)
	assert.Equal(t, tuple.NewVector(0, 0, 1), parsed.normals[0])
	assert.Equal(t, tuple.NewVector(0.707, 0, -0.707), parsed.normals[1])
	assert.Equal(t, tuple.NewVector(1, 2, 3), parsed.normals[2])
}

func TestFacesWithNormals(t *testing.T) {
	file := `
v 0 1 0
v -1 0 0
v 1 0 0
vn -1 0 0
vn 1 0 0
vn 0 1 0
f 1//3 2//1 3//2
f 1/0/3 2/102/1 3/14/2
	`

	parsed, err := parse(strings.NewReader(file))
	require.NoError(t, err)

	g := parsed.rootGroup
	t1 := g.Children()[0].(*primitive.SmoothTriangle)
	t2 := g.Children()[1].(*primitive.SmoothTriangle)

	v1, v2, v3 := t1.Vertices()
	assert.Equal(t, parsed.vertices[0], v1)
	assert.Equal(t, parsed.vertices[1], v2)
	assert.Equal(t, parsed.vertices[2], v3)

	n1, n2, n3 := t1.Normals()
	assert.Equal(t, parsed.normals[2], n1)
	assert.Equal(t, parsed.normals[0], n2)
	assert.Equal(t, parsed.normals[1], n3)

	assert.Equal(t, t2, t1)
}
