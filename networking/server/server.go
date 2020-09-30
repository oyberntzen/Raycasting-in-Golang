package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/oyberntzen/Raycasting-in-Golang/game/graphics"
	"github.com/oyberntzen/Raycasting-in-Golang/game/levels"
	"github.com/oyberntzen/Raycasting-in-Golang/game/physics"
	"github.com/oyberntzen/Raycasting-in-Golang/networking"
)

var (
	cells                           [][]uint8
	sprites                         []networking.Sprite
	players                         map[uint8]networking.Player
	playerInputs                    map[uint8][]networking.Input
	playerProts                     map[uint8]networking.Protocol
	playerLock, inputLock, protLock sync.Mutex
)

const (
	width    int = 500
	height   int = 500
	timeStep int = 250
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
	playersSlice := []networking.Player{}

	playerLock.Lock()
	for _, player := range players {
		playersSlice = append(playersSlice, player)
	}
	playerLock.Unlock()

	graphics.Draw2D(screen, cells, playersSlice, physics.PlayerSize, width, height)
}

//Layout returns the size of the canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	cells = levels.Level01.Cells
	sprites = levels.Level01.Sprites

	players = make(map[uint8]networking.Player)
	playerInputs = make(map[uint8][]networking.Input)
	playerProts = make(map[uint8]networking.Protocol)

	l, _ := net.Listen("tcp", ":8000")
	defer l.Close()

	go mainLoop()
	go handlePlayers(l)

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Raycasting")
	ebiten.SetRunnableOnUnfocused(true)
	//ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func mainLoop() {
	var frame uint64 = 0
	for {
		start := getTime()

		shots := []networking.Player{}
		inputLock.Lock()
		for id, inputs := range playerInputs {
			for i, input := range inputs {
				if input.Shoot {
					playerLock.Lock()
					player := physics.HandleInputs(players[id], inputs[:i], cells)
					playerLock.Unlock()
					player.LastInputs = []networking.Input{input}
					shots = append(shots, player)
				}
			}
		}

		for id, inputs := range playerInputs {
			if len(playerInputs[id]) == 0 {
				continue
			}

			for _, shot := range shots {
				if shot.PlayerID == id {
					continue
				}
				timeStamp := shot.LastInputs[0].TimeStamp
				i := 0
				for {
					if i < len(inputs)-1 {
						if inputs[i+1].TimeStamp < timeStamp {
							i++
						}
					} else {
						i = len(inputs) - 1
						break
					}
				}
				fmt.Println(i, timeStamp, inputs[i].TimeStamp)
			}

			playerInputs[id] = []networking.Input{inputs[len(inputs)-1]}
			player := players[id]

			player = physics.HandleInputs(player, inputs, cells)
			player.LastInputNumber = inputs[len(inputs)-1].Number
			player.LastInputs = inputs

			players[id] = player
		}
		inputLock.Unlock()

		protLock.Lock()
		playerLock.Lock()
		for id, prot := range playerProts {
			snapshot := networking.Snapshot{}
			snapshot.ThisPlayer = players[id]
			snapshot.OtherPlayers = []networking.Player{}
			snapshot.Frame = frame

			for otherid, otherPlayer := range players {
				if otherid != id {
					snapshot.OtherPlayers = append(snapshot.OtherPlayers, otherPlayer)
				}
			}

			err := prot.Send(snapshot, networking.SnapshotPacket)
			if err != nil {
				deletePlayer(id)
			}
		}
		protLock.Unlock()
		playerLock.Unlock()
		frame++

		end := getTime()
		time.Sleep(time.Millisecond*time.Duration(timeStep) - time.Duration((end-start)*1000000000)*time.Nanosecond)
	}
}

func handlePlayers(l net.Listener) {
	for i := 0; true; i++ {
		c, _ := l.Accept()
		go playerConnection(c, uint8(i))
	}
}

func playerConnection(c net.Conn, id uint8) {
	prot := networking.CreateProtocol(c)
	thisPlayer := networking.Player{PlayerID: id, X: 22.5, Y: 10.5, Z: 0, Angle: -math.Pi / 2, Pitch: 0, Health: 100}
	info := networking.ServerInfo{ThisPlayer: thisPlayer, Cells: cells, Sprites: sprites}

	playerLock.Lock()
	players[id] = thisPlayer
	playerLock.Unlock()

	lastTime := getTime()

	inputs := []networking.Input{networking.Input{TimeStamp: float32(lastTime)}}
	inputLock.Lock()
	playerInputs[id] = inputs
	inputLock.Unlock()

	protLock.Lock()
	playerProts[id] = prot
	protLock.Unlock()

	handleError(prot.Send(info, networking.ServerInfoPacket))

	for {
		//Handle message from client
		pid, data, err := prot.Recieve()
		if err != nil {
			deletePlayer(id)
			break
		}

		if pid != networking.NilPacket {
			if pid == networking.InputPacket {
				inputLock.Lock()
				inputs := playerInputs[id]
				input := prot.DecodeInput(data)
				inputs = append(inputs, input)
				playerInputs[id] = inputs
				inputLock.Unlock()
			} /*else if pid == networking.EventPacket {
				event := prot.DecodeEvent(data)

				if event.Event == networking.JumpEvent {
					if thisPlayer.Z <= 0 {
						thisPlayer.Vel = 0.05
					}
				} else if event.Event == networking.ShootEvent {
					players.Range(func(key interface{}, val interface{}) bool {
						player := val.(networking.Player)
						pid := key.(uint8)
						if pid != id {
							if physics.Hit(thisPlayer, player, cells) {
								player.Health = uint8(math.Max(float64(player.Health-10), 0))
								players.Store(pid, player)
							}
						}
						return true
					})
				}
			}*/

		}
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func rotate(x, y, a float64) (float64, float64) {
	return x*math.Cos(a) - y*math.Sin(a), y*math.Cos(a) + x*math.Sin(a)
}

func getTime() float64 {
	now := time.Now()
	return float64(now.Nanosecond())/float64(time.Second) + float64(now.Second())
}

func deletePlayer(id uint8) {
	playerLock.Lock()
	inputLock.Lock()
	protLock.Lock()

	delete(players, id)
	delete(playerInputs, id)
	delete(playerProts, id)

	playerLock.Unlock()
	inputLock.Unlock()
	protLock.Unlock()
}
