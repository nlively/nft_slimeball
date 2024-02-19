package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Net struct {
	image ebiten.Image
}

func (n *Net) Describe() string {
	return "Net"
}

func (n *Net) Dimensions() (width int, height int) {
	return n.image.Bounds().Dx(), n.image.Bounds().Dy()
}

func (b *Net) IntRect() IntRect {
	int_x := SCREEN_WIDTH/2 - NET_WIDTH/2
	int_y := SCREEN_HEIGHT - NET_HEIGHT
	int_x2 := SCREEN_WIDTH/2 + NET_WIDTH/2
	int_y2 := SCREEN_HEIGHT

	return IntRect{int_x, int_y, int_x2, int_y2}
}

func (n *Net) Rect() Rect {
	int_x := float64(SCREEN_WIDTH/2 - NET_WIDTH/2)
	int_y := float64(SCREEN_HEIGHT - NET_HEIGHT)
	int_x2 := float64(SCREEN_WIDTH/2 + NET_WIDTH/2)
	int_y2 := float64(SCREEN_HEIGHT)

	return Rect{int_x, int_y, int_x2, int_y2}
}
