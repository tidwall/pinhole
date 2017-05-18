package pinhole

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/fogleman/gg"
	"github.com/google/btree"
)

type line struct {
	x1, y1, z1 float64
	x2, y2, z2 float64
	nocaps     bool
}

func (l *line) Center() []float64 {
	min, max := l.Rect()
	return []float64{
		(max[0] + min[0]) / 2,
		(max[1] + min[1]) / 2,
		(max[2] + min[2]) / 2,
	}
}
func (l *line) Rect() (min, max []float64) {
	min, max = minMax(l.x1, l.y1, l.z1, l.x2, l.y2, l.z2)
	return
}

type Pinhole struct {
	lines  []*line
	stack  []int
	nocaps bool
	dirty  bool
}

func New() *Pinhole {
	return &Pinhole{}
}

func (p *Pinhole) Begin() {
	p.stack = append(p.stack, len(p.lines))
}
func (p *Pinhole) End() {
	if len(p.stack) > 0 {
		p.stack = p.stack[:len(p.stack)-1]
	}
}

func (p *Pinhole) Rotate(x, y, z float64) {
	var i int
	if len(p.stack) > 0 {
		i = p.stack[len(p.stack)-1]
	}
	for ; i < len(p.lines); i++ {
		l := p.lines[i]
		if x != 0 {
			l.x1, l.y1, l.z1 = rotate(l.x1, l.y1, l.z1, x, 0)
			l.x2, l.y2, l.z2 = rotate(l.x2, l.y2, l.z2, x, 0)
		}
		if y != 0 {
			l.x1, l.y1, l.z1 = rotate(l.x1, l.y1, l.z1, y, 1)
			l.x2, l.y2, l.z2 = rotate(l.x2, l.y2, l.z2, y, 1)
		}
		if z != 0 {
			l.x1, l.y1, l.z1 = rotate(l.x1, l.y1, l.z1, z, 2)
			l.x2, l.y2, l.z2 = rotate(l.x2, l.y2, l.z2, z, 2)
		}
		p.lines[i] = l
	}
}

func (p *Pinhole) Translate(x, y, z float64) {
	var i int
	if len(p.stack) > 0 {
		i = p.stack[len(p.stack)-1]
	}
	for ; i < len(p.lines); i++ {
		p.lines[i].x1 += x
		p.lines[i].y1 += y
		p.lines[i].z1 += z
		p.lines[i].x2 += x
		p.lines[i].y2 += y
		p.lines[i].z2 += z
	}
}

func (p *Pinhole) Scale(x, y, z float64) {
	var i int
	if len(p.stack) > 0 {
		i = p.stack[len(p.stack)-1]
	}
	for ; i < len(p.lines); i++ {
		p.lines[i].x1 *= x
		p.lines[i].y1 *= y
		p.lines[i].z1 *= z
		p.lines[i].x2 *= x
		p.lines[i].y2 *= y
		p.lines[i].z2 *= z
	}
}

func (p *Pinhole) Center() {
	var i int
	if len(p.stack) > 0 {
		i = p.stack[len(p.stack)-1]
	}
	minx, miny, minz := math.Inf(+1), math.Inf(+1), math.Inf(+1)
	maxx, maxy, maxz := math.Inf(-1), math.Inf(-1), math.Inf(-1)
	for ; i < len(p.lines); i++ {
		if p.lines[i].x1 < minx {
			minx = p.lines[i].x1
		}
		if p.lines[i].x1 > maxx {
			maxx = p.lines[i].x1
		}
		if p.lines[i].y1 < miny {
			miny = p.lines[i].y1
		}
		if p.lines[i].y1 > maxy {
			maxy = p.lines[i].y1
		}
		if p.lines[i].z1 < minz {
			minz = p.lines[i].z1
		}
		if p.lines[i].z2 > maxz {
			maxz = p.lines[i].z2
		}
		if p.lines[i].x2 < minx {
			minx = p.lines[i].x2
		}
		if p.lines[i].x2 > maxx {
			maxx = p.lines[i].x2
		}
		if p.lines[i].y2 < miny {
			miny = p.lines[i].y2
		}
		if p.lines[i].y2 > maxy {
			maxy = p.lines[i].y2
		}
		if p.lines[i].z2 < minz {
			minz = p.lines[i].z2
		}
		if p.lines[i].z2 > maxz {
			maxz = p.lines[i].z2
		}
	}
	x := (maxx + minx) / 2
	y := (maxy + miny) / 2
	z := (maxz + minz) / 2
	p.Translate(-x, -y, -z)
}

func minMax(x1, y1, z1, x2, y2, z2 float64) (min, max []float64) {
	min = []float64{x1, y1, z1}
	max = []float64{x2, y2, z2}
	for i := 0; i < 3; i++ {
		if min[i] > max[i] {
			min[i], max[i] = max[i], min[i]
		}
	}
	return
}

func (p *Pinhole) DrawLine(x1, y1, z1, x2, y2, z2 float64) {
	l := &line{
		x1: x1, y1: y1, z1: z1,
		x2: x2, y2: y2, z2: z2,
		nocaps: p.nocaps,
	}
	p.lines = append(p.lines, l)
}

func (p *Pinhole) DrawCircle(x, y, z float64, radius float64) {
	p.nocaps = true
	var fx, fy, fz float64
	var lx, ly, lz float64
	steps := 180.0
	for i := float64(0); i < steps; i++ {
		var dx, dy, dz float64
		dx, dy = destination(x, y, (math.Pi*2)/steps*i, radius)
		dz = z
		if i > 0 {
			p.DrawLine(lx, ly, lz, dx, dy, dz)
		} else {
			fx, fy, fz = dx, dy, dz
		}
		lx, ly, lz = dx, dy, dz
	}
	p.DrawLine(lx, ly, lz, fx, fy, fz)
	p.nocaps = false
}

type ImageOptions struct {
	FGColor   color.Color
	BGColor   color.Color
	LineWidth float64
	Scale     float64
}

var DefaultImageOptions = &ImageOptions{
	FGColor:   color.Black,
	BGColor:   color.Transparent,
	LineWidth: 1,
	Scale:     1,
}

func (p *Pinhole) Image(width, height int, opts *ImageOptions) *image.RGBA {
	if opts == nil {
		opts = DefaultImageOptions
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c := gg.NewContextForRGBA(img)
	if opts.BGColor != nil {
		c.SetColor(opts.BGColor)
		c.DrawRectangle(0, 0, float64(width), float64(height))
		c.Fill()
	}
	caps := newCapTree()
	if opts.FGColor != nil {
		c.SetColor(opts.FGColor)
	} else {
		c.SetRGB(0, 0, 0)
	}
	fwidth, fheight := float64(width), float64(height)
	focal := math.Min(fwidth, fheight) / 2
	for _, line := range p.lines {
		x1, y1, z1 := line.x1, line.y1, line.z1
		x2, y2, z2 := line.x2, line.y2, line.z2
		px1, py1 := projectPoint(x1, y1, z1, fwidth, fheight, focal, opts.Scale)
		px2, py2 := projectPoint(x2, y2, z2, fwidth, fheight, focal, opts.Scale)
		if !onscreen(fwidth, fheight, px1, py1, px2, py2) {
			continue
		}
		t1 := lineWidthAtZ(z1, focal) * opts.LineWidth
		t2 := lineWidthAtZ(z2, focal) * opts.LineWidth
		var cap1, cap2 bool
		if !line.nocaps {
			cap1 = caps.insert(x1, y1, z1)
			cap2 = caps.insert(x2, y2, z2)
		}
		drawUnbalancedLineSegment(c, px1, py1, px2, py2, t1, t2, cap1, cap2)
	}
	c.Fill()
	return img
}
func onscreen(w, h float64, x1, y1, x2, y2 float64) bool {
	amin := [2]float64{0, 0}
	amax := [2]float64{w, h}
	var bmin [2]float64
	var bmax [2]float64
	if x1 < x2 {
		bmin[0], bmax[0] = x1, x2
	} else {
		bmin[0], bmax[0] = x2, x1
	}
	if y1 < y2 {
		bmin[1], bmax[1] = y1, y2
	} else {
		bmin[1], bmax[1] = y2, y1
	}
	for i := 0; i < len(amin); i++ {
		if !(bmin[i] <= amax[i] && bmax[i] >= amin[i]) {
			return false
		}
	}
	return true
}

func (p *Pinhole) SavePNG(path string, width, height int, opts *ImageOptions) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, p.Image(width, height, opts))
}

// projectPoint projects a 3d point cartesian point to 2d screen coords.
//     Origin is the center
//     X is left/right
//     Y is down/up
//     Z is near/far, the 0 position is focal distance away from lens.
func projectPoint(
	x, y, z float64, // 3d point to project
	w, h, f float64, // width, height, focal
	scale float64, // scale
) (px, py float64) { // projected point
	x, y, z = x*scale*f, y*scale*f, z*scale*f
	zz := z + f
	if zz == 0 {
		zz = math.SmallestNonzeroFloat64
	}
	px = x*(f/zz) + w/2
	py = y*(f/zz) - h/2
	py *= -1
	return
}

func lineWidthAtZ(z float64, f float64) float64 {
	return ((z*-1 + 1) / 2) * f * 0.04
}

func drawUnbalancedLineSegment(c *gg.Context,
	x1, y1, x2, y2 float64,
	t1, t2 float64,
	cap1, cap2 bool,
) {
	if cap1 && t1 < 2 {
		cap1 = false
	}
	if cap2 && t2 < 2 {
		cap2 = false
	}
	a := lineAngle(x1, y1, x2, y2)
	dx1, dy1 := destination(x1, y1, a-math.Pi/2, t1/2)
	dx2, dy2 := destination(x1, y1, a+math.Pi/2, t1/2)
	dx3, dy3 := destination(x2, y2, a+math.Pi/2, t2/2)
	dx4, dy4 := destination(x2, y2, a-math.Pi/2, t2/2)
	c.MoveTo(dx1, dy1)
	if cap1 {
		ax1, ay1 := destination(dx1, dy1, a-math.Pi*2, t1*0.552284749831)
		ax2, ay2 := destination(dx2, dy2, a-math.Pi*2, t1*0.552284749831)
		c.CubicTo(ax1, ay1, ax2, ay2, dx2, dy2)
	} else {
		c.LineTo(dx2, dy2)
	}
	c.LineTo(dx3, dy3)
	if cap2 {
		ax1, ay1 := destination(dx3, dy3, a-math.Pi*2, -t2*0.552284749831)
		ax2, ay2 := destination(dx4, dy4, a-math.Pi*2, -t2*0.552284749831)
		c.CubicTo(ax1, ay1, ax2, ay2, dx4, dy4)
	} else {
		c.LineTo(dx4, dy4)
	}
	c.LineTo(dx1, dy1)
	c.ClosePath()
	a = a*180/math.Pi - 90
	if a < 0 {
		a += 360
	}
}

func lineAngle(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y1-y2, x1-x2)
}

func destination(x, y, angle, distance float64) (dx, dy float64) {
	dx = x + math.Cos(angle)*distance
	dy = y + math.Sin(angle)*distance
	return
}

// https://www.siggraph.org/education/materials/HyperGraph/modeling/mod_tran/3drota.htm
func rotate(x, y, z float64, q float64, which int) (dx, dy, dz float64) {
	switch which {
	case 0: // x
		dy = y*math.Cos(q) - z*math.Sin(q)
		dz = y*math.Sin(q) + z*math.Cos(q)
		dx = x
	case 1: // y
		dz = z*math.Cos(q) - x*math.Sin(q)
		dx = z*math.Sin(q) + x*math.Cos(q)
		dy = y
	case 2: // z
		dx = x*math.Cos(q) - y*math.Sin(q)
		dy = x*math.Sin(q) + y*math.Cos(q)
		dz = z
	}
	return
}

type capItem struct {
	point [3]float64
}

func (a *capItem) Less(v btree.Item) bool {
	b := v.(*capItem)
	for i := 2; i >= 0; i-- {
		if a.point[i] < b.point[i] {
			return true
		}
		if a.point[i] > b.point[i] {
			return false
		}
	}
	return false
}

// really lazy structure.
type capTree struct {
	tr *btree.BTree
}

func newCapTree() *capTree {
	return &capTree{
		tr: btree.New(9),
	}
}

func (tr *capTree) insert(x, y, z float64) bool {
	if tr.has(x, y, z) {
		return false
	}
	tr.tr.ReplaceOrInsert(&capItem{point: [3]float64{x, y, z}})
	return true
}

func (tr *capTree) has(x, y, z float64) bool {
	return tr.tr.Has(&capItem{point: [3]float64{x, y, z}})
}
