package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 600

	GRAVITY           = 9.8 * 3  // Gravity acceleration
	TIME_STEP         = 1 / 30.0 // Time step for the simulation, in seconds
	BOUNCE_EFFICIENCY = .85      // To simulate energy loss on bounce, 1.0 for a perfect bounce, less for damping
	FLOOR_FRICTION    = 0.98

	STATE_START = 0
	STATE_SERVE = 1
	STATE_PLAY  = 2

	BALL_RESTING = 0
	BALL_MOVING  = 1
)

type Point struct {
	x float64
	y float64
}

func (p *Point) Describe() string {
	return fmt.Sprintf("(%f, %f)", p.x, p.y)
}

type Player struct {
	points   int
	position Point
	color    color.RGBA
	image    ebiten.Image
}

type Ball struct {
	vector        *Point
	position      *Point
	last_position Point
	motion_state  int
	image         ebiten.Image
}

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
	// fmt.Printf("Player1 right boundary: %d\n", player1_right_boundary)
	// print player1 position and width
	// fmt.Printf("Player1 position: %s, width: %d\n", g.player1.position.Describe(), g.player1.image.Bounds().Dx())

	// player1_right_boundary := g.netImage.Bounds().Dx()
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

	// fmt.Printf("Gravity: %f, Time step: %f\n", GRAVITY, TIME_STEP)

	// fmt.Printf("Applying gravity force of %f to velocity vector: %f, and new result is %f\n", GRAVITY*TIME_STEP, vec.y, vec.y+GRAVITY*TIME_STEP)
	// fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	vec.y += GRAVITY * TIME_STEP
	pos.x += vec.x * TIME_STEP
	pos.y += vec.y * TIME_STEP

	// fmt.Printf("Ball position: %s and velocity vector: %s\n", pos.Describe(), vec.Describe())
	// fmt.Printf("Ball position: %s and velocity vector: %s\n", g.ball.position.Describe(), g.ball.vector.Describe())

	int_x := int(pos.x)
	int_y := int(pos.y)
	int_x2 := int_x + int(ball_width)
	int_y2 := int_y + int(ball_height)

	// Describe box bounds
	// fmt.Printf("Ball x,y / x2,y2: (%d, %d) / (%d, %d)\n", int_x, int_y, int_x2, int_y2)

	if int_x <= 0 || int_x2 >= SCREEN_WIDTH {
		vec.x = -vec.x * BOUNCE_EFFICIENCY
		// fmt.Println("X axis collision, new vector x: ", vec.x)

		if pos.x <= 0 {
			pos.x = 0
		}
		if int_x2 >= SCREEN_WIDTH {
			pos.x = SCREEN_WIDTH - ball_width
		}
	}

	if int_y < 0 || int_y2 >= SCREEN_HEIGHT {
		vec.y = -vec.y * BOUNCE_EFFICIENCY
		// fmt.Printf("Y axis collision at %s, new vector y: %f\n", pos.Describe(), vec.y)

		if int_y < 0 {
			pos.y = 0
			// g.ball.motion_state = BALL_RESTING
		}
		if int_y2 >= SCREEN_HEIGHT {
			pos.y = SCREEN_HEIGHT - ball_height
			vec.x *= FLOOR_FRICTION
		}
		// describe velocity vector
		// fmt.Printf("Velocity vector: %s\n", vec.Describe())
	}

	if (pos.x == g.ball.last_position.x) && (pos.y == g.ball.last_position.y) {
		g.ball.motion_state = BALL_RESTING
	}

	g.ball.last_position = *pos
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
	//get width of g.netImage
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

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawNet(screen)

	g.DrawPlayer(screen, g.player1)
	g.DrawPlayer(screen, g.player2)

	g.DrawBall(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	fmt.Println("NFT Slimeball!")

	commander_fren, err := loadImage("images/commander_fren.png", 0.5, 0.5)
	if err != nil {
		log.Fatal(err)
	}

	heckboy, err := loadImage("images/heckboy.png", 0.5, 0.5)
	if err != nil {
		log.Fatal(err)
	}

	ball_image := createBallImage(32) // loadImage("images/ball.png", 1, 1)

	player1 := &Player{
		position: Point{x: 200, y: float64(SCREEN_HEIGHT - commander_fren.Bounds().Max.Y)},
		image:    *commander_fren,
	}
	player2 := &Player{
		position: Point{x: SCREEN_WIDTH - 200, y: float64(SCREEN_HEIGHT - heckboy.Bounds().Max.Y)},
		image:    *heckboy,
	}

	netImage := createNetImage()

	ball_velocity := calculateVelocityComponents(450, -45)

	ball := &Ball{
		position:     &Point{x: SCREEN_WIDTH / 2, y: SCREEN_HEIGHT / 2},
		image:        *ball_image,
		vector:       &ball_velocity,
		motion_state: BALL_RESTING,
	}

	game := &Game{
		player1:  player1,
		player2:  player2,
		ball:     ball,
		netImage: netImage,
	}
	game.setup()

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("NFT Slimeball")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
