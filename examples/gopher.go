package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/pinhole"
)

type object struct {
	faces [][][3]float64
}

func loadObj(path string) (*object, error) {
	obj := &object{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var verts [][3]float64
	for ln, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "v ") {
			parts := strings.Split(line[2:], " ")
			if len(parts) >= 3 {
				verts = append(verts, [3]float64{})
				for j := 0; j < 3; j++ {
					if verts[len(verts)-1][j], err = strconv.ParseFloat(parts[j], 64); err != nil {
						return nil, fmt.Errorf("line %d: %s", ln+1, err.Error())
					}
				}
			}
		} else if strings.HasPrefix(line, "f ") {
			parts := strings.Split(line[2:], " ")
			if len(parts) >= 3 {
				obj.faces = append(obj.faces, [][3]float64{})
				for _, part := range parts {
					part = strings.Split(part, "/")[0]
					idx, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("line %d: %s", ln+1, err.Error())
					}
					if int(idx) > len(verts) {
						return nil, fmt.Errorf("line %d: invalid vert index: %d", ln+1, idx)
					}
					obj.faces[len(obj.faces)-1] = append(obj.faces[len(obj.faces)-1], verts[idx-1])
				}
			}
		}
	}
	return obj, nil
}
func main() {
	genGIF()
	genPNG()
}
func genPNG() {
	obj, err := loadObj("gopher.obj")
	if err != nil {
		log.Fatal(err)
	}
	p := pinhole.New()
	for _, faces := range obj.faces {
		var fx, fy, fz float64
		var lx, ly, lz float64
		var i int
		for _, face := range faces {
			if i == 0 {
				fx, fy, fz = face[0], face[1], face[2]
			} else {
				p.DrawLine(lx, ly, lz, face[0], face[1], face[2])
			}
			lx, ly, lz = face[0], face[1], face[2]
			i++
		}
		if i > 1 {
			p.DrawLine(lx, ly, lz, fx, fy, fz)
		}
	}
	opts := *pinhole.DefaultImageOptions
	opts.BGColor = color.White
	opts.LineWidth = 0.02 // thin lines
	// gopher is a little offscreen, scale, rotate, and center it
	p.Scale(0.25, 0.25, 0.25)
	p.Center()
	p.Rotate(0, math.Pi/2, 0)
	fmt.Printf("encoding PNG\n")
	p.SavePNG("gopher.png", 500, 500, &opts)

}
func genGIF() {
	obj, err := loadObj("gopher.obj")
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	n := 60
	p := pinhole.New()
	for _, faces := range obj.faces {
		var fx, fy, fz float64
		var lx, ly, lz float64
		var i int
		for _, face := range faces {
			if i == 0 {
				fx, fy, fz = face[0], face[1], face[2]
			} else {
				p.DrawLine(lx, ly, lz, face[0], face[1], face[2])
			}
			lx, ly, lz = face[0], face[1], face[2]
			i++
		}
		if i > 1 {
			p.DrawLine(lx, ly, lz, fx, fy, fz)
		}
	}
	opts := *pinhole.DefaultImageOptions
	opts.BGColor = color.White
	opts.LineWidth = 0.02 // thin lines
	var images []image.Image
	step := math.Pi * 2 / float64(n)
	for a := 0.0; a < math.Pi*2; a += step {
		// gopher is a little offscreen, scale and center it
		p.Scale(0.25, 0.25, 0.25)
		p.Center()
		p.Rotate(0, step, 0)
		fmt.Printf("frame %d/%d, %f\n", i, n, a)
		img := p.Image(500, 500, &opts)
		images = append(images, img)
		i++
	}
	fmt.Printf("encoding GIF\n")
	// load static image and construct outGif
	var palette = []color.Color{}
	for i := uint8(0); i < 255; i++ {
		palette = append(palette, color.RGBA{0xff / 255 * i, 0xff / 255 * i, 0xff / 255 * i, 0xff})
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
	f, _ := os.OpenFile("gopher.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}
