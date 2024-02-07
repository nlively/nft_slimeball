package main

import "fmt"

type Point struct {
	x float64
	y float64
}

func (p *Point) Describe() string {
	return fmt.Sprintf("(%f, %f)", p.x, p.y)
}
