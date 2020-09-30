package game

import (
	"log"
	"math"
)

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
