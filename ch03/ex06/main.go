package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
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
			img.Set(px, py, mandelbrot(z))
		}
	}
	// loop image pixel and average neighbors colors (supersampling)
	smoothedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			smoothedImg.Set(px, py, smoothPixel(img, px, py))
		}
	}
	png.Encode(os.Stdout, smoothedImg) // NOTE: ignoring errors
}

func smoothPixel(srcImg image.Image, px, py int) color.RGBA {
	const delta = 3
	var sr, sg, sb, n uint32
	var r, g, b uint32

	for i := px - delta; i <= px+delta; i++ {
		for j := py - delta; j <= py+delta; j++ {
			sr, sg, sb, _ = srcImg.At(px, py).RGBA()
			r += sr
			g += sg
			b += sb
			n++
		}
	}
	return color.RGBA{uint8(r / n >> 8), uint8(g / n >> 8), uint8(b / n >> 8), 255}
}

func mandelbrot(z complex128) color.Color {
	const iterations = 255
	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
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
