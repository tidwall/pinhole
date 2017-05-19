package pinhole

import (
	"fmt"
	"image/color"
	"math"
	"testing"
)

func TestPinhole(t *testing.T) {

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
	//p.Rotate(, a, 0)
	p.Colorize(color.RGBA{255, 0, 0, 255})
	p.End()

	p.Begin()
	p.DrawCircle(-0.4, 0.4, 0, 0.2)
	//p.Rotate(0, a, 0)
	p.End()

	p.Begin()
	p.DrawCircle(0, 0, 0, 0.1)
	//p.Rotate(0, a, 0)
	p.Translate(0.6, 0, 0)
	p.End()

	p.DrawCircle(0, 0.1, 0, 0.3)

	p.Begin()
	p.DrawDot(0, 0.1, -math.SmallestNonzeroFloat64, 0.01)
	p.Colorize(color.RGBA{255, 255, 0, 255})
	p.End()
	/*
		for i := float64(0); i < steps; i++ {
			p.PushRotation(math.Pi/steps*i, 0, 0)
			p.DrawCircle(0, 0, 0, 0.4)
			p.PopRotation()
		}
	*/
	p.Rotate(math.Pi/3, 0, 0)
	opts := *DefaultImageOptions
	/*
		opts.BGColor = color.White
		opts.FGColor = color.RGBA{255, 0, 0, 255}
		opts.LineWidth = 1
		opts.Scale = 2
	*/
	err := p.SavePNG(fmt.Sprintf("out.png"), 1000, 1000, &opts)
	if err != nil {
		t.Fatal(err)
	}
}
