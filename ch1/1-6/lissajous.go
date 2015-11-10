package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

var palette = []color.Color{
	color.Black,
}

const (
	blackIndex = 0
	greenIndex = 1
	blueIndex = 2
)

func genGradientPalette(n int) {
	ci := color.RGBA{G: 0xFF, A: 0x01}
	var stepSize uint8 = 0xFF % uint8(n)
	for i := 0; i < n; i++ {
		ci = color.RGBA{G: ci.G - stepSize, B: ci.B + stepSize, A: ci.A}
		palette = append(palette, ci)
	}
}

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 128
		delay   = 24
	)
	genGradientPalette(nframes)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(i+1))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
