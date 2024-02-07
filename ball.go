package main

import "github.com/hajimehoshi/ebiten/v2"

type Ball struct {
	vector        *Point
	position      *Point
	last_position Point
	motion_state  int
	image         ebiten.Image
}
