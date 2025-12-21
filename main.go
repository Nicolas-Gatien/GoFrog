package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const GameWidth = 320 / 2
const GameHeight = 240 / 2

type Vector2 struct {
	x float64
	y float64
}

func Euler(angle float64, magnitude float64) Vector2 {
	return Vector2{x: math.Cos(angle) * magnitude, y: math.Sin(angle) * magnitude}
}

type CatchEffect struct {
	position        Vector2
	currentFrame    int
	animationLength int
}

func NewCatchEffect(position Vector2) CatchEffect {
	return CatchEffect{
		position:        position,
		currentFrame:    0,
		animationLength: 4,
	}
}

type Fly struct {
	position        Vector2
	currentFrame    int
	animationLength int
	lifetime        int
	state           string
}

const IDLE = "searching"
const ATTACKING = "attacking"
const RETREATING = "retreating"
const HIT = "hit"

func CenterScreen() Vector2 {
	return Vector2{x: GameWidth / 2, y: GameHeight / 2}
}

type Frog struct {
	open               bool
	angle              float64
	tongueLength       float64
	tongueTargetLength float64
	state              string
	health             int
}

func (frog *Frog) tonguePosition() Vector2 {
	magnitude := Euler(frog.angle+(90*(math.Pi/180)), frog.tongueLength)
	center := CenterScreen()
	position := Vector2{center.x + magnitude.x, center.y + magnitude.y}
	return position
}

type Game struct {
	FROG_IMAGE      *ebiten.Image
	OPEN_FROG_IMAGE *ebiten.Image
	FLY_IMAGE       *ebiten.Image
	TONGUE_IMAGE    *ebiten.Image
	CATCH_IMAGE     *ebiten.Image
	flies           []Fly
	catchEffects    []CatchEffect
	time            int
	frog            Frog
	audioContext    *audio.Context
	audioChomp      *audio.Player
	audioCatch      *audio.Player
}

func Distance(pos1 Vector2, pos2 Vector2) float64 {
	return math.Sqrt(math.Pow(pos2.x-pos1.x, 2) + math.Pow(pos2.y-pos1.y, 2))
}

func GetAngleTo(target Vector2, self Vector2) float64 {
	return math.Atan2(float64(target.y-self.y), float64(target.x-self.x))
}

func MoveFly(fly *Fly) {
	sway := math.Sin(float64(fly.lifetime/10)) / 2
	fly.position.y += sway

	angle := GetAngleTo(CenterScreen(), fly.position)
	movement := Euler(angle, 1)

	fly.position.x += movement.x * 0.2
	fly.position.y += movement.y * 0.2
}

func (g *Game) Update() error {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.time += 1
	cursorX, cursorY := ebiten.CursorPosition()

	if math.Mod(float64(g.time), 60) == 0 {
		spawnAngle := (random.Float64() * 360) * (math.Pi / 180)
		spawnMagnitude := Euler(spawnAngle, (GameWidth/2)+64)
		spawnPosition := Vector2{CenterScreen().x + spawnMagnitude.x, CenterScreen().y + spawnMagnitude.y}

		g.flies = append(g.flies, Fly{position: spawnPosition, animationLength: 6, state: ATTACKING})
	}

	killEffectIndex := []int{}
	for i := range g.catchEffects {
		effect := &g.catchEffects[i]
		if math.Mod(float64(g.time), 3.0) == 0 {
			if effect.currentFrame+1 < effect.animationLength {
				effect.currentFrame += 1
			} else {
				killEffectIndex = append(killEffectIndex, i)
			}
		}
	}

	for i := range killEffectIndex {
		g.catchEffects[killEffectIndex[i]] = g.catchEffects[len(g.catchEffects)-1]
		g.catchEffects = g.catchEffects[:len(g.catchEffects)-1]
	}

	killFlyIndex := []int{}
	for i := range g.flies {
		fly := &g.flies[i]

		switch fly.state {
		case HIT:
			if Distance(fly.position, CenterScreen()) < 5 {
				killFlyIndex = append(killFlyIndex, i)
				continue
			}
			fly.position = g.frog.tonguePosition()
		case ATTACKING:
			if g.frog.state == ATTACKING && Distance(fly.position, g.frog.tonguePosition()) < 8 {
				fly.state = HIT
				g.frog.state = RETREATING
				g.catchEffects = append(g.catchEffects, NewCatchEffect(fly.position))
				g.audioCatch.Rewind()
				g.audioCatch.Play()
				continue
			}
			if Distance(CenterScreen(), fly.position) < 5 {
				g.frog.health -= 1
				killFlyIndex = append(killFlyIndex, i)
				continue
			}

			MoveFly(fly)

			if math.Mod(float64(g.time), 3.0) == 0 {
				if fly.currentFrame+1 >= fly.animationLength {
					fly.currentFrame = 0
				} else {
					fly.currentFrame += 1
				}
			}
		}

		fly.lifetime += 1
	}

	for i := range killFlyIndex {
		g.flies[killFlyIndex[i]] = g.flies[len(g.flies)-1]
		g.flies = g.flies[:len(g.flies)-1]
		g.audioCatch.Pause()
		g.audioChomp.Rewind()
		g.audioChomp.Play()
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
		offset := 90 * (math.Pi / 180)
		g.frog.angle = GetAngleTo(Vector2{x: float64(cursorX), y: float64(cursorY)}, CenterScreen()) - offset

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			g.frog.state = ATTACKING
			g.frog.tongueTargetLength = Distance(Vector2{x: GameWidth / 2, y: GameHeight / 2}, Vector2{x: float64(cursorX), y: float64(cursorY)})
		}
	}

	if g.frog.health < 1 {
		g.Reset()
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
		opts.GeoM.Translate(-8, -8)
		screen.DrawImage(g.FLY_IMAGE.SubImage(image.Rect(fly.currentFrame*16, 0, (fly.currentFrame+1)*16, 16)).(*ebiten.Image), &opts)
	}

	for _, effect := range g.catchEffects {
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Translate(effect.position.x, effect.position.y)
		opts.GeoM.Translate(-8, -8)
		screen.DrawImage(g.CATCH_IMAGE.SubImage(image.Rect(effect.currentFrame*16, 0, (effect.currentFrame+1)*16, 16)).(*ebiten.Image), &opts)
	}

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.frog.health))

	/*debug := ebiten.NewImage(1, 1)
	debug.Fill(color.RGBA{255, 0, 0, 255})
	debugOpts := ebiten.DrawImageOptions{}
	debugOpts.GeoM.Translate(g.frog.tonguePosition().x, g.frog.tonguePosition().y)
	screen.DrawImage(debug, &debugOpts)*/
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

func (game *Game) Reset() {
	game.flies = []Fly{
		{position: Vector2{x: 10, y: 10}, animationLength: 6, currentFrame: 0, state: ATTACKING},
	}
	game.time = 0
	game.frog = Frog{open: false, state: IDLE, health: 3}
	game.catchEffects = []CatchEffect{}

	game.audioContext = audio.NewContext(48000)
	game.audioChomp = CreatePlayerForSoundFile("assets/sounds/chomp.wav", game.audioContext)
	game.audioCatch = CreatePlayerForSoundFile("assets/sounds/catch.wav", game.audioContext)
}

func CreatePlayerForSoundFile(path string, context *audio.Context) *audio.Player {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	d, err := wav.DecodeF32(file)
	if err != nil {
		log.Fatal(err)
	}

	player, err := context.NewPlayerF32(d)
	if err != nil {
		log.Fatal(err)
	}
	return player
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pond Frog")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	frogImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/frog.png")
	if err != nil {
		log.Fatal(err)
	}

	openFrogImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/frog_open.png")
	if err != nil {
		log.Fatal(err)
	}

	flyImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/fly_animation.png")
	if err != nil {
		log.Fatal(err)
	}

	tongueImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/tongue.png")
	if err != nil {
		log.Fatal(err)
	}

	catchImage, _, err := ebitenutil.NewImageFromFile("assets/sprites/catch_effect.png")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		TONGUE_IMAGE:    tongueImage,
		FROG_IMAGE:      frogImage,
		OPEN_FROG_IMAGE: openFrogImage,
		FLY_IMAGE:       flyImage,
		CATCH_IMAGE:     catchImage,
	}
	game.Reset()

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
