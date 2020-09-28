package game

import (
	"log"
	"math"
)

//Sprite contains position and texture of a sprite
type Sprite struct {
	X, Y, Z, W, H float64
	Texture       uint8
}

//Player contains a snapshot for a player. This includes xyz, pitch and angle
type Player struct {
	PlayerID                   uint8
	X, Y, Z, Angle, Pitch, Vel float64
	Health                     uint8
}

//Rotate rotates a vector
func Rotate(x, y, a float64) (float64, float64) {
	return x*math.Cos(a) - y*math.Sin(a), y*math.Cos(a) + x*math.Sin(a)
}

//HandleError handles an error if ocurred
func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
