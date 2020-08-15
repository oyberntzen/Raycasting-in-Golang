package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"

	"os"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

//Level01 is the first level
var Level01 [][]int = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0},
	{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
	{0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 1, 0},
	{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 1, 0},
	{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
	{0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

//Enviroment is the enviromental variables
type Enviroment struct {
	cellsx, cellsy int
	cells          [][]int
	textures       []image.Image
}

func (env *Enviroment) init(cellsx, cellsy int) {
	env.cellsx = cellsx
	env.cellsy = cellsy

	for y := 0; y < env.cellsy; y++ {
		row := make([]int, env.cellsx)
		for x := 0; x < env.cellsx; x++ {
			row[x] = rand.Intn(2)
		}
		env.cells = append(env.cells, row)
	}

	f, _ := os.Open("./floor.png")
	f2, _ := os.Open("./wall.jpg")
	wall, _, _ := image.Decode(f)
	bricks2, _ := jpeg.Decode(f2)
	env.textures = append(env.textures, wall)
	env.textures = append(env.textures, bricks2)
}

func (env *Enviroment) draw2D(screen *ebiten.Image) {
	cellsizex := screen.Bounds().Max.X / env.cellsx
	cellsizey := screen.Bounds().Max.Y / env.cellsy
	for y := 0; y < env.cellsy; y++ {
		for x := 0; x < env.cellsx; x++ {
			if env.cells[y][x] == 1 {
				ebitenutil.DrawRect(screen, float64(x*cellsizex), float64(y*cellsizey),
					float64(cellsizex), float64(cellsizey), color.RGBA{255, 0, 0, 255})
			}
		}
	}
}
