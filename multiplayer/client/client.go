package main

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/oyberntzen/Raycasting-in-Golang/game"
)

var env game.Enviroment
var player game.Player
var players []game.Player
var playerSprites []game.Sprite
var notDraw int

const width int = 500
const height int = 500
const scaleDown int = 1

//Game is the struct that implements ebiten.Game
type Game struct{}

//Exit implements error interface
type Exit struct{}

func (e *Exit) Error() string {
	return "Exit game"
}

//Update handles the logic. 60fps
func (g *Game) Update(screen *ebiten.Image) error {
	player.Update(screen, &env)
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return &Exit{}
	}
	return nil
}

//Draw handles displaying each frame. 60fps
func (g *Game) Draw(screen *ebiten.Image) {
	player.Draw3D(screen, &env, playerSprites)
	//env.Draw2D(screen)
	//player.Draw2D(screen, &env)
}

//Layout returns the size of the canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / scaleDown, outsideHeight / scaleDown
}

func main() {
	env = game.Enviroment{}
	player = game.Player{}

	path, _ := os.Getwd()
	imagesPath := filepath.Dir(filepath.Dir(path)) + "/images/"
	env.InitTextures(imagesPath)
	gob.Register(game.Enviroment{})

	c, _ := net.Dial("tcp", "localhost:8000")
	defer c.Close()

	enc := gob.NewEncoder(c)
	dec := gob.NewDecoder(c)

	handleError(enc.Encode(width / scaleDown))
	handleError(enc.Encode(height / scaleDown))

	handleError(dec.Decode(&env))
	handleError(dec.Decode(&player))
	handleError(dec.Decode(&notDraw))
	player.InitConstants()

	go serverConnection(enc, dec)

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Raycasting")
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	//ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}

func serverConnection(enc *gob.Encoder, dec *gob.Decoder) {
	for {
		handleError(enc.Encode(player))
		handleError(dec.Decode(&players))
		updateSprites()
		//updateSprites(lastLenght)
		//lastLenght = len(players)
	}

}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func updateSprites() {
	playerSprites = []game.Sprite{}
	for i, player := range players {
		if i != notDraw {
			playerSprites = append(playerSprites, game.NewSprite(player.PosX, player.PosY, 11))
		}
	}
}
