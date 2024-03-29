package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Ball struct {
	vector         *Vector
	position       *Point
	last_positions []Point
	motion_state   int
	image          ebiten.Image
}

func (b *Ball) hasStoppedMoving() bool {
	positions := len(b.last_positions)

	if positions < 5 {
		return false
	}

	for i := 1; i < positions; i++ {
		if !compareFloat64(b.last_positions[i].x, b.last_positions[i-1].x, 0.1) {
			return false
		}
		if !compareFloat64(b.last_positions[i].y, b.last_positions[i-1].y, 0.1) {
			return false
		}
	}

	return true
}

func compareFloat64(a, b, epsilon float64) bool {
	return (a-b) <= epsilon && (b-a) <= epsilon

}

func (b *Ball) savePosition(pos Point) {
	// Save up to 5 positions, then remove the oldest one
	if len(b.last_positions) == 5 {
		b.last_positions = b.last_positions[1:]
	}

	b.last_positions = append(b.last_positions, pos)
}

func (b *Ball) Describe() string {
	return "Ball"
}

func (b *Ball) setXVector(x float64, reason string) {
	fmt.Printf("%s: Changing vector x from %f to %f\n", reason, b.vector.x, x)
	b.vector.x = x
}

func (b *Ball) setYVector(y float64, reason string) {
	fmt.Printf("%s: Changing vector y from %f to %f\n", reason, b.vector.y, y)
	b.vector.y = y
}

func (b *Ball) Dimensions() (width int, height int) {
	return b.image.Bounds().Dx(), b.image.Bounds().Dy()
}

func (b *Ball) IntRect() IntRect {
	pos := b.position

	int_x := int(pos.x)
	int_y := int(pos.y)
	int_x2 := int_x + b.image.Bounds().Dx()
	int_y2 := int_y + b.image.Bounds().Dy()

	return IntRect{int_x, int_y, int_x2, int_y2}
}

func (b *Ball) Rect() Rect {
	pos := b.position

	int_x := pos.x
	int_y := pos.y
	int_x2 := int_x + float64(b.image.Bounds().Dx())
	int_y2 := int_y + float64(b.image.Bounds().Dy())

	return Rect{int_x, int_y, int_x2, int_y2}
}

func (b *Ball) bounce(collisionVector Vector) {
	newXVector := -b.vector.x * BOUNCE_EFFICIENCY
	newYVector := -b.vector.y * BOUNCE_EFFICIENCY

	switch collisionVector.x {
	case -1:
		b.setXVector(newXVector, "collision w left boundary")
	case 1:
		b.setXVector(newXVector, "collision w right boundary")
	}

	switch collisionVector.y {
	case -1:
		b.setYVector(newYVector, "collision w top boundary")
	case 1:
		b.setYVector(newYVector, "collision w bottom boundary")
	}
}
