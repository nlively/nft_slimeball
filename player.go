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

func (p *Player) Describe() string {
	return p.name
}

func (p *Player) Dimensions() (width int, height int) {
	return p.image.Bounds().Dx(), p.image.Bounds().Dy()
}

func (b *Player) IntRect() IntRect {
	pos := b.position

	int_x := int(pos.x)
	int_y := int(pos.y)
	int_x2 := int_x + b.image.Bounds().Dx()
	int_y2 := int_y + b.image.Bounds().Dy()

	return IntRect{int_x, int_y, int_x2, int_y2}
}

func (b *Player) Rect() Rect {
	pos := b.position

	int_x := pos.x
	int_y := pos.y
	int_x2 := int_x + float64(b.image.Bounds().Dx())
	int_y2 := int_y + float64(b.image.Bounds().Dy())

	return Rect{int_x, int_y, int_x2, int_y2}
}
