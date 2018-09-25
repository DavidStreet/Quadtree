package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/h8gi/canvas"
)

var name, numOfp string
var points *Arr
var qt *Quadtree
var c *canvas.Canvas
var scl float64

func main() {
	numOfp = "10"
	scl = 0.15
	points = &Arr{make([]Pos, 0), 0, 0, 0, 0}
	name = "../DataSets/ConjuntoDeDatosCon" + numOfp + "abejas.txt"
	fil, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer fil.Close()
	rdr := csv.NewReader(fil)
	rdr.Comment = 'C'
	rdr.FieldsPerRecord = 2

	reg0, err := rdr.Read()
	if err != nil {
		panic(err)
	}

	x, e1 := strconv.ParseFloat(reg0[0], 64)
	y, e2 := strconv.ParseFloat(reg0[1], 64)
	if e1 == nil && e2 == nil {
		points.slice = append(points.slice, Pos{x, y})
		points.xmax = x
		points.xmin = x
		points.ymax = y
		points.ymin = y
	}

	for reg, err := rdr.Read(); err != io.EOF; reg, err = rdr.Read() {
		if err != nil {
			panic(err)
		}

		lat, e1 := strconv.ParseFloat(reg[0], 64)
		lon, e2 := strconv.ParseFloat(reg[1], 64)

		if e1 != nil || e2 != nil {
			continue
		}

		points.slice = append(points.slice, Pos{x: lat, y: lon})
		if lat < points.xmin {
			points.xmin = lat
		}
		if lat > points.xmax {
			points.xmax = lat
		}
		if lon < points.ymin {
			points.ymin = lon
		}
		if lon > points.ymax {
			points.ymax = lon
		}
	}

	xm := math.Ceil((points.xmax-points.xmin)*1000000.0/9.0) + 10
	ym := math.Ceil((points.ymax-points.ymin)*1000000.0/9.0) + 10
	fmt.Println(int(ym))
	c = canvas.New(&canvas.NewCanvasOptions{
		Width:     int(xm*scl + 30),
		Height:    int(ym*scl + 30),
		FrameRate: 1,
		Title:     "Quadtree",
	})

	qt = NewQuadtree(0, 0, xm, ym)
	for i, p := range points.slice {
		x1, y1 := points.Convert(p.x, p.y)
		qt.Insert(NewPoint(x1, y1, i+1))
	}
	fmt.Println(len(points.slice))

	dng := qt.DetectCollitions(100)
	for _, dp := range dng {
		fmt.Printf("%v - %v\n", dp.a.data, dp.b.data)
	}
	fmt.Printf("%v\n------------------\n", len(dng))
	c.Draw(func(ctx *canvas.Context) {
		qt.Show(ctx, scl)
	})
}

// func main() {
// 	qt := NewQuadtree(0, 0, 100, 100)
// 	qt.Insert(NewPoint(20, 20, 1))
// 	qt.Insert(NewPoint(60, 60, 2))
// 	dt := qt.DetectCollitions(100)
// 	c := canvas.New(&canvas.NewCanvasOptions{
// 		Width:     450,
// 		Height:    450,
// 		FrameRate: 1,
// 	})

// 	c.Draw(func(ctx *canvas.Context) {
// 		qt.Show(ctx, 4)
// 	})

// 	for e := dt.Front(); e != nil; e = e.Next() {
// 		val := e.Value.(*Duple)
// 		fmt.Println(val.a)
// 		fmt.Println(val.b)
// 	}
// }
