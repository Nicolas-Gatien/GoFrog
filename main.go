package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const GameWidth = 320 / 2
const GameHeight = 240 / 2

type Vector2 struct {
	x float64
	y float64
}

type Fly struct {
	position        Vector2
	currentFrame    int
	animationLength int
}

type Frog struct {
	open  bool
	angle float64
}

type Game struct {
	FROG_IMAGE      *ebiten.Image
	OPEN_FROG_IMAGE *ebiten.Image
	FLY_IMAGE       *ebiten.Image
	TONGUE_IMAGE    *ebiten.Image
	flies           []Fly
	time            int
	frog            Frog
}

func (g *Game) Update() error {
	g.time += 1
	for i, fly := range g.flies {
		if math.Mod(float64(g.time), 3.0) == 0 {
			if fly.currentFrame+1 >= fly.animationLength {
				g.flies[0].currentFrame = 0
			} else {
				g.flies[i].currentFrame += 1
			}
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		g.frog.open = true
	} else {
		g.frog.open = false
		cursorX, cursorY := ebiten.CursorPosition()
		centerX := GameWidth / 2
		centerY := GameHeight / 2
		offset := 90 * (math.Pi / 180)
		g.frog.angle = math.Atan2(float64(cursorY-centerY), float64(cursorX-centerX)) - offset
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	background := color.RGBA{103, 114, 169, 255}
	screen.Fill(background)

	opts := ebiten.DrawImageOptions{}

	opts.GeoM.Rotate(g.frog.angle)
	opts.GeoM.Translate(float64(GameWidth/2)-(8.0*math.Cos(g.frog.angle))+8.0*math.Sin(g.frog.angle), float64(GameHeight/2)-(8.0*math.Cos(g.frog.angle))-8.0*math.Sin(g.frog.angle))

	if g.frog.open {
		screen.DrawImage(g.OPEN_FROG_IMAGE, &opts)
	} else {
		screen.DrawImage(g.FROG_IMAGE, &opts)
	}

	for _, fly := range g.flies {
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Translate(fly.position.x, fly.position.y)
		screen.DrawImage(g.FLY_IMAGE.SubImage(image.Rect(fly.currentFrame*16, 0, (fly.currentFrame+1)*16, 16)).(*ebiten.Image), &opts)
	}

	tongue := g.TONGUE_IMAGE
	cursorX, cursorY := ebiten.CursorPosition()
	distance := math.Sqrt(math.Pow(float64(cursorX)-GameWidth/2, 2) + math.Pow(float64(cursorY)-GameHeight/2, 2))
	tongueOptions := ebiten.DrawImageOptions{}
	tongueOptions.GeoM.Scale(1, distance)
	tongueOptions.GeoM.Rotate(g.frog.angle)
	tongueOptions.GeoM.Translate(float64(GameWidth/2)-(3.0*math.Cos(g.frog.angle))+3.0*math.Sin(g.frog.angle), float64(GameHeight/2)-(3.0*math.Cos(g.frog.angle))-3.0*math.Sin(g.frog.angle))

	screen.DrawImage(tongue, &tongueOptions)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pond Frog")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	frogImage, _, err := ebitenutil.NewImageFromFile("frog.png")
	if err != nil {
		log.Fatal(err)
	}

	openFrogImage, _, err := ebitenutil.NewImageFromFile("frog_open.png")
	if err != nil {
		log.Fatal(err)
	}

	flyImage, _, err := ebitenutil.NewImageFromFile("fly_animation.png")
	if err != nil {
		log.Fatal(err)
	}

	tongueImage, _, err := ebitenutil.NewImageFromFile("tongue.png")
	if err != nil {
		log.Fatal(err)
	}

	flies := []Fly{
		{position: Vector2{30, 50}, animationLength: 6},
	}

	if err := ebiten.RunGame(&Game{TONGUE_IMAGE: tongueImage, FROG_IMAGE: frogImage, OPEN_FROG_IMAGE: openFrogImage, FLY_IMAGE: flyImage, flies: flies, frog: Frog{open: false}}); err != nil {
		log.Fatal(err)
	}
}
