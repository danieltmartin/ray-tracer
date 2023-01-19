package ppm

import (
	"bytes"
	"strings"
	"testing"

	"github.com/danieltmartin/ray-tracer/canvas"
	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	var b bytes.Buffer
	c := canvas.New(5, 3)

	c.WritePixel(4, 2, floatcolor.New(-0.5, 0, 1))
	err := Encode(&b, c)

	expected := `P3
5 3
65535
`

	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(b.String(), expected))
}

func TestPixelData(t *testing.T) {
	var b bytes.Buffer

	c := canvas.New(5, 3)
	c.WritePixel(0, 0, floatcolor.New(1.5, 0, 0))
	c.WritePixel(2, 1, floatcolor.New(0, 0.5, 0))
	c.WritePixel(4, 2, floatcolor.New(-0.5, 0, 1))

	err := Encode(&b, c)

	expected := `P3
5 3
65535
65535 0 0 0 0 0 0 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 32767 0 0 0 0 0 0 0
0 0 0 0 0 0 0 0 0 0 0 0 0 0 65535

`

	assert.NoError(t, err)
	assert.Equal(t, expected, b.String())
}

func TestSplitLongLines(t *testing.T) {
	var b bytes.Buffer

	c := canvas.New(10, 2)
	color := floatcolor.New(1, 0.8, 0.6)
	for x := 0; x < 10; x++ {
		for y := 0; y < 2; y++ {
			c.WritePixel(x, y, color)
		}
	}

	err := Encode(&b, c)

	expected := `P3
10 2
65535
65535 52428 39321 65535 52428 39321 65535 52428 39321 65535 52428
39321 65535 52428 39321 65535 52428 39321 65535 52428 39321 65535
52428 39321 65535 52428 39321 65535 52428 39321
65535 52428 39321 65535 52428 39321 65535 52428 39321 65535 52428
39321 65535 52428 39321 65535 52428 39321 65535 52428 39321 65535
52428 39321 65535 52428 39321 65535 52428 39321

`

	assert.NoError(t, err)
	assert.Equal(t, expected, b.String())
}
