package packages

import (
	"image"
	"image/color"

	"os"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

//Enviroment is the enviromental variables
type Enviroment struct {
	cellsx, cellsy int
	cells          [][]int
	sprites        []Sprite
	textures       []image.Image
}

//Init initializes the enviroment with the specified level
func (env *Enviroment) Init(level Level) {
	env.cellsx = len(level.cells[0])
	env.cellsy = len(level.cells)

	env.cells = level.cells
	env.sprites = level.sprites

	env.textures = append(env.textures,
		load("./images/eagle.png"),
		load("./images/redbrick.png"),
		load("./images/purplestone.png"),
		load("./images/greystone.png"),
		load("./images/bluestone.png"),
		load("./images/mossy.png"),
		load("./images/wood.png"),
		load("./images/colorstone.png"),
		load("./images/barrel.png"),
		load("./images/pillar.png"),
		load("./images/greenlight.png"))
}

//Draw2D draws the walls in a 2D map on the screen
func (env *Enviroment) Draw2D(screen *ebiten.Image) {
	cellsizex := screen.Bounds().Max.X / env.cellsx
	cellsizey := screen.Bounds().Max.Y / env.cellsy
	for y := 0; y < env.cellsy; y++ {
		for x := 0; x < env.cellsx; x++ {
			if env.cells[y][x] != 0 {
				ebitenutil.DrawRect(screen, float64(x*cellsizex), float64(y*cellsizey),
					float64(cellsizex), float64(cellsizey), color.RGBA{255, 0, 0, 255})
			}
		}
	}
}

//Level is a struct for levels
type Level struct {
	cells   [][]int
	sprites []Sprite
}

//Sprite is a struct a sprite
type Sprite struct {
	posX     float64
	posY     float64
	texture  int
	distance float64
}

func load(path string) image.Image {
	file, _ := os.Open(path)
	img, _, _ := image.Decode(file)
	return img
}

//Level01 is the first level
var Level01 Level = Level{[][]int{
	{8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 4, 4, 6, 4, 4, 6, 4, 6, 4, 4, 4, 6, 4},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{8, 0, 3, 3, 0, 0, 0, 0, 0, 8, 8, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6},
	{8, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6},
	{8, 0, 3, 3, 0, 0, 0, 0, 0, 8, 8, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 4, 0, 0, 0, 0, 0, 6, 6, 6, 0, 6, 4, 6},
	{8, 8, 8, 8, 0, 8, 8, 8, 8, 8, 8, 4, 4, 4, 4, 4, 4, 6, 0, 0, 0, 0, 0, 6},
	{7, 7, 7, 7, 0, 7, 7, 7, 7, 0, 8, 0, 8, 0, 8, 0, 8, 4, 0, 4, 0, 6, 0, 6},
	{7, 7, 0, 0, 0, 0, 0, 0, 7, 8, 0, 8, 0, 8, 0, 8, 8, 6, 0, 0, 0, 0, 0, 6},
	{7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 6, 0, 0, 0, 0, 0, 4},
	{7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 6, 0, 6, 0, 6, 0, 6},
	{7, 7, 0, 0, 0, 0, 0, 0, 7, 8, 0, 8, 0, 8, 0, 8, 8, 6, 4, 6, 0, 6, 6, 6},
	{7, 7, 7, 7, 0, 7, 7, 7, 7, 8, 8, 4, 0, 6, 8, 4, 8, 3, 3, 3, 0, 3, 3, 3},
	{2, 2, 2, 2, 0, 2, 2, 2, 2, 4, 6, 4, 0, 0, 6, 0, 6, 3, 0, 0, 0, 0, 0, 3},
	{2, 2, 0, 0, 0, 0, 0, 2, 2, 4, 0, 0, 0, 0, 0, 0, 4, 3, 0, 0, 0, 0, 0, 3},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 4, 0, 0, 0, 0, 0, 0, 4, 3, 0, 0, 0, 0, 0, 3},
	{1, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 4, 6, 0, 6, 3, 3, 0, 0, 0, 3, 3},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 1, 2, 2, 2, 6, 6, 0, 0, 5, 0, 5, 0, 5},
	{2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 2, 2, 0, 5, 0, 5, 0, 0, 0, 5, 5},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 5, 0, 5, 0, 5, 0, 5, 0, 5},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 5, 0, 5, 0, 5, 0, 5, 0, 5},
	{2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 2, 2, 0, 5, 0, 5, 0, 0, 0, 5, 5},
	{2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 5, 5, 5, 5, 5, 5, 5, 5, 5}},
	[]Sprite{
		Sprite{20.5, 11.5, 10, 0},
		Sprite{18.5, 4.5, 10, 0},
		Sprite{10.0, 4.5, 10, 0},
		Sprite{10.0, 12.5, 10, 0},
		Sprite{3.5, 6.5, 10, 0},
		Sprite{3.5, 20.5, 10, 0},
		Sprite{3.5, 14.5, 10, 0},
		Sprite{14.5, 20.5, 10, 0},
		Sprite{18.5, 10.5, 9, 0},
		Sprite{18.5, 11.5, 9, 0},
		Sprite{18.5, 12.5, 9, 0},
		Sprite{21.5, 1.5, 8, 0},
		Sprite{15.5, 1.5, 8, 0},
		Sprite{16.0, 1.8, 8, 0},
		Sprite{16.2, 1.2, 8, 0},
		Sprite{3.5, 2.5, 8, 0},
		Sprite{9.5, 15.5, 8, 0},
		Sprite{10.0, 15.1, 8, 0},
		Sprite{10.5, 15.8, 8, 0},
	},
}
