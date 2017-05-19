package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"math"
	"os"

	"github.com/tidwall/pinhole"
)

func main() {
	p := pinhole.New()
	f, err := os.Open("suzanne.obj")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err = p.LoadObj(f); err != nil {
		log.Fatal(err)
	}

	// suzanne is a little offscreen, scale and center it
	p.Scale(0.50, 0.50, 0.50)
	p.Center()
	opts := *pinhole.DefaultImageOptions
	opts.LineWidth = 0.3 // thin lines
	var n = 60
	var i int
	var images []image.Image
	var step = math.Pi * 2 / float64(n)
	for a := 0.0; a < math.Pi*2; a += step {
		p.Rotate(0, step, 0)
		fmt.Printf("frame %d/%d, %f\n", i, n, a)
		if i == 0 {
			p.SavePNG("suzanne.png", 750, 750, &opts)
		}
		img := p.Image(750, 750, &opts)
		images = append(images, img)
		i++
	}
	fmt.Printf("encoding GIF\n")
	// load static image and construct outGif
	var palette = []color.Color{}
	colors := uint8(8)
	for i := uint8(0); i < colors; i++ {
		palette = append(palette, color.RGBA{0xff / colors * i, 0xff / colors * i, 0xff / colors * i, 0xff})
	}
	palette = append(palette, color.RGBA{0xff, 0xff, 0xff, 0xff})
	outGif := &gif.GIF{}
	for i := 0; i < len(images); i++ {
		inPng := images[i]
		inGif := image.NewPaletted(inPng.Bounds(), palette)
		draw.Draw(inGif, inPng.Bounds(), inPng, image.Point{}, draw.Src)
		outGif.Image = append(outGif.Image, inGif)
		outGif.Delay = append(outGif.Delay, 0)
	}
	f, _ = os.OpenFile("suzanne.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}
