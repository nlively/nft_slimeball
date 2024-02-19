package main

import (
	"fmt"
	"math"
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

// addVectorToRect moves the rectangle by the vector.
func addVectorToRect(rect Rect, vector Vector) Rect {
	return Rect{
		x1: rect.x1 + vector.x,
		y1: rect.y1 + vector.y,
		x2: rect.x2 + vector.x,
		y2: rect.y2 + vector.y,
	}
}

// detectCollision checks if two rectangles are colliding with a minimum 10% overlap.
// It returns a boolean indicating collision and the collision vector.
func detectCollision(objA, objB GameObject, vectorA, vectorB Vector) (bool, Vector) {
	rectA := objA.Rect()
	rectB := objB.Rect()

	movedRectA := addVectorToRect(rectA, vectorA)
	movedRectB := addVectorToRect(rectB, vectorB)

	overlapX := math.Min(movedRectA.x2, movedRectB.x2) - math.Max(movedRectA.x1, movedRectB.x1)
	overlapY := math.Min(movedRectA.y2, movedRectB.y2) - math.Max(movedRectA.y1, movedRectB.y1)

	if overlapX > 0 && overlapY > 0 {
		// Calculate the percentage overlap for each axis
		widthA := rectA.x2 - rectA.x1
		heightA := rectA.y2 - rectA.y1
		widthB := rectB.x2 - rectB.x1
		heightB := rectB.y2 - rectB.y1

		percentOverlapX := (overlapX / math.Min(widthA, widthB)) * 100
		percentOverlapY := (overlapY / math.Min(heightA, heightB)) * 100

		// Check for at least 10% overlap on both axes
		if percentOverlapX >= 10 && percentOverlapY >= 10 {
			vec := Vector{x: overlapX, y: overlapY}

			fmt.Println("Collision between ", objA.Describe(), " and ", objB.Describe(), " detected with vector:", collisionVector)

			return true, vec
		}
	}

	return false, Vector{}
}
