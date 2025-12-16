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

type Game struct {
	FROG_IMAGE *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	background := color.RGBA{103, 114, 169, 255}
	screen.Fill(background)

	cursorX, cursorY := ebiten.CursorPosition()
	frogX := GameWidth / 2
	frogY := GameHeight / 2

	opts := ebiten.DrawImageOptions{}
	offset := 90 * (math.Pi / 180)
	angle := math.Atan2(float64(cursorY-frogY), float64(cursorX-frogX)) - offset
	opts.GeoM.Rotate(angle)
	opts.GeoM.Translate(float64(frogX)-(8.0*math.Cos(angle))+8.0*math.Sin(angle), float64(frogY)-(8.0*math.Cos(angle))-8.0*math.Sin(angle))

	screen.DrawImage(g.FROG_IMAGE, &opts)
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

	if err := ebiten.RunGame(&Game{FROG_IMAGE: frogImage}); err != nil {
		log.Fatal(err)
	}
}
