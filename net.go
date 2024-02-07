package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Net struct {
	image ebiten.Image
}

func (n *Net) Dimensions() (width int, height int) {
	return n.image.Bounds().Dx(), n.image.Bounds().Dy()
}

func (b *Net) Rect() (x1 int, y1 int, x2 int, y2 int) {
	int_x := b.image.Bounds().Min.X
	int_y := b.image.Bounds().Min.Y
	int_x2 := int_x + b.image.Bounds().Dx()
	int_y2 := int_y + b.image.Bounds().Dy()

	return int_x, int_y, int_x2, int_y2
}
