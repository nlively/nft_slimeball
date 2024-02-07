package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameElement interface {
	Dimensions() (int, int)
	Rect() (int, int, int, int)
}

func Describe(e GameElement) {
	w, h := e.Dimensions()
	x1, y1, x2, y2 := e.Rect()
	fmt.Printf("Dimensions: %d, %d\n", w, h)
	fmt.Printf("Rect: %d, %d, %d, %d\n", x1, y1, x2, y2)
}

func (g *Game) Update() error {
	switch g.state {
	case STATE_SERVE:
		g.handleServeInput()
	case STATE_PLAY:
		g.handlePlayInput()
		if g.ball.motion_state == BALL_MOVING {
			g.moveBall()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawNet(screen)
	g.DrawPlayer(screen, g.player1)
	g.DrawPlayer(screen, g.player2)
	g.DrawBall(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}
