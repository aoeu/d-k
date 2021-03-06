package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	var host string
	flag.StringVar(&host, "host", "localhost:8000", "The hostname and port to serve data from over HTTP.")
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		}
		s := newLissajous()
		msg := ""
		for name, value := range r.Form {
			var i int
			var f float64
			var err error
			errFmtStr := "Error converting %s value %v : %v\n"
			n := strings.ToLower(name) // TODO(aoeu): Keep case insensitivity like headers?
			switch n {
			case "cycles", "size", "nframes", "delay":
				if i, err = strconv.Atoi(value[0]); err != nil {
					msg += fmt.Sprintf(errFmtStr, name, value[0], err)
					continue
				}
			case "res", "freq":
				if f, err = strconv.ParseFloat(value[0], 64); err != nil {
					msg += fmt.Sprintf(errFmtStr, name, value[0], err)
					continue
				}
			case "help":
				msg += "Add 'size', 'cycles', 'nframes', 'delay', and 'res' as query string parameters:"
			}
			switch n {
			case "cycles":
				s.cycles = i
			case "size":
				s.size = i
			case "nframes":
				s.nframes = i
			case "delay":
				s.delay = i
			case "res":
				s.res = f
			case "freq":
				s.freq = f
			}
		}
		log.Printf("%v\n", s.URL(host))
		fmt.Fprintln(os.Stderr, msg)
		s.renderAnim(w)
	})
	log.Fatal(http.ListenAndServe(host, nil))
}

var palette = []color.Color{
	color.Black,
	color.RGBA{G: 0xFF, A: 0x01},
}

const (
	blackIndex = 0
	greenIndex = 1
)

type lissajous struct {
	cycles  int
	res     float64
	size    int
	nframes int
	delay   int
	freq    float64
}

func newLissajous() *lissajous {
	return &lissajous{
		cycles:  5,
		res:     0.001,
		size:    100,
		nframes: 64,
		delay:   8,
		freq:    rand.Float64() * 3.0,
	}
}

func (s lissajous) URL(host string) string {
	// TODO(aoeu): Would url.URL provide implementation advantages?
	return fmt.Sprintf("%v?size=%v&cycles=%v&nframes=%v&delay=%v&res=%v&freq=%v",
		host, s.size, s.cycles, s.nframes, s.delay, s.res, s.freq)
}

func (s *lissajous) renderAnim(out io.Writer) {
	anim := gif.GIF{LoopCount: s.nframes}
	phase := 0.0
	for i := 0; i < s.nframes; i++ {
		rect := image.Rect(0, 0, 2*s.size+1, 2*s.size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(s.cycles)*2*math.Pi; t += s.res {
			x := math.Sin(t)
			y := math.Sin(t*s.freq + phase)
			xi := s.size + int(x*float64(s.size)+0.5)
			yi := s.size + int(y*float64(s.size)+0.5)
			img.SetColorIndex(xi, yi, greenIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, s.delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
