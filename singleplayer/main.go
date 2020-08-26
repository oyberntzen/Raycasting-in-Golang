package main

import (
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/oyberntzen/Raycasting-in-Golang/game"
)

var player game.Player
var enviroment game.Enviroment

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
	player.Update(screen, &enviroment)
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return &Exit{}
	}
	return nil
}

//Draw handles displaying each frame. 60fps
func (g *Game) Draw(screen *ebiten.Image) {
	player.Draw3D(screen, &enviroment, []game.Sprite{})
	//env.draw2D(screen)
	//player.draw2D(screen)
}

//Layout returns the size of the canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / scaleDown, outsideHeight / scaleDown
}

func main() {
	player = game.Player{}
	player.Init(22.5, 10.5, -math.Pi/2, width/scaleDown, height/scaleDown)

	enviroment = game.Enviroment{}
	path, _ := os.Getwd()
	imagesPath := filepath.Dir(path) + "/images/"
	enviroment.Init(game.Level01, imagesPath)

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Raycasting")
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
