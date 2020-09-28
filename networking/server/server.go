package main

import (
	"log"
	"math"
	"net"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/oyberntzen/Raycasting-in-Golang/game"
	"github.com/oyberntzen/Raycasting-in-Golang/game/graphics"
	"github.com/oyberntzen/Raycasting-in-Golang/game/levels"
	"github.com/oyberntzen/Raycasting-in-Golang/game/physics"
	"github.com/oyberntzen/Raycasting-in-Golang/networking"
)

var (
	cells        [][]uint8
	sprites      []game.Sprite
	players      sync.Map
	playerInputs sync.Map
	playerProts  sync.Map
)

const (
	width  int = 500
	height int = 500
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
	playersSlice := []game.Player{}
	players.Range(func(key interface{}, val interface{}) bool {
		player := val.(game.Player)
		playersSlice = append(playersSlice, player)
		return true
	})
	graphics.Draw2D(screen, cells, playersSlice, physics.PlayerSize, width, height)
}

//Layout returns the size of the canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	cells = levels.Level01.Cells
	sprites = levels.Level01.Sprites

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
	for {
		time.Sleep(time.Second / 10)
		playerInputs.Range(func(key, val interface{}) bool {
			id := key.(uint8)
			inputs := val.([]networking.Input)
			playerInputs.Store(id, []networking.Input{inputs[len(inputs)-1]})

			val, ok := players.Load(id)
			if !ok {
				players.Delete(id)
				playerInputs.Delete(id)
				playerProts.Delete(id)
				return true
			}

			player := val.(game.Player)

			player = physics.HandleInputs(player, inputs, cells)

			players.Store(id, player)

			return true
		})

		playerProts.Range(func(key, val interface{}) bool {
			prot := val.(networking.Protocol)
			id := key.(uint8)

			val, ok := players.Load(id)
			if !ok {
				players.Delete(id)
				playerInputs.Delete(id)
				playerProts.Delete(id)
				return true
			}

			snapshot := networking.Snapshot{}
			snapshot.ThisPlayer = val.(game.Player)
			snapshot.OtherPlayers = []game.Player{}

			players.Range(func(key interface{}, val interface{}) bool {
				otherPlayer := val.(game.Player)
				otherid := key.(uint8)
				if otherid != id {
					snapshot.OtherPlayers = append(snapshot.OtherPlayers, otherPlayer)
				}
				return true
			})

			err := prot.Send(snapshot, networking.SnapshotPacket)
			if err != nil {
				players.Delete(id)
				playerInputs.Delete(id)
				playerProts.Delete(id)
			}

			return true
		})
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
	thisPlayer := game.Player{PlayerID: id, X: 22.5, Y: 10.5, Z: 0, Angle: -math.Pi / 2, Pitch: 0, Health: 100}
	info := networking.ServerInfo{ThisPlayer: thisPlayer, Cells: cells, Sprites: sprites}
	players.Store(id, thisPlayer)

	lastTime := getTime()

	inputs := []networking.Input{networking.Input{TimeStamp: float32(lastTime)}}
	playerInputs.Store(id, inputs)

	playerProts.Store(id, prot)

	handleError(prot.Send(info, networking.ServerInfoPacket))

	for {
		//Handle message from client
		pid, data, err := prot.Recieve()
		if err != nil {
			players.Delete(id)
			playerInputs.Delete(id)
			playerProts.Delete(id)
			break
		}

		if pid != networking.NilPacket {
			if pid == networking.InputPacket {
				val, _ := playerInputs.Load(id)
				inputs := val.([]networking.Input)
				input := prot.DecodeInput(data)
				inputs = append(inputs, input)
				playerInputs.Store(id, inputs)
				/*time.Sleep(time.Second / 60)

				//Handle input sent from client
				input := prot.DecodeInput(data)

				nowTime := getTime()
				delta := nowTime - lastTime
				if delta < 0 {
					delta = nowTime - (lastTime - 60)
				}
				lastTime = getTime()

				thisPlayer.Angle += float64(input.MouseX) * 0.002
				if thisPlayer.Angle < -math.Pi {
					thisPlayer.Angle += math.Pi * 2
				} else if thisPlayer.Angle > math.Pi {
					thisPlayer.Angle -= math.Pi * 2
				}

				thisPlayer.Pitch = math.Max(math.Min(thisPlayer.Pitch-float64(input.MouseY)*0.002, 1), -1)

				dirXFor, dirYFor := math.Cos(thisPlayer.Angle), math.Sin(thisPlayer.Angle)
				nextX := dirXFor*float64(input.Up) - dirXFor*float64(input.Down)
				nextY := dirYFor*float64(input.Up) - dirYFor*float64(input.Down)

				dirXLeft, dirYLeft := math.Cos(thisPlayer.Angle-math.Pi/2), math.Sin(thisPlayer.Angle-math.Pi/2)
				nextX += dirXLeft*float64(input.Left) - dirXLeft*float64(input.Right)
				nextY += dirYLeft*float64(input.Left) - dirYLeft*float64(input.Right)

				angle := math.Atan2(nextY, nextX)
				if nextX != 0 {
					nextX = math.Cos(angle) * physics.PlayerSpeed * delta
				}
				if nextY != 0 {
					nextY = math.Sin(angle) * physics.PlayerSpeed * delta
				}

				thisPlayer.X, thisPlayer.Y = physics.Collision(thisPlayer.X, thisPlayer.Y, thisPlayer.X+nextX, thisPlayer.Y+nextY, cells)

				players.Store(id, thisPlayer)

				//Update Z based on velocity
				thisPlayer.Vel = math.Min(math.Max(thisPlayer.Vel-0.17*delta, -0.1), 0.1)
				thisPlayer.Z = math.Min(math.Max(thisPlayer.Z+thisPlayer.Vel, 0), 0.4)*/
			} else if pid == networking.EventPacket {
				event := prot.DecodeEvent(data)

				if event.Event == networking.JumpEvent {
					if thisPlayer.Z <= 0 {
						thisPlayer.Vel = 0.05
					}
				} else if event.Event == networking.ShootEvent {
					players.Range(func(key interface{}, val interface{}) bool {
						player := val.(game.Player)
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
			}

		}

		/*//Send snapshot
		snapshot := networking.Snapshot{}
		snapshot.ThisPlayer = thisPlayer
		snapshot.OtherPlayers = []game.Player{}

		players.Range(func(key interface{}, val interface{}) bool {
			player := val.(game.Player)
			pid := key.(uint8)
			if pid != id {
				snapshot.OtherPlayers = append(snapshot.OtherPlayers, player)
			}
			return true
		})

		err = prot.Send(snapshot, networking.SnapshotPacket)
		if err != nil {
			players.Delete(id)
			break

		}*/

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
