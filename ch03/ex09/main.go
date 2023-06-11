package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"net/url"
	"strconv"
)

func main() {
	http.HandleFunc("/", handler) // each request calls handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	xmin, xmax, ymin, ymax := -2, +2, -2, +2

	queries, err := url.ParseQuery(r.URL.RawQuery)
	if err == nil {
		if v, ok := queries["xmin"]; ok {
			xmin, _ = strconv.Atoi(v[0])
		}
		if v, ok := queries["xmax"]; ok {
			xmax, _ = strconv.Atoi(v[0])
		}
		if v, ok := queries["ymin"]; ok {
			ymin, _ = strconv.Atoi(v[0])
		}
		if v, ok := queries["ymax"]; ok {
			ymax, _ = strconv.Atoi(v[0])
		}
	}
	genMandelbortPNG(w, xmin, xmax, ymin, ymax)
}

func genMandelbortPNG(out io.Writer, xmin, xmax, ymin, ymax int) {
	const (
		width, height = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := (float64(py) / height * float64(ymax-ymin)) + float64(ymin)
		for px := 0; px < width; px++ {
			x := (float64(px) / width * float64(xmax-xmin)) + float64(xmin)
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(out, img) // NOTE: ignoring errors
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
