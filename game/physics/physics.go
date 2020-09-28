package physics

import (
	"math"

	"github.com/oyberntzen/Raycasting-in-Golang/game"
	"github.com/oyberntzen/Raycasting-in-Golang/game/graphics"
	"github.com/oyberntzen/Raycasting-in-Golang/networking"
)

const (
	//PlayerSize is the width of the player
	PlayerSize float64 = 0.3
	//PlayerSpeed is the maximum speed of the player in units per second
	PlayerSpeed float64 = 3
)

//Collision calculates new position based on walls
func Collision(startX, startY, endX, endY float64, cells [][]uint8) (float64, float64) {
	return doCollide(startX, endX, startY, true, cells), doCollide(startY, endY, startX, false, cells)
}

func doCollide(start, end, other float64, isX bool, cells [][]uint8) float64 {
	cend := end
	if end > start {
		cend += PlayerSize / 2
	} else {
		cend -= PlayerSize / 2
	}
	if round(cend) == round(start) {
		return end
	} else if cend > 0 && cend < float64(len(cells[0])) {
		if isX {
			if cells[int(other)][int(cend)] == 0 {
				return end
			}
		} else {
			if cells[int(cend)][int(other)] == 0 {
				return end
			}
		}
	} //else if end > 0 && end < float64(len(cells[0])) {
	//return end
	//}
	if end > start {
		return float64(int(cend)) - PlayerSize/2
	}
	return float64(int(start)) + PlayerSize/2
}

func round(number float64) int {
	if number < 0 {
		number--
	}
	return int(number)
}

//Hit calculates if a player aims at another player
func Hit(shootPlayer game.Player, otherPlayer game.Player, cells [][]uint8) bool {
	relX := otherPlayer.X - shootPlayer.X
	relY := otherPlayer.Y - shootPlayer.Y

	playerAngle := math.Atan2(relY, relX)
	dist := math.Sqrt(math.Pow(relX, 2) + math.Pow(relY, 2))
	angleWidth := math.Atan((PlayerSize / 2) / dist)

	dirX, dirY := game.Rotate(1, 0, shootPlayer.Angle)
	wallDist, _, _ := graphics.Ray(shootPlayer, cells, dirX, dirY)

	return (shootPlayer.Angle > playerAngle-angleWidth || shootPlayer.Angle > playerAngle-angleWidth+math.Pi*2) &&
		(shootPlayer.Angle < playerAngle+angleWidth || shootPlayer.Angle < playerAngle+angleWidth-math.Pi*2) && dist < wallDist
}

//HandleInputs moves the player with the inputs
func HandleInputs(player game.Player, inputs []networking.Input, cells [][]uint8) game.Player {
	for i, input := range inputs {
		if i == 0 {
			continue
		}
		delta := float64(input.TimeStamp - inputs[i-1].TimeStamp)
		if delta > 50 {
			delta -= 60
		}

		player.Angle += float64(input.MouseX) * 0.002
		if player.Angle < -math.Pi {
			player.Angle += math.Pi * 2
		} else if player.Angle > math.Pi {
			player.Angle -= math.Pi * 2
		}

		player.Pitch = math.Max(math.Min(player.Pitch-float64(input.MouseY)*0.002, 1), -1)

		dirXFor, dirYFor := math.Cos(player.Angle), math.Sin(player.Angle)
		nextX := dirXFor*float64(input.Up) - dirXFor*float64(input.Down)
		nextY := dirYFor*float64(input.Up) - dirYFor*float64(input.Down)

		dirXLeft, dirYLeft := math.Cos(player.Angle-math.Pi/2), math.Sin(player.Angle-math.Pi/2)
		nextX += dirXLeft*float64(input.Left) - dirXLeft*float64(input.Right)
		nextY += dirYLeft*float64(input.Left) - dirYLeft*float64(input.Right)

		angle := math.Atan2(nextY, nextX)
		if nextX != 0 {
			nextX = math.Cos(angle) * PlayerSpeed * delta
		}
		if nextY != 0 {
			nextY = math.Sin(angle) * PlayerSpeed * delta
		}

		player.X, player.Y = Collision(player.X, player.Y, player.X+nextX, player.Y+nextY, cells)

		if player.Z <= 0 && input.Jump {
			player.Vel = 0.05
		}
		player.Vel = math.Min(math.Max(player.Vel-0.17*delta, -0.1), 0.1)
		player.Z = math.Min(math.Max(player.Z+player.Vel, 0), 0.4)
	}

	return player
}
