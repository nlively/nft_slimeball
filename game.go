package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
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
	g.ball.position.x = SCREEN_WIDTH / 2
	g.ball.position.y = SCREEN_HEIGHT / 2
	g.state = STATE_PLAY
	g.ball.motion_state = BALL_MOVING
}

func (g *Game) handlePlayInput() {
	nx1, _, nx2, _ := g.net.Rect()
	p1x1, _, p1x2, _ := g.player1.Rect()
	p2x1, _, p2x2, _ := g.player2.Rect()

	// Allow player1 to move left and right using A and D keys
	if ebiten.IsKeyPressed(ebiten.KeyA) && p1x1 >= 0 {
		g.player1.position.x -= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && p1x2 <= nx1 {
		g.player1.position.x += 4
	}

	// Allow player2 to move left and right using arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && p2x1 >= nx2 {
		g.player2.position.x -= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && p2x2 <= SCREEN_WIDTH {
		g.player2.position.x += 4
	}
}

func (g *Game) detectWallCollision() {
	vec := g.ball.vector
	bx1, by1, bx2, by2 := g.ball.Rect()
	ball_width, ball_height := g.ball.Dimensions()

	if bx1 <= 0 || bx2 >= SCREEN_WIDTH {
		fmt.Println("Wall collision detected")
		vec.x = -vec.x * BOUNCE_EFFICIENCY

		if bx1 <= 0 {
			g.ball.position.x = 0
		}
		if bx2 >= SCREEN_WIDTH {
			g.ball.position.x = float64(SCREEN_WIDTH - ball_width)
		}
	}

	if by1 < 0 || by2 >= SCREEN_HEIGHT {
		fmt.Println("Wall collision detected")
		vec.y = -vec.y * BOUNCE_EFFICIENCY

		if by1 < 0 {
			g.ball.position.y = 0
		}
		if by2 >= SCREEN_HEIGHT {
			g.ball.position.y = float64(SCREEN_HEIGHT - ball_height)
			vec.x *= FLOOR_FRICTION
		}
	}
}

func (g *Game) detectNetCollision() {
	// pos := g.ball.position
	vec := g.ball.vector
	bx1, _, bx2, by2 := g.ball.Rect()
	nx1, ny1, nx2, _ := g.net.Rect()

	// ball_width, ball_height := g.ball.Dimensions()

	// get x position of net's left boundary

	// print net boundaries
	// fmt.Println("Net left, right, top: ", net_left_boundary, net_right_boundary, net_top_boundary)

	// get y position of net's top boundary

	if by2 > ny1 {
		// fmt.Println("Net collision detected")
		if bx2 > nx1 || bx1 < nx2 {
			vec.x = -vec.x * BOUNCE_EFFICIENCY
			// } else {
			// vec.y = -vec.y * BOUNCE_EFFICIENCY
		}
		// detect if ball is touching top of net
	}
}

func (g *Game) detectPlayerCollision(player *Player) {
	// pos := g.ball.position
	vec := g.ball.vector
	bx1, by1, bx2, by2 := g.ball.Rect()
	// ball_width, ball_height := g.ball.Dimensions()

	px1, py1, px2, py2 := player.Rect()
	// player_width, player_height := player.Dimensions()

	// ball's right boundary has collided with player's left boundary
	// or ball's left boundary has collided with player's right boundary
	// or ball's bottom boundary has collided with player's top boundary
	if bx2 >= px1 && bx1 <= px2 && by2 >= py1 && by1 <= py2 {
		fmt.Println("Player collision detected ", player.name)
		if bx2 >= px1 && bx1 <= px2 && by2 >= py1 {
			vec.x = -vec.x * BOUNCE_EFFICIENCY
		} else {
			vec.y = -vec.y * BOUNCE_EFFICIENCY
		}
		// vec.y = -vec.y * BOUNCE_EFFICIENCY
		// if by1 < py1 {
		// 	g.ball.position.y = float64(py1 - ball_height)
		// } else {
		// 	g.ball.position.y = float64(py2)
		// }
		// if bx1 < px1 {
		// 	g.ball.position.x = float64(px1 - ball_width)
		// } else {
		// 	g.ball.position.x = float64(px2)
		// }
	}
}

func (g *Game) moveBall() {
	pos := g.ball.position
	vec := g.ball.vector

	vec.y += GRAVITY * TIME_STEP
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
	opts := &ebiten.DrawImageOptions{}

	midScreenX := float64(SCREEN_WIDTH) / 2
	midScreenY := float64(SCREEN_HEIGHT) / 2
	lineThickness, _ := g.net.Dimensions()

	// Position the scaled image to the middle of the screen, starting at the halfway point on the y-axis
	opts.GeoM.Translate(midScreenX-float64(lineThickness)/2, midScreenY)

	// Draw the line on the screen
	screen.DrawImage(&g.net.image, opts)
}

func (g *Game) DrawBall(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	x, y := float64(g.ball.position.x), float64(g.ball.position.y)
	opts.GeoM.Translate(x, y)

	screen.DrawImage(&g.ball.image, opts)
}
