package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 10000, 10000
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		wg.Add(1)
		local_py := py
		go func() {
			defer wg.Done()
			y := float64(local_py)/height*(ymax-ymin) + ymin
			for px := 0; px < width; px++ {
				x := float64(px)/width*(xmax-xmin) + xmin
				z := complex(x, y)
				// Image point (px, py) represents complex value z.
				img.Set(px, local_py, mandelbrot(z))
			}
		}()
	}
	wg.Wait()
	png.Encode(os.Stdout, img) // NOTE: ignoring errors
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
