package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
)

var player Player
var env Enviroment

//Game is the struct that implements ebiten.Game
type Game struct{}

//Update handles the logic. 60fps
func (g *Game) Update(screen *ebiten.Image) error {
	player.update(screen, &env)
	return nil
}

//Draw handles displaying each frame. 60fps
func (g *Game) Draw(screen *ebiten.Image) {
	player.draw3D(screen, &env)
	//enviroment.draw2D(screen)
	//player.draw2D(screen)
}

//Layout returns the size of the canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920 / 3, 1080 / 3
}

func main() {
	player = Player{}
	player.init(1920/1080, 10, 10, 0)

	env = Enviroment{}
	env.init(20, 20)
	env.cells = Level01

	//ebiten.SetWindowSize(500, 500)
	//ebiten.SetWindowTitle("Raycasting")
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
