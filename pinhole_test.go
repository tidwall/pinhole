package pinhole

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/image/draw"
)

func TestPinhole(t *testing.T) {
	start := time.Now()
	var ai int
	for a := 0.0; a < math.Pi*2; a += math.Pi * 2 / 60 {

		p := New()
		p.Begin()
		p.DrawLine(-0.2, -0.2, -0.2, -0.2, -0.2, 0.2)
		p.DrawLine(0.2, -0.2, -0.2, 0.2, -0.2, 0.2)
		p.DrawLine(0.2, 0.2, -0.2, 0.2, 0.2, 0.2)
		p.DrawLine(-0.2, 0.2, -0.2, -0.2, 0.2, 0.2)

		p.DrawLine(-0.2, -0.2, -0.2, 0.2, -0.2, -0.2)
		p.DrawLine(0.2, -0.2, -0.2, 0.2, 0.2, -0.2)
		p.DrawLine(0.2, 0.2, -0.2, -0.2, 0.2, -0.2)
		p.DrawLine(-0.2, 0.2, -0.2, -0.2, -0.2, -0.2)

		p.DrawLine(-0.2, -0.2, 0.2, 0.2, -0.2, 0.2)
		p.DrawLine(0.2, -0.2, 0.2, 0.2, 0.2, 0.2)
		p.DrawLine(0.2, 0.2, 0.2, -0.2, 0.2, 0.2)
		p.DrawLine(-0.2, 0.2, 0.2, -0.2, -0.2, 0.2)
		p.Rotate(a/4, a, 0)
		p.End()

		p.Begin()
		p.DrawCircle(-0.4, 0.4, 0, 0.2)
		p.Rotate(0, a, 0)
		p.End()

		p.Begin()
		p.DrawCircle(0, 0, 0, 0.1)
		p.Rotate(0, a, 0)
		p.Translate(0.6, 0, 0)
		p.End()

		/*
			for i := float64(0); i < steps; i++ {
				p.PushRotation(math.Pi/steps*i, 0, 0)
				p.DrawCircle(0, 0, 0, 0.4)
				p.PopRotation()
			}
		*/
		var opts ImageOptions
		opts.BGColor = color.White
		opts.FGColor = color.RGBA{255, 0, 0, 255}
		opts.LineWidth = 1
		opts.Scale = 2
		err := p.SavePNG(fmt.Sprintf("out%d.png", ai), 500, 500, &opts)
		if err != nil {
			t.Fatal(err)
		}
		ai++
		break
	}
	println("done", time.Since(start).String())
	return
	// load static image and construct outGif
	var palette = []color.Color{}
	for i := uint8(0); i < 255; i++ {
		palette = append(palette, color.RGBA{0xff / 255 * i, 0xff / 255 * i, 0xff / 255 * i, 0xff})
	}
	palette = append(palette, color.RGBA{0xff, 0xff, 0xff, 0xff})
	outGif := &gif.GIF{}
	for i := 0; i < ai; i++ {
		f, err := os.Open(fmt.Sprintf("out%d.png", i))
		if err != nil {
			t.Fatal(err)
		}
		inPng, err := png.Decode(f)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		inGif := image.NewPaletted(inPng.Bounds(), palette)
		draw.Draw(inGif, inPng.Bounds(), inPng, image.Point{}, draw.Src)
		outGif.Image = append(outGif.Image, inGif)
		outGif.Delay = append(outGif.Delay, 0)
	}

	// save to out.gif
	f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}

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

func TestGopher(t *testing.T) {
	start := time.Now()
	obj, err := loadObj("gopher.obj")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("read object %s\n", time.Since(start))
	start = time.Now()
	p := New()
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
	fmt.Printf("drawn faces %s\n", time.Since(start))
	start = time.Now()
	p.Rotate(0, math.Pi/2, 0)
	fmt.Printf("rotated %s\n", time.Since(start))
	start = time.Now()
	p.Translate(0, 0.2, 4)
	fmt.Printf("translated %s\n", time.Since(start))

	start = time.Now()
	opts := *DefaultImageOptions
	opts.LineWidth = 0.01
	opts.NoCaps = true

	//opts.Straight = true
	if err := p.SavePNG("out.png", 500, 500, &opts); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("imaged %s\n", time.Since(start))
	/*
		data, err := ioutil.ReadFile("gopher.obj")
		if err != nil {
			t.Fatal(err)
		}
		p := New()
		var fx, fy, fz float64
		var lx, ly, lz float64
		var i int
		for _, line := range bytes.Split(data, []byte{'\n'}) {
			if bytes.HasPrefix(line, []byte{'v', ' '}) {
				parts := bytes.Split(line[2:], []byte{' '})
				if len(parts) >= 3 {
					x, err1 := strconv.ParseFloat(string(parts[0]), 64)
					y, err2 := strconv.ParseFloat(string(parts[1]), 64)
					z, err3 := strconv.ParseFloat(string(parts[2]), 64)
					if err1 == nil && err2 == nil && err3 == nil {
						if i == 3 {
							i = 0
						}
						if i == 0 {
							fx, fy, fz = x, y, z
						} else {
							p.DrawLine(lx, ly, lz, x, y, z)
						}
						lx, ly, lz = x, y, z
						i++
					}
				}
			}
		}
		if i > 1 {
			p.DrawLine(lx, ly, lz, fx, fy, fz)
		}
		p.Rotate(0, math.Pi/2, 0)
		p.Translate(0, 0, 4)
		opts := *DefaultImageOptions
		opts.LineWidth = 0.01

		if err := p.SavePNG("out.png", 500, 500, &opts); err != nil {
			t.Fatal(err)
		}
	*/
}
