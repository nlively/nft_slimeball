package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player1       *Player
	player2       *Player
	ball          *Ball
	net           *Net
	state         int
	servingPlayer *Player
}

func (g *Game) setup() {
	g.state = STATE_SERVE
	g.servingPlayer = g.player1
}

func (g *Game) handleServeInput() {
	g.state = STATE_PLAY
	g.ball.motion_state = BALL_MOVING
	fmt.Println("Throwing ball at vector", g.ball.vector)
}

func (g *Game) handlePlayInput() {
	nr := g.net.IntRect()
	p1r := g.player1.IntRect()
	p2r := g.player2.IntRect()

	br := g.ball.IntRect()

	ballWithinPlayerHeight := br.y2 > p1r.y1 && br.y1 < p1r.y2
	ballBetweenPlayer1AndLeftWall := br.x1 <= p1r.x1 && br.x1 >= 0 && ballWithinPlayerHeight
	ballBetweenPlayer1AndNet := br.x1 >= p1r.x2 && br.x2 <= nr.x1 && ballWithinPlayerHeight
	ballBetweenPlayer2AndNet := br.x1 <= p2r.x1 && br.x2 >= nr.x2 && ballWithinPlayerHeight
	ballBetweenPlayer2AndRightWall := br.x2 >= p2r.x2 && br.x2 <= SCREEN_WIDTH && ballWithinPlayerHeight

	// Allow player1 to move left and right using A and D keys
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if (!ballBetweenPlayer1AndLeftWall && p1r.x1 >= 0) ||
			(ballBetweenPlayer1AndLeftWall && p1r.x1-4 >= br.x2) {
			g.player1.position.x -= 4
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if (!ballBetweenPlayer1AndNet && p1r.x2 <= nr.x1) ||
			(ballBetweenPlayer1AndNet && p1r.x2+4 <= br.x1) {
			fmt.Println("ball between player1 and net", ballBetweenPlayer1AndNet, "player right side", p1r.x2, "net left side", nr.x1, "ball left side", br.x1, "ball right side", br.x2)
			g.player1.position.x += 4
		}
	}

	// Allow player2 to move left and right using arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if (!ballBetweenPlayer2AndNet && p2r.x1 >= nr.x2) ||
			(ballBetweenPlayer2AndNet && p2r.x1-4 >= br.x2) {
			g.player2.position.x -= 4
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if (!ballBetweenPlayer2AndRightWall && p2r.x2 <= SCREEN_WIDTH) ||
			(ballBetweenPlayer2AndRightWall && p2r.x2+4 <= br.x1) {
			g.player2.position.x += 4
		}
	}
}

func (g *Game) detectWallCollision() {
	window := &GameWindow{}

	collision, collisionVector := detectInnerCollision(g.ball, window, *g.ball.vector, ZeroVector)

	if collision {
		g.ball.bounce(collisionVector)
	}
}

func (g *Game) detectNetCollision() {
	collision, collisionVector := detectCollision(g.ball, g.net, *g.ball.vector, ZeroVector, 5)

	if collision {
		g.ball.bounce(collisionVector)
	}
}

func (g *Game) detectPlayerCollision(player *Player) {
	collision, collisionVector := detectCollision(g.ball, player, *g.ball.vector, ZeroVector, 5)

	if collision {
		g.ball.bounce(collisionVector)
	}
}

func (g *Game) moveBall() {
	pos := g.ball.position
	vec := g.ball.vector

	accelerationY := GRAVITY * TIME_STEP

	vec.y += accelerationY
	pos.x += vec.x * TIME_STEP
	pos.y += vec.y * TIME_STEP

	g.detectWallCollision()
	g.detectNetCollision()
	g.detectPlayerCollision(g.player1)
	g.detectPlayerCollision(g.player2)

	// If the ball has stopped moving, set its motion state to resting
	if (pos.x == g.ball.last_position.x) && (pos.y == g.ball.last_position.y) {
		g.ball.motion_state = BALL_RESTING
	}

	// Note the ball's position for the next frame
	g.ball.last_position = *pos
}

func (g *Game) DrawPlayer(screen *ebiten.Image, player *Player) {
	// Create an options struct
	opts := &ebiten.DrawImageOptions{}

	x, y := float64(player.position.x), float64(player.position.y)
	opts.GeoM.Translate(x, y)

	// Draw the image to the screen with the specified options
	screen.DrawImage(&player.image, opts)
}

func (g *Game) DrawNet(screen *ebiten.Image) {
	// Get the rect of the net
	nr := g.net.Rect()
	nx1, ny1, nx2, ny2 := nr.x1, nr.y1, nr.x2, nr.y2

	// using ebiten v2, draw a rectangle based on the net rect
	ebitenutil.DrawRect(screen, nx1, ny1, nx2-nx1, ny2-ny1, color.RGBA{0, 255, 0, 255})
}

func (g *Game) DrawBall(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	x, y := float64(g.ball.position.x), float64(g.ball.position.y)
	opts.GeoM.Translate(x, y)

	screen.DrawImage(&g.ball.image, opts)
}
