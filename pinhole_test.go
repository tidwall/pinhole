package pinhole

//import (
//	"fmt"
//	"image"
//	"image/color"
//	"image/gif"
//	"image/png"
//	"math"
//	"os"
//	"testing"
//	"time"
//
//	"golang.org/x/image/draw"
//)
//
//func TestPinhole(t *testing.T) {
//	start := time.Now()
//	var ai int
//	for a := 0.0; a < math.Pi*2; a += math.Pi * 2 / 60 {
//
//		p := New()
//
//		p.DrawLine(-0.6, -0.6, 0, -0.6, 0.6, 0)
//
//		p.Begin()
//		p.DrawLine(-0.2, -0.2, -0.2, -0.2, -0.2, 0.2)
//		p.DrawLine(0.2, -0.2, -0.2, 0.2, -0.2, 0.2)
//		p.DrawLine(0.2, 0.2, -0.2, 0.2, 0.2, 0.2)
//		p.DrawLine(-0.2, 0.2, -0.2, -0.2, 0.2, 0.2)
//
//		p.DrawLine(-0.2, -0.2, -0.2, 0.2, -0.2, -0.2)
//		p.DrawLine(0.2, -0.2, -0.2, 0.2, 0.2, -0.2)
//		p.DrawLine(0.2, 0.2, -0.2, -0.2, 0.2, -0.2)
//		p.DrawLine(-0.2, 0.2, -0.2, -0.2, -0.2, -0.2)
//
//		p.DrawLine(-0.2, -0.2, 0.2, 0.2, -0.2, 0.2)
//		p.DrawLine(0.2, -0.2, 0.2, 0.2, 0.2, 0.2)
//		p.DrawLine(0.2, 0.2, 0.2, -0.2, 0.2, 0.2)
//		p.DrawLine(-0.2, 0.2, 0.2, -0.2, -0.2, 0.2)
//		p.Rotate(a/4, a, 0)
//		p.End()
//
//		p.Begin()
//		p.DrawCircle(-0.4, 0.4, 0, 0.2)
//		p.Rotate(0, a, 0)
//		p.End()
//
//		p.Begin()
//		p.DrawCircle(0, 0, 0, 0.1)
//		p.Rotate(0, a, 0)
//		p.Translate(0.6, 0, 0)
//		p.End()
//
//		/*
//			for i := float64(0); i < steps; i++ {
//				p.PushRotation(math.Pi/steps*i, 0, 0)
//				p.DrawCircle(0, 0, 0, 0.4)
//				p.PopRotation()
//			}
//		*/
//		//p.Rotate(math.Pi/3, 0, 0)
//		opts := *DefaultImageOptions
//		/*
//			opts.BGColor = color.White
//			opts.FGColor = color.RGBA{255, 0, 0, 255}
//			opts.LineWidth = 1
//			opts.Scale = 2
//		*/
//		err := p.SavePNG(fmt.Sprintf("out%d.png", ai), 500, 500, &opts)
//		if err != nil {
//			t.Fatal(err)
//		}
//		ai++
//		break
//	}
//	println("done", time.Since(start).String())
//	return
//	// load static image and construct outGif
//	var palette = []color.Color{}
//	for i := uint8(0); i < 255; i++ {
//		palette = append(palette, color.RGBA{0xff / 255 * i, 0xff / 255 * i, 0xff / 255 * i, 0xff})
//	}
//	palette = append(palette, color.RGBA{0xff, 0xff, 0xff, 0xff})
//	outGif := &gif.GIF{}
//	for i := 0; i < ai; i++ {
//		f, err := os.Open(fmt.Sprintf("out%d.png", i))
//		if err != nil {
//			t.Fatal(err)
//		}
//		inPng, err := png.Decode(f)
//		if err != nil {
//			t.Fatal(err)
//		}
//		f.Close()
//		inGif := image.NewPaletted(inPng.Bounds(), palette)
//		draw.Draw(inGif, inPng.Bounds(), inPng, image.Point{}, draw.Src)
//		outGif.Image = append(outGif.Image, inGif)
//		outGif.Delay = append(outGif.Delay, 0)
//	}
//
//	// save to out.gif
//	f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
//	defer f.Close()
//	gif.EncodeAll(f, outGif)
//}
