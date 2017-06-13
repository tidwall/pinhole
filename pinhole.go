package pinhole

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/google/btree"
)

const circleSteps = 45

var gof = func() *truetype.Font {
	gof, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	return gof
}()

type line struct {
	x1, y1, z1 float64
	x2, y2, z2 float64
	nocaps     bool
	color      color.Color
	str        string
	scale      float64
	circle     bool
	cfirst     *line
	cprev      *line
	cnext      *line

	drawcoords *fourcorners
}

func (l *line) Rect() (min, max [3]float64) {
	if l.x1 < l.x2 {
		min[0], max[0] = l.x1, l.x2
	} else {
		min[0], max[0] = l.x2, l.x1
	}
	if l.y1 < l.y2 {
		min[1], max[1] = l.y1, l.y2
	} else {
		min[1], max[1] = l.y2, l.y1
	}
	if l.z1 < l.z2 {
		min[2], max[2] = l.z1, l.z2
	} else {
		min[2], max[2] = l.z2, l.z1
	}
	return
}

func (l *line) Center() []float64 {
	min, max := l.Rect()
	return []float64{
		(max[0] + min[0]) / 2,
		(max[1] + min[1]) / 2,
		(max[2] + min[2]) / 2,
	}
}

type Pinhole struct {
	lines []*line
	stack []int
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
		if len(p.lines[i].str) > 0 {
			p.lines[i].scale *= math.Min(x, y)
		}
	}
}

func (p *Pinhole) Colorize(color color.Color) {
	var i int
	if len(p.stack) > 0 {
		i = p.stack[len(p.stack)-1]
	}
	for ; i < len(p.lines); i++ {
		p.lines[i].color = color
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

func (p *Pinhole) DrawString(x, y, z float64, s string) {
	if s != "" {
		p.DrawLine(x, y, z, x, y, z)
		//p.lines[len(p.lines)-1].scale = 10 / 0.1 * radius
		p.lines[len(p.lines)-1].str = s
	}
}
func (p *Pinhole) DrawRect(minx, miny, maxx, maxy, z float64) {
	p.DrawLine(minx, maxy, z, maxx, maxy, z)
	p.DrawLine(maxx, maxy, z, maxx, miny, z)
	p.DrawLine(maxx, miny, z, minx, miny, z)
	p.DrawLine(minx, miny, z, minx, maxy, z)
}
func (p *Pinhole) DrawCube(minx, miny, minz, maxx, maxy, maxz float64) {
	p.DrawLine(minx, maxy, minz, maxx, maxy, minz)
	p.DrawLine(maxx, maxy, minz, maxx, miny, minz)
	p.DrawLine(maxx, miny, minz, minx, miny, minz)
	p.DrawLine(minx, miny, minz, minx, maxy, minz)
	p.DrawLine(minx, maxy, maxz, maxx, maxy, maxz)
	p.DrawLine(maxx, maxy, maxz, maxx, miny, maxz)
	p.DrawLine(maxx, miny, maxz, minx, miny, maxz)
	p.DrawLine(minx, miny, maxz, minx, maxy, maxz)
	p.DrawLine(minx, maxy, minz, minx, maxy, maxz)
	p.DrawLine(maxx, maxy, minz, maxx, maxy, maxz)
	p.DrawLine(maxx, miny, minz, maxx, miny, maxz)
	p.DrawLine(minx, miny, minz, minx, miny, maxz)
}

func (p *Pinhole) DrawDot(x, y, z float64, radius float64) {
	p.DrawLine(x, y, z, x, y, z)
	p.lines[len(p.lines)-1].scale = 10 / 0.1 * radius
}

func (p *Pinhole) DrawLine(x1, y1, z1, x2, y2, z2 float64) {
	l := &line{
		x1: x1, y1: y1, z1: z1,
		x2: x2, y2: y2, z2: z2,
		color: color.Black,
		scale: 1,
	}
	p.lines = append(p.lines, l)
}
func (p *Pinhole) DrawCircle(x, y, z float64, radius float64) {
	var fx, fy, fz float64
	var lx, ly, lz float64
	var first, prev *line
	// we go one beyond the steps because we need to join at the end
	for i := float64(0); i <= circleSteps; i++ {
		var dx, dy, dz float64
		dx, dy = destination(x, y, (math.Pi*2)/circleSteps*i, radius)
		dz = z
		if i > 0 {
			if i == circleSteps {
				p.DrawLine(lx, ly, lz, fx, fy, fz)
			} else {
				p.DrawLine(lx, ly, lz, dx, dy, dz)
			}
			line := p.lines[len(p.lines)-1]
			line.nocaps = true
			line.circle = true
			if first == nil {
				first = line
			}
			line.cfirst = first
			line.cprev = prev
			if prev != nil {
				prev.cnext = line
			}
			prev = line

		} else {
			fx, fy, fz = dx, dy, dz
		}
		lx, ly, lz = dx, dy, dz
	}
}

type ImageOptions struct {
	BGColor   color.Color
	LineWidth float64
	Scale     float64
}

var DefaultImageOptions = &ImageOptions{
	BGColor:   color.White,
	LineWidth: 1,
	Scale:     1,
}

type byDistance []*line

func (a byDistance) Len() int {
	return len(a)
}
func (a byDistance) Less(i, j int) bool {
	imin, imax := a[i].Rect()
	jmin, jmax := a[j].Rect()
	for i := 2; i >= 0; i-- {
		if imax[i] > jmax[i] {
			return i == 2
		}
		if imax[i] < jmax[i] {
			return i != 2
		}
		if imin[i] > jmin[i] {
			return i == 2
		}
		if imin[i] < jmin[i] {
			return i != 2
		}
	}
	return false
}
func (a byDistance) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (p *Pinhole) Image(width, height int, opts *ImageOptions) *image.RGBA {
	if opts == nil {
		opts = DefaultImageOptions
	}
	sort.Sort(byDistance(p.lines))
	for _, line := range p.lines {
		line.drawcoords = nil
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c := gg.NewContextForRGBA(img)
	if opts.BGColor != nil {
		c.SetColor(opts.BGColor)
		c.DrawRectangle(0, 0, float64(width), float64(height))
		c.Fill()
	}

	capsMap := make(map[color.Color]*capTree)
	var ccolor color.Color
	var caps *capTree
	fwidth, fheight := float64(width), float64(height)
	focal := math.Min(fwidth, fheight) / 2
	maybeDraw := func(line *line) *fourcorners {
		x1, y1, z1 := line.x1, line.y1, line.z1
		x2, y2, z2 := line.x2, line.y2, line.z2
		px1, py1 := projectPoint(x1, y1, z1, fwidth, fheight, focal, opts.Scale)
		px2, py2 := projectPoint(x2, y2, z2, fwidth, fheight, focal, opts.Scale)
		if !onscreen(fwidth, fheight, px1, py1, px2, py2) && !line.circle && line.str == "" {
			return nil
		}
		t1 := lineWidthAtZ(z1, focal) * opts.LineWidth * line.scale
		t2 := lineWidthAtZ(z2, focal) * opts.LineWidth * line.scale
		if line.str != "" {
			sz := 10 * t1
			c.SetFontFace(truetype.NewFace(gof, &truetype.Options{Size: sz}))
			w, h := c.MeasureString(line.str)
			c.DrawString(line.str, px1-w/2, py1+h*.4)
			return nil
		}
		var cap1, cap2 bool
		if !line.nocaps {
			cap1 = caps.insert(x1, y1, z1)
			cap2 = caps.insert(x2, y2, z2)
		}
		return drawUnbalancedLineSegment(c,
			px1, py1, px2, py2,
			t1, t2,
			cap1, cap2,
			line.circle,
		)
	}
	for _, line := range p.lines {
		if line.color != ccolor {
			ccolor = line.color
			caps = capsMap[ccolor]
			if caps == nil {
				caps = newCapTree()
				capsMap[ccolor] = caps
			}
			c.SetColor(ccolor)
		}
		if line.circle {
			if line.drawcoords == nil {
				// need to process the coords for all segments belonging to
				// the current circle segment.
				// first get the basic estimates
				var coords []*fourcorners
				seg := line.cfirst
				for seg != nil {
					seg.drawcoords = maybeDraw(seg)
					if seg.drawcoords == nil {
						panic("nil!")
					}
					coords = append(coords, seg.drawcoords)
					seg = seg.cnext
				}
				// next reprocess to join the midpoints
				for i := 0; i < len(coords); i++ {
					var line1, line2 *fourcorners
					if i == 0 {
						line1 = coords[len(coords)-1]
					} else {
						line1 = coords[i-1]
					}
					line2 = coords[i]
					midx1 := (line2.x1 + line1.x4) / 2
					midy1 := (line2.y1 + line1.y4) / 2
					midx2 := (line2.x2 + line1.x3) / 2
					midy2 := (line2.y2 + line1.y3) / 2
					line2.x1 = midx1
					line2.y1 = midy1
					line1.x4 = midx1
					line1.y4 = midy1
					line2.x2 = midx2
					line2.y2 = midy2
					line1.x3 = midx2
					line1.y3 = midy2

				}
			}
			// draw the cached coords
			c.MoveTo(line.drawcoords.x1-math.SmallestNonzeroFloat64, line.drawcoords.y1-math.SmallestNonzeroFloat64)
			c.LineTo(line.drawcoords.x2-math.SmallestNonzeroFloat64, line.drawcoords.y2-math.SmallestNonzeroFloat64)
			c.LineTo(line.drawcoords.x3+math.SmallestNonzeroFloat64, line.drawcoords.y3+math.SmallestNonzeroFloat64)
			c.LineTo(line.drawcoords.x4+math.SmallestNonzeroFloat64, line.drawcoords.y4+math.SmallestNonzeroFloat64)
			c.LineTo(line.drawcoords.x1-math.SmallestNonzeroFloat64, line.drawcoords.y1-math.SmallestNonzeroFloat64)
			c.ClosePath()
		} else {
			maybeDraw(line)
		}
		c.Fill()
	}
	return img
}

type fourcorners struct {
	x1, y1, x2, y2, x3, y3, x4, y4 float64
}

func drawUnbalancedLineSegment(c *gg.Context,
	x1, y1, x2, y2 float64,
	t1, t2 float64,
	cap1, cap2 bool,
	circleSegment bool,
) *fourcorners {
	if x1 == x2 && y1 == y2 {
		c.DrawCircle(x1, y1, t1/2)
		return nil
	}

	a := lineAngle(x1, y1, x2, y2)
	dx1, dy1 := destination(x1, y1, a-math.Pi/2, t1/2)
	dx2, dy2 := destination(x1, y1, a+math.Pi/2, t1/2)
	dx3, dy3 := destination(x2, y2, a+math.Pi/2, t2/2)
	dx4, dy4 := destination(x2, y2, a-math.Pi/2, t2/2)
	if circleSegment {
		return &fourcorners{dx1, dy1, dx2, dy2, dx3, dy3, dx4, dy4}
	}
	const cubicCorner = 1.0 / 3 * 2 //0.552284749831
	if cap1 && t1 < 2 {
		cap1 = false
	}
	if cap2 && t2 < 2 {
		cap2 = false
	}
	c.MoveTo(dx1, dy1)
	if cap1 {
		ax1, ay1 := destination(dx1, dy1, a-math.Pi*2, t1*cubicCorner)
		ax2, ay2 := destination(dx2, dy2, a-math.Pi*2, t1*cubicCorner)
		c.CubicTo(ax1, ay1, ax2, ay2, dx2, dy2)
	} else {
		c.LineTo(dx2, dy2)
	}
	c.LineTo(dx3, dy3)
	if cap2 {
		ax1, ay1 := destination(dx3, dy3, a-math.Pi*2, -t2*cubicCorner)
		ax2, ay2 := destination(dx4, dy4, a-math.Pi*2, -t2*cubicCorner)
		c.CubicTo(ax1, ay1, ax2, ay2, dx4, dy4)
	} else {
		c.LineTo(dx4, dy4)
	}
	c.LineTo(dx1, dy1)
	c.ClosePath()
	return nil
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

func (p *Pinhole) LoadObj(r io.Reader) error {
	var faces [][][3]float64
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	var verts [][3]float64
	for ln, line := range strings.Split(string(data), "\n") {
		for {
			nline := strings.Replace(line, "  ", " ", -1)
			if len(nline) < len(line) {
				line = nline
				continue
			}
			break
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "v ") {
			parts := strings.Split(line[2:], " ")
			if len(parts) >= 3 {
				verts = append(verts, [3]float64{})
				for j := 0; j < 3; j++ {
					if verts[len(verts)-1][j], err = strconv.ParseFloat(parts[j], 64); err != nil {
						return fmt.Errorf("line %d: %s", ln+1, err.Error())
					}
				}
			}
		} else if strings.HasPrefix(line, "f ") {
			parts := strings.Split(line[2:], " ")
			if len(parts) >= 3 {
				faces = append(faces, [][3]float64{})
				for _, part := range parts {
					part = strings.Split(part, "/")[0]
					idx, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						return fmt.Errorf("line %d: %s", ln+1, err.Error())
					}
					if int(idx) > len(verts) {
						return fmt.Errorf("line %d: invalid vert index: %d", ln+1, idx)
					}
					faces[len(faces)-1] = append(faces[len(faces)-1], verts[idx-1])
				}
			}
		}
	}
	for _, faces := range faces {
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
	return nil
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
