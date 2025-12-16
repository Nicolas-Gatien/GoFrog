package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const GameWidth = 320
const GameHeight = 240

type Vector2 struct {
	x float64
	y float64
}

type Fly struct {
	position Vector2
}

type Game struct {
	FROG_IMAGE *ebiten.Image
	FLY_IMAGE  *ebiten.Image
	flies      []Fly
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	background := color.RGBA{103, 114, 169, 255}
	screen.Fill(background)

	cursorX, cursorY := ebiten.CursorPosition()
	centerX := GameWidth / 2
	centerY := GameHeight / 2

	opts := ebiten.DrawImageOptions{}
	offset := 90 * (math.Pi / 180)
	angle := math.Atan2(float64(cursorY-centerY), float64(cursorX-centerX)) - offset
	opts.GeoM.Rotate(angle)
	opts.GeoM.Translate(float64(centerX)-(8.0*math.Cos(angle))+8.0*math.Sin(angle), float64(centerY)-(8.0*math.Cos(angle))-8.0*math.Sin(angle))

	screen.DrawImage(g.FROG_IMAGE, &opts)

	for _, fly := range g.flies {
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Translate(fly.position.x, fly.position.y)
		screen.DrawImage(g.FLY_IMAGE, &opts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pond Frog")

	frogImage, _, err := ebitenutil.NewImageFromFile("frog.png")
	if err != nil {
		log.Fatal(err)
	}

	flyImage, _, err := ebitenutil.NewImageFromFile("fly.png")
	if err != nil {
		log.Fatal(err)
	}

	flies := []Fly{
		Fly{position: Vector2{30, 50}},
	}

	if err := ebiten.RunGame(&Game{FROG_IMAGE: frogImage, FLY_IMAGE: flyImage, flies: flies}); err != nil {
		log.Fatal(err)
	}
}
