package main

import (
	"image/color"
	"math"

	"github.com/tidwall/pinhole"
)

func main() {
	p := pinhole.New()

	p.DrawLine(-0.3, -0.3, -0.3, -0.3, -0.3, 0.3)
	p.DrawLine(0.3, -0.3, -0.3, 0.3, -0.3, 0.3)
	p.DrawLine(0.3, 0.3, -0.3, 0.3, 0.3, 0.3)
	p.DrawLine(-0.3, 0.3, -0.3, -0.3, 0.3, 0.3)
	p.DrawLine(-0.3, -0.3, -0.3, 0.3, -0.3, -0.3)
	p.DrawLine(0.3, -0.3, -0.3, 0.3, 0.3, -0.3)
	p.DrawLine(0.3, 0.3, -0.3, -0.3, 0.3, -0.3)
	p.DrawLine(-0.3, 0.3, -0.3, -0.3, -0.3, -0.3)
	p.DrawLine(-0.3, -0.3, 0.3, 0.3, -0.3, 0.3)
	p.DrawLine(0.3, -0.3, 0.3, 0.3, 0.3, 0.3)
	p.DrawLine(0.3, 0.3, 0.3, -0.3, 0.3, 0.3)
	p.DrawLine(-0.3, 0.3, 0.3, -0.3, -0.3, 0.3)
	p.Rotate(math.Pi/3, math.Pi/6, 0)

	p.Begin()
	p.DrawCircle(0, 0, 0, 0.2)
	p.Rotate(0, math.Pi/2, 0)
	p.Translate(-0.6, -0.4, 0)
	p.Colorize(color.RGBA{255, 0, 0, 255})
	p.End()
	p.SavePNG("cube.png", 500, 500, nil)
}
