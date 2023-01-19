package ppm

import (
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"
)

const maxLineLength = 70

func Encode(w io.Writer, m image.Image) error {
	io.WriteString(w, "P3\n")
	bounds := m.Bounds()
	io.WriteString(w, fmt.Sprintf("%v %v\n", bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y))
	io.WriteString(w, "65535\n")

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var line strings.Builder
		flushLine := func() {
			io.WriteString(w, line.String())
			io.WriteString(w, "\n")
			line.Reset()
		}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

			for _, v := range []uint32{r, g, b} {
				s := strconv.FormatUint(uint64(v), 10)
				if line.Len()+len(s)+1 > maxLineLength {
					flushLine()
				} else if line.Len() > 0 {
					line.WriteRune(' ')
				}
				line.WriteString(s)
			}
		}
		if line.Len() > 0 {
			flushLine()
		}
	}

	io.WriteString(w, "\n")

	return nil
}
