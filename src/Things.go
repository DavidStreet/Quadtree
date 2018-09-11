package main

type Rectangle struct {
	Pos
	h, w float64
}

type Point struct {
	Pos
	data   int
	status bool
}

type Pos struct {
	x, y float64
}

type Circle struct {
	Pos
	r float64
}

type Duple struct {
	a *Point
	b *Point
}

type Arr struct {
	slice []Pos
	xmin  float64
	xmax  float64
	ymin  float64
	ymax  float64
}

type Duples struct {
	*Duple
	next *Duples
	len  int
}

type Points struct {
	*Point
	next *Points
	len  int
}

func NewRect(x, y, w, h float64) *Rectangle {
	return &Rectangle{Pos: Pos{x, y}, h: h, w: w}
}

func NewPoint(x, y float64, data int) *Point {
	return &Point{Pos: Pos{x, y}, data: data, status: false}
}

func NewCircle(x, y, r float64) *Circle {
	return &Circle{Pos{x, y}, r}
}

func NewDuple(a, b *Point) *Duple {
	return &Duple{a, b}
}

func NewPoints() *Points {
	return &Points{Point: nil, next: nil, len: 0}
}

func NewDuples() *Duples {
	return &Duples{Duple: nil, next: nil, len: 0}
}

func DistSq(a, b Pos) float64 {
	return (a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y)
}

func (r *Rectangle) Contains(p Pos) bool {
	return (p.x <= r.x+r.w &&
		p.y <= r.y+r.h &&
		p.x >= r.x &&
		p.y >= r.y)
}

func (c *Circle) Contains(p Pos) bool {

	return DistSq(c.Pos, p) <= c.r*c.r
}

func (c *Circle) Intersects(r *Rectangle) bool {
	ra := c.r
	xr := r.x
	yr := r.y
	h := r.h
	w := r.w

	r1 := Rectangle{Pos{xr - ra, yr}, h, w + 2*ra}
	r2 := Rectangle{Pos{xr, yr - ra}, h + 2*ra, w}
	return (r1.Contains(c.Pos) ||
		r2.Contains(c.Pos) ||
		DistSq(r.Pos, c.Pos) <= ra*ra ||
		DistSq(Pos{xr + w, yr}, c.Pos) <= ra*ra ||
		DistSq(Pos{xr, yr + h}, c.Pos) <= ra*ra ||
		DistSq(Pos{xr + w, yr + h}, c.Pos) <= ra*ra)

}

func (a *Arr) Convert(x, y float64) (float64, float64) {
	p1 := mapp(x, a.xmin, a.xmax, 0, (a.xmax-a.xmin)*1000000.0/9.0)
	p2 := mapp(y, a.ymin, a.ymax, 0, (a.ymax-a.ymin)*1000000.0/9.0)
	return p1, p2
}

func mapp(value, istart, istop, ostart, ostop float64) float64 {
	return ostart + (ostop-ostart)*((value-istart)/(istop-istart))
}

func (ps *Points) Add(p *Point) {
	ps.next = NewPoints()
	ps.next.Point = p
	ps.len++
}

func (ds *Duples) Add(d *Duple) {
	ds.next = NewDuples()
	ds.next.Duple = d
	ds.len++
}
