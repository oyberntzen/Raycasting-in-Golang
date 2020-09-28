package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/oyberntzen/Raycasting-in-Golang/game/physics"

	"github.com/enriquebris/goconcurrentqueue"

	"github.com/hajimehoshi/ebiten"
	"github.com/oyberntzen/Raycasting-in-Golang/game"
	"github.com/oyberntzen/Raycasting-in-Golang/game/graphics"
	"github.com/oyberntzen/Raycasting-in-Golang/networking"
)

var (
	cells   [][]uint8
	sprites []game.Sprite
	player  game.Player
	players []game.Player

	playerID       uint8
	gameState      int
	input          networking.Input
	mouseX, mouseY int16
	events         *goconcurrentqueue.FIFO
	pressedSpace   bool
	pressedShoot   bool
	firstShake     bool
	prot           networking.Protocol
)

const (
	width     int = 500
	height    int = 500
	scaleDown int = 2
)

//Game is the struct that implements ebiten.Game
type Game struct{}

//Exit implements error interface
type Exit struct{}

func (e *Exit) Error() string {
	return "Exit game"
}

//Update handles the logic
func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return &Exit{}
	}
	return nil
}

//Draw handles displaying each frame
func (g *Game) Draw(screen *ebiten.Image) {
	if gameState == 1 {
		graphics.Draw3D(screen, player, cells, sprites, players, width/scaleDown, height/scaleDown, physics.PlayerSize)
	}
}

//Layout returns the size of the canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / scaleDown, outsideHeight / scaleDown
}

func main() {

	c, _ := net.Dial("tcp", "localhost:8000")
	defer c.Close()

	go serverConnection(c)

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Raycasting")
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetRunnableOnUnfocused(true)
	//ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}

func serverConnection(conn net.Conn) {
	prot = networking.CreateProtocol(conn)

	for {

		//Handle incoming packet
		id, data, err := prot.Recieve()
		handleError(err)

		if id != networking.NilPacket {
			if id == networking.ServerInfoPacket && gameState == 0 {
				serverInfo := prot.DecodeServerInfo(data)

				path, _ := os.Getwd()
				imagesPath := filepath.Dir(filepath.Dir(path)) + "\\images\\"
				graphics.Init(imagesPath, width/scaleDown, height/scaleDown)

				cells = serverInfo.Cells
				sprites = serverInfo.Sprites

				player = serverInfo.ThisPlayer

				events = goconcurrentqueue.NewFIFO()
				playerID = serverInfo.ThisPlayer.PlayerID
				input = networking.Input{}
				players = []game.Player{}

				go updateInput()

				gameState = 1
			}
			if id == networking.SnapshotPacket {
				snapshot := prot.DecodeSnapshot(data)
				player = snapshot.ThisPlayer
				players = snapshot.OtherPlayers
			}
		}

		//Send pending event
		if events.GetLen() > 0 {
			val, err := events.Dequeue()
			handleError(err)
			if event, ok := val.(networking.Event); ok {
				handleError(prot.Send(event, networking.EventPacket))
			}
		} /*else {
			//Otherwise send input
			newX, newY := ebiten.CursorPosition()
			deltaX, deltaY := int16(newX*scaleDown)-mouseX, int16(newY*scaleDown)-mouseY
			if (deltaX < 100 && deltaY < 100 && deltaX > -100 && deltaY > -100) || firstShake {
				input.MouseX, input.MouseY = deltaX, deltaY
			} else {
				firstShake = true
			}
			mouseX, mouseY = int16(newX*scaleDown), int16(newY*scaleDown)

			handleError(prot.Send(input, networking.InputPacket))
			input = networking.Input{}
		}*/

	}
}

func updateInput() {
	last := float32(getTime())
	for {
		time.Sleep(time.Second / 120)
		now := float32(getTime())
		delta := now - last
		if delta > 50 {
			delta = (now - 60) - last
		}
		last = float32(getTime())
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			input.Up = delta
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			input.Down = delta
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			input.Left = delta
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			input.Right = delta
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			if !pressedSpace {
				//events.Enqueue(networking.Event{Event: networking.JumpEvent})
				input.Jump = true
				pressedSpace = true
			}
		} else {
			pressedSpace = false
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if !pressedShoot {
				//events.Enqueue(networking.Event{Event: networking.ShootEvent})
				input.Shoot = true
				pressedShoot = true
			}
		} else {
			pressedShoot = false
		}
		newX, newY := ebiten.CursorPosition()
		deltaX, deltaY := int16(newX*scaleDown)-mouseX, int16(newY*scaleDown)-mouseY
		if (deltaX < 100 && deltaY < 100 && deltaX > -100 && deltaY > -100) || firstShake {
			input.MouseX, input.MouseY = deltaX, deltaY
		} else {
			firstShake = true
		}
		mouseX, mouseY = int16(newX*scaleDown), int16(newY*scaleDown)

		input.TimeStamp = now

		handleError(prot.Send(input, networking.InputPacket))
		input = networking.Input{}
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getTime() float64 {
	now := time.Now()
	return float64(now.Nanosecond())/float64(time.Second) + float64(now.Second())
}
