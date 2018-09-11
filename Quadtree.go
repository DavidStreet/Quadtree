package main

import (
	"fmt"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

type Quadtree struct {
	nw      *Quadtree
	ne      *Quadtree
	sw      *Quadtree
	se      *Quadtree
	divided bool
	limit   *Rectangle
	point   *Point
}

func NewQuadtree(x, y, w, h float64) *Quadtree {
	return &Quadtree{nw: nil, ne: nil, sw: nil, se: nil,
		divided: false, limit: NewRect(x, y, w, h), point: nil}
}

func (qt *Quadtree) Insert(p *Point) bool {
	if !qt.limit.Contains(p.Pos) {
		return false
	}
	if qt.point == nil {
		qt.point = p
		return true
	}
	if !qt.divided {
		qt.Subdivide()
	}
	if qt.nw.Insert(p) {
		return true
	}
	if qt.ne.Insert(p) {
		return true
	}
	if qt.sw.Insert(p) {
		return true
	}
	if qt.se.Insert(p) {
		return true
	}
	return false
}

func (qt *Quadtree) Subdivide() {
	x := qt.limit.x
	y := qt.limit.y
	h := qt.limit.h
	w := qt.limit.w
	qt.ne = NewQuadtree(x+w/2, y, w/2, h/2)
	qt.nw = NewQuadtree(x, y, w/2, h/2)
	qt.se = NewQuadtree(x+w/2, y+h/2, w/2, h/2)
	qt.sw = NewQuadtree(x, y+h/2, w/2, h/2)
	qt.divided = true
}

func (qt *Quadtree) Print() {
	qt.print(0)
}

func (qt *Quadtree) print(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("--")
	}
	if qt.point == nil {
		fmt.Println("null")
	} else {
		fmt.Printf("n: %d. ID: %d\n", n, qt.point.data)
	}
	if qt.divided {
		qt.nw.print(n + 1)
		qt.ne.print(n + 1)
		qt.sw.print(n + 1)
		qt.se.print(n + 1)
	}
}

// func (qt *Quadtree) Query(rang *Circle) *list.List {
// 	list := list.New()
// 	qt.queryAux(rang, list)
// 	//fmt.Println(list)
// 	return list
// }

// func (qt *Quadtree) queryAux(rang *Circle, list *list.List) {
// 	if !rang.Intersects(qt.limit) || qt.point == nil {
// 		return
// 	}
// 	if rang.Contains(qt.point.Pos) {
// 		list.PushBack(qt.point)
// 		fmt.Printf("an element found: %v.\n", qt.point.data)
// 	}
// 	if qt.divided {
// 		qt.nw.queryAux(rang, list)
// 		qt.ne.queryAux(rang, list)
// 		qt.sw.queryAux(rang, list)
// 		qt.se.queryAux(rang, list)
// 	}
// }

// func (qt *Quadtree) DetectCollitions(rang float64) *list.List {
// 	list := list.New()
// 	qt.detectAux(rang, list)
// 	return list
// }

// func (qt *Quadtree) detectAux(rang float64, list *list.List) {
// 	if qt.point == nil {
// 		return
// 	}

// 	if qt.divided {
// 		qt.nw.detectAux(rang, list)
// 		qt.nw.detectAux(rang, list)
// 		qt.nw.detectAux(rang, list)
// 		qt.nw.detectAux(rang, list)
// 	}
// 	qr := qt.Query(NewCircle(qt.point.x, qt.point.y, rang))
// 	fmt.Println(qr.Len())
// 	for e := qr.Front(); e != nil; e = e.Next() {
// 		pt := e.Value.(*Point)
// 		if qt.point.data < pt.data {
// 			fmt.Printf("adding element %v - %v.", qt.point.data, pt.data)
// 			list.PushBack(NewDuple(qt.point, pt))
// 			qt.point.status = true
// 			pt.status = true
// 		}

// 	}
// }

func (qt *Quadtree) Query(rang *Circle) []*Point {
	list := qt.queryAux(rang)
	//fmt.Println(list)
	return list
}

func (qt *Quadtree) queryAux(rang *Circle) []*Point {
	list := make([]*Point, 0)
	if !rang.Intersects(qt.limit) || qt.point == nil {
		return list
	}
	if rang.Contains(qt.point.Pos) {
		list = append(list, qt.point)
		fmt.Printf("an element found: %v.\n", qt.point.data)
	}
	if qt.divided {
		list = append(list, qt.nw.queryAux(rang)...)
		list = append(list, qt.ne.queryAux(rang)...)
		list = append(list, qt.sw.queryAux(rang)...)
		list = append(list, qt.se.queryAux(rang)...)
	}
	return list
}

func (qt *Quadtree) DetectCollitions(rang float64) []*Duple {
	list := qt.detectAux(rang)
	return list
}

func (qt *Quadtree) detectAux(rang float64) []*Duple {
	list := make([]*Duple, 0)
	if qt.point == nil {
		return list
	}

	qr := qt.Query(NewCircle(qt.point.x, qt.point.y, rang))
	fmt.Printf("List: %v\n len: %v.\n", qr, len(qr))
	for _, e := range qr {
		if qt.point.data < e.data {
			fmt.Printf("adding element %v - %v.", qt.point.data, e.data)
			list = append(list, NewDuple(qt.point, e))
			qt.point.status = true
			e.status = true
		}

	}
	if qt.divided {
		list = append(list, qt.nw.detectAux(rang)...)
		list = append(list, qt.nw.detectAux(rang)...)
		list = append(list, qt.nw.detectAux(rang)...)
		list = append(list, qt.nw.detectAux(rang)...)
	}
	return list
}

func (qt *Quadtree) Show(c *canvas.Context, scl float64) {
	c.SetColor(colornames.White)
	c.Stroke()
	c.DrawRectangle(qt.limit.x*scl, qt.limit.y*scl, qt.limit.w*scl, qt.limit.h*scl)

	if qt.divided {
		qt.nw.Show(c, scl)
		qt.ne.Show(c, scl)
		qt.sw.Show(c, scl)
		qt.se.Show(c, scl)
	}
	if qt.point != nil {
		r := 1.0
		if qt.point.status {
			c.SetColor(colornames.Orange)
			r = 10
		}
		c.DrawCircle(qt.point.x*scl, qt.point.y*scl, r)
	}
}
