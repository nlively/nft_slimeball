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

func detectInnerCollision(objA, objB GameObject, vectorA, vectorB Vector) (bool, Vector) {
	rectA := objA.Rect()
	rectB := objB.Rect()

	if vectorA.x > 0 && rectB.x2-rectA.x2 <= 0 {
		return true, Vector{-1, 0}
	} else if vectorA.x < 0 && rectA.x1-rectB.x1 <= 0 {
		return true, Vector{1, 0}
	}

	if vectorA.y > 0 && rectB.y2-rectA.y2 <= 0 {
		return true, Vector{0, -1}
	} else if vectorA.y < 0 && rectA.y1-rectB.y1 <= 0 {
		return true, Vector{0, 1}
	}

	return false, Vector{0, 0}
}

// detectCollision checks if two rectangles are colliding with a minimum 10% overlap.
// It returns a boolean indicating collision and the collision vector.
func detectCollision(objA, objB GameObject, vectorA, vectorB Vector, overlapThreshold float64) (bool, Vector) {
	rectA := objA.Rect()
	rectB := objB.Rect()

	// Calculate the relative movement of rectA towards rectB
	relativeVector := Vector{x: vectorA.x - vectorB.x, y: vectorA.y - vectorB.y}

	overlapX := math.Min(rectA.x2, rectB.x2) - math.Max(rectA.x1, rectB.x1)
	overlapY := math.Min(rectA.y2, rectB.y2) - math.Max(rectA.y1, rectB.y1)

	if overlapX > overlapThreshold && overlapY > overlapThreshold {

		fmt.Println("Collision between", objA.Describe(), "and", objB.Describe())

		// print relative vector and overlap
		fmt.Println("Relative Vector:", relativeVector, "OverlapX:", overlapX, "OverlapY:", overlapY)

		// Standard collision detection logic for outer walls
		if overlapX < overlapY {
			if relativeVector.x > 0 {
				return true, Vector{-1, 0} // Collision on A's right side
			} else {
				return true, Vector{1, 0} // Collision on A's left side
			}
		} else {
			if relativeVector.y > 0 {
				return true, Vector{0, -1} // Collision on A's bottom side
			} else {
				return true, Vector{0, 1} // Collision on A's top side
			}
		}
	}

	return false, Vector{0, 0}
}
