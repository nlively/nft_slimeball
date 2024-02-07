package main

import "github.com/hajimehoshi/ebiten/v2"

type Ball struct {
	vector        *Point
	position      *Point
	last_position Point
	motion_state  int
	image         ebiten.Image
}

func (b *Ball) Dimensions() (width int, height int) {
	return b.image.Bounds().Dx(), b.image.Bounds().Dy()
}

func (b *Ball) Rect() (x1 int, y1 int, x2 int, y2 int) {
	pos := b.position

	int_x := int(pos.x)
	int_y := int(pos.y)
	int_x2 := int_x + b.image.Bounds().Dx()
	int_y2 := int_y + b.image.Bounds().Dy()

	return int_x, int_y, int_x2, int_y2
}
