package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func calculateVelocityComponents(feetPerSecond float64, angleDegrees float64) Vector {
	angleRadians := angleDegrees * math.Pi / 180 // Convert angle from degrees to radians
	vx := feetPerSecond * math.Cos(angleRadians)
	vy := feetPerSecond * math.Sin(angleRadians)

	// Help me translate the feet per second to match the screen size
	// 1 foot = 32 pixels
	// 1 foot per second = 32 pixels per second
	// 1 second = 60 frames
	// 1 foot per second = 32 / 60 pixels per frame
	// 1 foot per second = 0.5333333333333333 pixels per frame

	fmt.Println("Vx:", vx, "Vy:", vy)
	return Vector{vx, vy}
}

func createBallImage(size int) *ebiten.Image {
	circleImage := ebiten.NewImage(size, size)
	circleColor := color.RGBA{0, 255, 0, 255} // Green color
	r := float64(size / 2)
	// Fill the circle image with green
	ebitenutil.DrawCircle(circleImage, r, r, r, circleColor)

	return circleImage
}

func createNetImage() *ebiten.Image {
	width := 15
	height := SCREEN_HEIGHT / 2

	netImage := ebiten.NewImage(width, height)

	// draw a rectangle that is 15 pixels wide and half the height of the screen
	ebitenutil.DrawRect(netImage, 0, 0, float64(width), float64(height), color.RGBA{0, 255, 0, 255})
	return netImage
}

func loadImage(path string, scaleX, scaleY float64) (*ebiten.Image, error) {
	// Open the image file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Get the original image dimensions
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// Calculate the new dimensions based on the scaling factors
	newWidth := int(float64(originalWidth) * scaleX)
	newHeight := int(float64(originalHeight) * scaleY)

	// Create a new Ebiten image with the scaled dimensions
	scaledImg := ebiten.NewImage(newWidth, newHeight)

	// Create drawing options and apply scaling
	opts := &ebiten.DrawImageOptions{}
	// The scale is already applied by the size of the new image,
	// so here we just draw the original image onto the scaled image.
	opts.GeoM.Scale(scaleX, scaleY)

	// Draw the original image onto the scaled image
	scaledImg.DrawImage(ebiten.NewImageFromImage(img), opts)

	return scaledImg, nil
}
