package main

import "github.com/hajimehoshi/ebiten/v2"

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
