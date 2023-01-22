package canvas

import (
	"testing"

	"github.com/danieltmartin/ray-tracer/floatcolor"
	"github.com/stretchr/testify/assert"
)

func TestCreateCanvas(t *testing.T) {
	c := New(10, 20)

	assert.Equal(t, c.Width(), uint(10))
	assert.Equal(t, c.Height(), uint(20))

	for x := uint(0); x < 10; x++ {
		for y := uint(0); y < 20; y++ {
			assert.Equal(t, floatcolor.Black, c.PixelAt(x, y))
		}
	}
}

func TestWritePixel(t *testing.T) {
	c := New(10, 20)

	c.WritePixel(2, 3, floatcolor.Red)

	assert.Equal(t, floatcolor.Red, c.PixelAt(2, 3))
}
