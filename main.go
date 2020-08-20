package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
)

var player Player
var enviroment Enviroment

const width int = 1920
const height int = 1080
const scaleDown int = 3

//Game is the struct that implements ebiten.Game
type Game struct{}

//Exit implements error interface
type Exit struct{}

func (e *Exit) Error() string {
	return "Exit game"
}

//Update handles the logic. 60fps
func (g *Game) Update(screen *ebiten.Image) error {
	player.update(screen, &enviroment)
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return &Exit{}
	}
	return nil
}

//Draw handles displaying each frame. 60fps
func (g *Game) Draw(screen *ebiten.Image) {
	player.draw3D(screen, &enviroment)
	//env.draw2D(screen)
	//player.draw2D(screen)
}

//Layout returns the size of the canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / scaleDown, outsideHeight / scaleDown
}

func main() {
	player = Player{}
	player.init(22.5, 10.5, -math.Pi/2, width/scaleDown, height/scaleDown)

	enviroment = Enviroment{}
	enviroment.init(Level01)

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Raycasting")
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
