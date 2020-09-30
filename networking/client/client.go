package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/oyberntzen/Raycasting-in-Golang/game/physics"

	"github.com/enriquebris/goconcurrentqueue"

	"github.com/hajimehoshi/ebiten"
	"github.com/oyberntzen/Raycasting-in-Golang/game/graphics"
	"github.com/oyberntzen/Raycasting-in-Golang/networking"
)

var (
	cells   [][]uint8
	sprites []networking.Sprite
	player  networking.Player
	players []networking.Player

	playerID       uint8
	gameState      int
	input          networking.Input
	mouseX, mouseY int16
	events         *goconcurrentqueue.FIFO
	pressedSpace   bool
	pressedShoot   bool
	firstShake     bool
	prot           networking.Protocol

	oldPlayers []networking.Player
	oldInputs  []networking.Input
	inputLock  sync.Mutex

	frame            uint64
	lastOtherPlayers []networking.Player
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
				players = []networking.Player{}

				go updateInput()

				gameState = 1
			}
			if id == networking.SnapshotPacket {
				snapshot := prot.DecodeSnapshot(data)
				newPlayer := snapshot.ThisPlayer
				frame = snapshot.Frame

				players = make([]networking.Player, len(lastOtherPlayers))
				for i := 0; i < len(lastOtherPlayers); i++ {
					players[i] = lastOtherPlayers[i]
					for _, p := range snapshot.OtherPlayers {
						if p.PlayerID == players[i].PlayerID && p.PlayerID != playerID {
							players[i].LastInputs = make([]networking.Input, len(p.LastInputs))
							copy(players[i].LastInputs, p.LastInputs)
							go updateOtherPlayer(i, frame)
						}
					}

				}

				lastOtherPlayers = make([]networking.Player, len(snapshot.OtherPlayers))
				for i := 0; i < len(snapshot.OtherPlayers); i++ {
					lastOtherPlayers[i] = snapshot.OtherPlayers[i]
				}

				if len(oldPlayers) > 0 {
					inputLock.Lock()
					firstInputNumber := oldPlayers[0].LastInputNumber
					if firstInputNumber < newPlayer.LastInputNumber {
						diff := newPlayer.LastInputNumber - firstInputNumber
						checkPlayer := oldPlayers[diff]
						if !(newPlayer.Angle == checkPlayer.Angle && newPlayer.Pitch == checkPlayer.Pitch &&
							newPlayer.X == checkPlayer.X && newPlayer.Y == checkPlayer.Y && newPlayer.Z == checkPlayer.Z) {
							player = physics.HandleInputs(newPlayer, oldInputs[diff:], cells)
						}
						oldPlayers = oldPlayers[diff+1:]
						oldInputs = oldInputs[diff+1:]
					}
					inputLock.Unlock()
				}

			}
		}

		//Send pending event
		/*if events.GetLen() > 0 {
			val, err := events.Dequeue()
			handleError(err)
			if event, ok := val.(networking.Event); ok {
				handleError(prot.Send(event, networking.EventPacket))
			}
		} else {
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
	lastInput := networking.Input{TimeStamp: last}
	var number uint64 = 0
	for {
		time.Sleep(time.Second / 120)
		now := float32(getTime())
		delta := now - last

		if delta < 0 {
			delta += 60
		}
		last = now

		if ebiten.IsKeyPressed(ebiten.KeyW) {
			input.Up = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			input.Down = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			input.Left = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			input.Right = true
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
		input.Number = number

		player = physics.HandleInputs(player, []networking.Input{lastInput, input}, cells)
		player.LastInputNumber = number
		lastInput = input

		inputLock.Lock()
		oldPlayers = append(oldPlayers, player)
		oldInputs = append(oldInputs, input)
		inputLock.Unlock()

		handleError(prot.Send(input, networking.InputPacket))
		input = networking.Input{}

		number++

	}
}

func updateOtherPlayer(i int, f uint64) {
	inputs := players[i].LastInputs
	for num, input := range inputs {
		if frame != f {
			return
		}
		if num > 0 {
			players[i] = physics.HandleInputs(players[i], []networking.Input{inputs[num-1], input}, cells)
			delta := input.TimeStamp - inputs[num-1].TimeStamp
			if delta < 0 {
				delta += 60
			}
			time.Sleep(time.Duration(delta*1000000000) * time.Nanosecond)
		}
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
