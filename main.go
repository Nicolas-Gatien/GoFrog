package main

import (
	"image/color"
	"log"

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
	ebitenutil.DebugPrintAt(screen, "Here", cursorX, cursorY)

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64((GameWidth/2.0)-(g.FROG_IMAGE.Bounds().Dx()/2.0)), float64((GameHeight/2.0)-(g.FROG_IMAGE.Bounds().Dy()/2.0)))
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
