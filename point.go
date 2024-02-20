package main

import (
	"fmt"
)

type Point struct {
	x float64
	y float64
}

func (p *Point) Describe() string {
	return fmt.Sprintf("(%f, %f)", p.x, p.y)
}

// Rect represents a rectangle with its top-left (x1, y1) and bottom-right (x2, y2) corners.
type Rect struct {
	x1, y1, x2, y2 float64
}

type IntRect struct {
	x1, y1, x2, y2 int
}

// Vector represents a 2D motion vector.
type Vector struct {
	x, y float64
}

var ZeroVector = Vector{0, 0}

// interface for a game object with Rect, IntRect and Describe methods
type GameObject interface {
	Rect() Rect
	IntRect() IntRect
	Describe() string
}
