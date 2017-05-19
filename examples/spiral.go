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

	"github.com/fogleman/ease"
	"github.com/tidwall/pinhole"
)

func makeSpiral() *pinhole.Pinhole {
	p := pinhole.New()
	n := 360.0
	for i, z := 0.0, -0.2; i < n && z <= 1; i, z = i+1, z+0.003 {
		d := 0.5 * (1 - (i / n / 2))  // distance of circle from origin
		a := math.Pi * 2 / 30 * i     // angle of circle from origin
		r := 0.03 * (1 - (i / n / 2)) // radius of circle
		p.DrawCircle(math.Cos(a)*d, math.Sin(a)*d, z, r)
	}
	return p
}
func main() {
	p := makeSpiral()
	var imgs []image.Image
	n := 60
	rotate := math.Pi / 3
	for i := 0; i < n; i++ {
		fmt.Printf("frame %d/%d\n", i, n)
		t := float64(i) / float64(n)
		if t < 0.5 {
			t = ease.InSine(t * 2)
		} else {
			t = 1 - ease.OutSine((t-0.5)*2)
		}
		a := rotate * t
		p.Rotate(a, 0, 0)
		if i == 0 {
			if err := p.SavePNG("spiral.png", 750, 750, nil); err != nil {
				log.Fatal(err)
			}
		}
		imgs = append(imgs, p.Image(750, 750, nil))
		p.Rotate(-a, 0, 0)
	}
	fmt.Printf("encoding GIF\n")
	if err := encodeGIF(imgs, "spiral.gif"); err != nil {
		log.Fatal(err)
	}
}

func encodeGIF(imgs []image.Image, path string) error {
	// load static image and construct outGif
	var palette = []color.Color{}
	colors := uint8(8)
	for i := uint8(0); i < colors; i++ {
		palette = append(palette, color.RGBA{0xff / colors * i, 0xff / colors * i, 0xff / colors * i, 0xff})
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
