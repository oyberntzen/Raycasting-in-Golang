package main

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/oyberntzen/Raycasting-in-Golang/game"
)

var env game.Enviroment
var player game.Player
var players map[int]game.Player
var playerSprites []game.Sprite
var notDraw int

var lastTime float64 = getTime()

const width int = 500
const height int = 500

var scaleDown int = 2

//Game is the struct that implements ebiten.Game
type Game struct{}

//Exit implements error interface
type Exit struct{}

func (e *Exit) Error() string {
	return "Exit game"
}

//Update handles the logic
func (g *Game) Update(screen *ebiten.Image) error {

	now := getTime()
	delta := now - lastTime
	if delta > 50 {
		delta = (now - 60) - lastTime
	}
	player.Update(screen, &env, delta, scaleDown)
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return &Exit{}
	}
	lastTime = getTime()
	return nil
}

//Draw handles displaying each frame
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
	players = make(map[int]game.Player)

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
		players = make(map[int]game.Player)
		handleError(dec.Decode(&players))
		updateSprites()
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
			playerSprites = append(playerSprites, game.NewSprite(player.PosX, player.PosY, 11, player.PosZ))
		}
	}
}

func getTime() float64 {
	now := time.Now()
	return float64(now.Nanosecond())/float64(time.Second) + float64(now.Second())
}
