package main

type GameWindow struct {
}

func (g *GameWindow) Rect() Rect {
	return Rect{0, 0, SCREEN_WIDTH, SCREEN_HEIGHT}
}

func (g *GameWindow) IntRect() IntRect {
	return IntRect{0, 0, SCREEN_WIDTH, SCREEN_HEIGHT}
}

func (g *GameWindow) Describe() string {
	return "Walls"
}
