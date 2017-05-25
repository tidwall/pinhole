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
	var imgs []image.Image
	n := 60
	//rotate := math.Pi / 3
	for i := 0; i < n; i++ {

		fmt.Printf("frame %d/%d\n", i, n)

		p := pinhole.New()
		p.Begin()
		p.DrawCube(-0.2, -0.2, -0.2, 0.2, 0.2, 0.2)
		p.Rotate(0, math.Pi*2/(float64(n)/float64(i)), 0)
		p.Colorize(color.RGBA{255, 0, 0, 255})
		p.End()

		p.Begin()
		p.DrawCircle(0, 0, 0, 0.2)
		p.Rotate(math.Pi*2/(float64(n)/float64(i)), math.Pi*4/(float64(n)/float64(i)), 0)
		p.End()

		p.Begin()
		p.DrawCircle(0, 0, 0, 0.2)
		p.Rotate(-math.Pi*2/(float64(n)/float64(i)), math.Pi*4/(float64(n)/float64(i)), 0)
		p.End()

		p.Scale(1.75, 1.75, 1.75)

		if i == 0 {
			if err := p.SavePNG("shapes.png", 750, 750, nil); err != nil {
				log.Fatal(err)
			}
			//return
		}
		imgs = append(imgs, p.Image(750, 750, nil))
	}
	fmt.Printf("encoding GIF\n")
	if err := encodeGIF(imgs, "shapes.gif"); err != nil {
		log.Fatal(err)
	}
}

func encodeGIF(imgs []image.Image, path string) error {
	// load static image and construct outGif
	var palette = []color.Color{}
	colors := uint8(16)
	for i := uint8(0); i < colors; i++ {
		palette = append(palette, color.RGBA{0xff / colors * i, 0xff / colors * i, 0xff / colors * i, 0xff})
	}
	for i := uint8(0); i < colors; i++ {
		palette = append(palette, color.RGBA{0xff / colors * i, 0, 0, 0xff})
	}
	for i := uint8(0); i < colors; i++ {
		palette = append(palette, color.RGBA{0, 0, 0xff / colors * i, 0xff})
	}
	palette = append(palette, color.RGBA{0xff, 0xff, 0xff, 0xff})
	outGif := &gif.GIF{}
	for i := 0; i < len(imgs); i++ {
		inPng := imgs[i]
		inGif := image.NewPaletted(inPng.Bounds(), palette)
		draw.Draw(inGif, inPng.Bounds(), inPng, image.Point{}, draw.Src)
		outGif.Image = append(outGif.Image, inGif)
		outGif.Delay = append(outGif.Delay, 0)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, outGif)
}
