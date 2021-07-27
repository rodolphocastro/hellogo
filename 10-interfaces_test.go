package main

import (
	"errors"
	"math"
	"testing"
)

type Point struct {
	x int
	y int
}

type Line struct {
	start, end Point
}

func (l Line) size() int {
	deltaX :=
		l.start.x - l.end.x
	deltaY :=
		l.start.y - l.end.y

	return int(math.Abs(float64(deltaX)) + math.Abs(float64(deltaY)))
}

func newPoint(xPos int, yPos int) Point {
	return Point{
		x: xPos,
		y: yPos,
	}
}

func newLine(start, end Point) Line {
	return Line{
		start: start,
		end:   end,
	}
}

// Interfaces in goLang are collection of methods one must supply.
type measure interface {
	size() int
}

// For instance DoMath trusts that anyone calling it will supply whatever object they have
// as far as said object matches the interface defined by "measure"
func DoMath(mes measure) int {
	return mes.size()
}

type Polygon struct {
	lines []Line
}

// Finally, in goLang, you don't need to specifically tell the compiler "yo, I'll implement that interface",
// you just need to create a method (func with receiver parameter) and call it a day
func (poly Polygon) size() int {
	result := 0
	for _, line := range poly.lines {
		result += line.size()
	}
	return result
}

func newPolygon(lines ...Line) (*Polygon, error) {
	if len(lines) <= 2 {
		return nil, errors.New("Polygons must have at least 3 lines")
	}
	poly := Polygon{
		lines: lines,
	}

	return &poly, nil
}

func TestInterfaces(t *testing.T) {
	pointA := newPoint(1, 1)
	pointB := newPoint(2, 1)
	lineAToB := newLine(pointA, pointB)
	// The size of lineAToB should be 1
	if DoMath(lineAToB) != 1 {
		t.Errorf("expected lineAtoB to have a size of %v, but found %v", 1, lineAToB.size())
	}
	// The size of a Polygon should be the sum of its lines
	poly, _ := newPolygon(lineAToB, lineAToB, lineAToB)
	if DoMath(poly) != 3 {
		t.Error("expected the polygon to have a size, found something else")
	}
}
