package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	points   int
	position Point
	color    color.RGBA
	image    ebiten.Image
}
