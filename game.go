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
	// g.ball.position.x = SCREEN_WIDTH / 2
	// g.ball.position.y = SCREEN_HEIGHT / 2
	g.state = STATE_PLAY
	g.ball.motion_state = BALL_MOVING
	fmt.Println("Throwing ball at vector", g.ball.vector)
}

func (g *Game) handlePlayInput() {
	nr := g.net.IntRect()
	p1r := g.player1.IntRect()
	p2r := g.player2.IntRect()

	// Allow player1 to move left and right using A and D keys
	if ebiten.IsKeyPressed(ebiten.KeyA) && p1r.x1 >= 0 {
		g.player1.position.x -= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && p1r.x2 <= nr.x1 {
		g.player1.position.x += 4
	}

	// Allow player2 to move left and right using arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && p2r.x1 >= nr.x2 {
		g.player2.position.x -= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && p2r.x2 <= SCREEN_WIDTH {
		g.player2.position.x += 4
	}
}

func (g *Game) detectWallCollision() {
	vec := g.ball.vector
	br := g.ball.IntRect()
	bx1, by1, bx2, by2 := br.x1, br.y1, br.x2, br.y2
	ball_width, ball_height := g.ball.Dimensions()

	if bx1 <= 0 || bx2 >= SCREEN_WIDTH {
		// fmt.Println("Wall collision detected, vec.x was", vec.x, " and now ", -vec.x*BOUNCE_EFFICIENCY)
		g.ball.setXVector(-vec.x*BOUNCE_EFFICIENCY, "wall collision")

		if bx1 <= 0 {
			g.ball.position.x = 0
		}
		if bx2 >= SCREEN_WIDTH {
			g.ball.position.x = float64(SCREEN_WIDTH - ball_width)
		}
	}

	if by1 < 0 || by2 >= SCREEN_HEIGHT {
		fmt.Println("Wall collision detected")
		g.ball.setYVector(-vec.y*BOUNCE_EFFICIENCY, "wall collision")

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
	br := g.ball.IntRect()
	nr := g.net.IntRect()
	bx1, by1, bx2, by2 := br.x1, br.y1, br.x2, br.y2
	nx1, ny1, nx2, ny2 := nr.x1, nr.y1, nr.x2, nr.y2

	if by2 > ny1 {
		// show ball and net rect
		fmt.Println("Ball rect: ", bx1, by1, bx2, by2, "; Net rect: ", nx1, ny1, nx2, ny2)

		// fmt.Println("Net collision detected")
		if bx2 > nx1 || bx1 < nx2 {
			g.ball.setXVector(-vec.x*BOUNCE_EFFICIENCY, "net collision")
			// } else {
			// vec.y = -vec.y * BOUNCE_EFFICIENCY
		}
		// detect if ball is touching top of net
	}
}

func (g *Game) detectPlayerCollision(player *Player) {
	vec := g.ball.vector
	br := g.ball.IntRect()
	bx1, by1, bx2, by2 := br.x1, br.y1, br.x2, br.y2

	pr := player.IntRect()
	px1, py1, px2, py2 := pr.x1, pr.y1, pr.x2, pr.y2

	collision, collisionVector := detectCollision(g.ball.Rect(), player.Rect(), *g.ball.vector, ZeroVector)
	if collision {
		fmt.Println("Player ", player.name, " detected with vector:", collisionVector)
	}

	playerLeftBoundaryCollision := bx2 >= px1 && bx1 <= px1 && by2 >= py1 && by1 <= py2
	playerRightBoundaryCollision := bx1 <= px2 && bx2 >= px2 && by2 >= py1 && by1 <= py2
	playerTopBoundaryCollision := by2 >= py1 && by1 <= py1 && bx2 >= px1 && bx1 <= px2

	if playerLeftBoundaryCollision || playerTopBoundaryCollision || playerRightBoundaryCollision {
		if playerLeftBoundaryCollision && vec.x > 0 {
			g.ball.setXVector(-vec.x, "collision w player left boundary")
			fmt.Println("Ball rect: ", bx1, by1, bx2, by2, "; Player rect: ", px1, py1, px2, py2)
		} else if playerRightBoundaryCollision && vec.x < 0 {
			g.ball.setXVector(-vec.x, "collision w player right boundary")
		}

		if playerTopBoundaryCollision && vec.y > 0 {
			g.ball.setYVector(-vec.y, "collision w player top boundary")
		}
	}
}

func (g *Game) moveBall() {
	pos := g.ball.position
	vec := g.ball.vector

	accelerationY := GRAVITY * TIME_STEP

	vec.y += accelerationY
	pos.x += vec.x * TIME_STEP
	pos.y += vec.y * TIME_STEP

	// print position of ball and vector
	// fmt.Println("Acceleration: ", accelerationY, "; Ball position: ", pos.x, pos.y, "; Ball vector: ", vec.x, vec.y)

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
