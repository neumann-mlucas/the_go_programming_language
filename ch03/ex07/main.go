package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, newtow(z))
		}
	}
	png.Encode(os.Stdout, img) // NOTE: ignoring errors
}

func newtow(z complex128) color.Color {
	const iterations = 255
	const lim = 1e-6

	v := z - fn(z)/dfn(z)
	for n := uint8(0); n < iterations; n++ {
		v = v - fn(v)/dfn(v)
		if cmplx.Abs(fn(v)) < lim {
			return calcColor(n)
		}
	}
	return color.Black
}

func fn(x complex128) complex128 {
	return cmplx.Pow(x, 3) - 1
}

func dfn(x complex128) complex128 {
	return 3*cmplx.Pow(x, 2) - 1

}
func calcColor(n uint8) color.Color {
	const contrast = 15
	blue := 255 - contrast*n
	red := 255 - blue
	return color.RGBA{red, 0, blue, 255}
}
