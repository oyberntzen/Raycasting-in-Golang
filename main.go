package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
)

var player Player
var env Enviroment

//Game is the struct that implements ebiten.Game
type Game struct{}

//Exit implements error interface
type Exit struct{}

func (e *Exit) Error() string {
	return "Exit game"
}

//Update handles the logic. 60fps
func (g *Game) Update(screen *ebiten.Image) error {
	player.update(screen, &env)
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return &Exit{}
	}
	return nil
}

//Draw handles displaying each frame. 60fps
func (g *Game) Draw(screen *ebiten.Image) {
	player.draw3D(screen, &env)
	//env.draw2D(screen)
	//player.draw2D(screen)
}

//Layout returns the size of the canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 255, 255
}

func main() {
	player = Player{}
	player.init(22.5, 10.5, -math.Pi/2)

	env = Enviroment{}
	env.init(Level01)

	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Raycasting")
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	//ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
