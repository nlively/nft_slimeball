package main

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	player1       *Player
	player2       *Player
	ball          *Ball
	netImage      *ebiten.Image
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
	// get coordinates of left boundary of net
	player1_right_boundary := float64(SCREEN_WIDTH/2 - (g.netImage.Bounds().Dx() / 2))
	player2_left_boundary := float64(SCREEN_WIDTH/2 + (g.netImage.Bounds().Dx() / 2))

	player1_width := float64(g.player1.image.Bounds().Dx())
	player2_width := float64(g.player2.image.Bounds().Dx())

	// Allow player1 to move left and right using A and D keys
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.player1.position.x >= 0 {
		g.player1.position.x -= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && g.player1.position.x+player1_width <= player1_right_boundary {
		g.player1.position.x += 4
	}

	// Allow player2 to move left and right using arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.player2.position.x >= player2_left_boundary {
		g.player2.position.x -= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.player2.position.x+player2_width <= SCREEN_WIDTH {
		g.player2.position.x += 4
	}
}

func (g *Game) moveBall() {
	pos := g.ball.position
	vec := g.ball.vector

	ball_width := float64(g.ball.image.Bounds().Dx())
	ball_height := float64(g.ball.image.Bounds().Dy())

	vec.y += GRAVITY * TIME_STEP
	pos.x += vec.x * TIME_STEP
	pos.y += vec.y * TIME_STEP

	int_x := int(pos.x)
	int_y := int(pos.y)
	int_x2 := int_x + int(ball_width)
	int_y2 := int_y + int(ball_height)

	if int_x <= 0 || int_x2 >= SCREEN_WIDTH {
		vec.x = -vec.x * BOUNCE_EFFICIENCY

		if pos.x <= 0 {
			pos.x = 0
		}
		if int_x2 >= SCREEN_WIDTH {
			pos.x = SCREEN_WIDTH - ball_width
		}
	}

	if int_y < 0 || int_y2 >= SCREEN_HEIGHT {
		vec.y = -vec.y * BOUNCE_EFFICIENCY

		if int_y < 0 {
			pos.y = 0
		}
		if int_y2 >= SCREEN_HEIGHT {
			pos.y = SCREEN_HEIGHT - ball_height
			vec.x *= FLOOR_FRICTION
		}
	}

	if (pos.x == g.ball.last_position.x) && (pos.y == g.ball.last_position.y) {
		g.ball.motion_state = BALL_RESTING
	}

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
	lineThickness := float64(g.netImage.Bounds().Dx())

	// Position the scaled image to the middle of the screen, starting at the halfway point on the y-axis
	opts.GeoM.Translate(midScreenX-lineThickness/2, midScreenY)

	// Draw the line on the screen
	screen.DrawImage(g.netImage, opts)
}

func (g *Game) DrawBall(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	x, y := float64(g.ball.position.x), float64(g.ball.position.y)
	opts.GeoM.Translate(x, y)

	screen.DrawImage(&g.ball.image, opts)
}
