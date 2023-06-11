package main

import (
	"image"
	"image/color"
	"image/png"
	"math/big"
	"os"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := (float64(py) / height * (ymax - ymin)) + ymin
		for px := 0; px < width; px++ {
			x := (float64(px) / width * (xmax - xmin)) + xmin
			// Image point (px, py) represents complex value z.
			x_, y_ := big.NewFloat(x), big.NewFloat(y)
			img.Set(px, py, mandelbrot(x_, y_))
		}
	}
	png.Encode(os.Stdout, img) // NOTE: ignoring errors
}

func mandelbrot(za, zb *big.Float) color.Color {
	const iterations = 255
	va, vb := new(big.Float), new(big.Float)
	a2, b2, ab := new(big.Float), new(big.Float), new(big.Float)
	two, four := big.NewFloat(2.0), big.NewFloat(4.0)

	for n := uint8(0); n < iterations; n++ {
		// temp variables
		a2.Mul(va, va)              // a^2
		b2.Mul(vb, vb)              // b^2
		ab.Mul(va, vb).Mul(ab, two) // 2ab

		va.Sub(a2, b2).Add(va, za) // a = a^2 - b^2 + real(z)
		vb.Add(ab, zb)             // b = 2ab + img(z)

		// abs, square sum
		a2.Mul(va, va)
		b2.Mul(vb, vb)
		if new(big.Float).Add(b2, a2).Cmp(four) > 0 {
			return calcColor(n)
		}
	}
	return color.Black
}

func calcColor(n uint8) color.Color {
	const contrast = 15
	blue := 255 - contrast*n
	red := 255 - blue
	return color.RGBA{red, 0, blue, 255}
}
