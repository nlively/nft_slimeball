package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

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

	ball_image := createBallImage(32)

	player1 := &Player{
		name:     "Commander Fren",
		position: Point{x: 200, y: float64(SCREEN_HEIGHT - commander_fren.Bounds().Max.Y)},
		image:    *commander_fren,
	}
	player2 := &Player{
		name:     "Heckboy",
		position: Point{x: SCREEN_WIDTH - 200, y: float64(SCREEN_HEIGHT - heckboy.Bounds().Max.Y)},
		image:    *heckboy,
	}

	netImage := createNetImage()

	ball_velocity := calculateVelocityComponents(450, 60)

	ball := &Ball{
		position:     &Point{x: 50, y: 50},
		image:        *ball_image,
		vector:       &ball_velocity,
		motion_state: BALL_RESTING,
	}

	net := &Net{
		image: *netImage,
	}

	game := &Game{
		player1: player1,
		player2: player2,
		ball:    ball,
		net:     net,
	}
	game.setup()

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("NFT Slimeball")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
