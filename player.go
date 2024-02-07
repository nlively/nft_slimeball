package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	name     string
	points   int
	position Point
	color    color.RGBA
	image    ebiten.Image
}

func (p *Player) Dimensions() (width int, height int) {
	return p.image.Bounds().Dx(), p.image.Bounds().Dy()
}

func (p *Player) Rect() (x1 int, y1 int, x2 int, y2 int) {
	pos := p.position

	int_x := int(pos.x)
	int_y := int(pos.y)
	int_x2 := int_x + p.image.Bounds().Dx()
	int_y2 := int_y + p.image.Bounds().Dy()

	return int_x, int_y, int_x2, int_y2
}
