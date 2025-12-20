package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	lifetime        int
}

const IDLE = "searching"
const ATTACKING = "attacking"
const RETREATING = "retreating"

func CenterScreen() Vector2 {
	return Vector2{x: GameWidth / 2, y: GameHeight / 2}
}

type Frog struct {
	open               bool
	angle              float64
	tongueLength       float64
	tongueTargetLength float64
	state              string
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

func Distance(pos1 Vector2, pos2 Vector2) float64 {
	return math.Sqrt(math.Pow(pos2.x-pos1.x, 2) + math.Pow(pos2.y-pos1.y, 2))
}

func MoveFly(fly *Fly) {
	//move := (math.Sin() * 2) + 1
	fly.position.y += math.Sin(float64(fly.lifetime/10)) / 2
}

func (g *Game) Update() error {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.time += 1
	cursorX, cursorY := ebiten.CursorPosition()

	if math.Mod(float64(g.time), 60) == 0 {
		g.flies = append(g.flies, Fly{position: Vector2{x: random.Float64() * GameWidth, y: random.Float64() * GameHeight}, animationLength: 6})
	}

	for i := range g.flies {
		fly := &g.flies[i]
		fly.lifetime += 1

		MoveFly(fly)

		if math.Mod(float64(g.time), 3.0) == 0 {
			if fly.currentFrame+1 >= fly.animationLength {
				fly.currentFrame = 0
			} else {
				fly.currentFrame += 1
			}
		}
	}

	switch g.frog.state {
	case ATTACKING:
		g.frog.open = true
		if g.frog.tongueLength < g.frog.tongueTargetLength {
			g.frog.tongueLength += 5
		} else {
			g.frog.state = RETREATING
		}
	case RETREATING:
		if g.frog.tongueLength > 0 {
			g.frog.tongueLength -= 10
		} else {
			g.frog.state = IDLE
			g.frog.tongueLength = 0
		}
	case IDLE:
		g.frog.open = false
		centerX := GameWidth / 2
		centerY := GameHeight / 2
		offset := 90 * (math.Pi / 180)
		g.frog.angle = math.Atan2(float64(cursorY-centerY), float64(cursorX-centerX)) - offset

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			g.frog.state = ATTACKING
			g.frog.tongueTargetLength = Distance(Vector2{x: GameWidth / 2, y: GameHeight / 2}, Vector2{x: float64(cursorX), y: float64(cursorY)})
		}
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
		tongue := g.TONGUE_IMAGE
		tongueOptions := ebiten.DrawImageOptions{}
		tongueOptions.GeoM.Scale(1, g.frog.tongueLength)
		tongueOptions.GeoM.Translate(0, 8.0)
		tongueOptions.GeoM.Rotate(g.frog.angle)
		tongueOptions.GeoM.Translate(float64(GameWidth/2)-(3.0*math.Cos(g.frog.angle))+3.0*math.Sin(g.frog.angle), float64(GameHeight/2)-(3.0*math.Cos(g.frog.angle))-3.0*math.Sin(g.frog.angle))
		screen.DrawImage(tongue, &tongueOptions)
	} else {
		screen.DrawImage(g.FROG_IMAGE, &opts)
	}

	for _, fly := range g.flies {
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Translate(fly.position.x, fly.position.y)
		screen.DrawImage(g.FLY_IMAGE.SubImage(image.Rect(fly.currentFrame*16, 0, (fly.currentFrame+1)*16, 16)).(*ebiten.Image), &opts)
	}

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

	if err := ebiten.RunGame(&Game{TONGUE_IMAGE: tongueImage, FROG_IMAGE: frogImage, OPEN_FROG_IMAGE: openFrogImage, FLY_IMAGE: flyImage, flies: flies, frog: Frog{open: false, state: IDLE}}); err != nil {
		log.Fatal(err)
	}
}
