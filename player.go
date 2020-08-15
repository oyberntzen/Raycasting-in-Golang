package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

//Player is the player in the game
type Player struct {
	x     float64
	y     float64
	angle float64

	rays [][2]float64
	fov  float64

	frontRay     float64
	maxDist      int
	intersection [2]float64
}

func (player *Player) init(ratio float64, x int, y int, angle float64) {
	player.x = float64(x)
	player.y = float64(y)
	player.angle = angle
	player.fov = 65 * math.Pi / 180 * ratio
	player.maxDist = 20
}

func (player *Player) update(screen *ebiten.Image, env *Enviroment) {
	changed := false
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		changed = true
		player.x += math.Cos(player.angle) * 0.02
		player.y += math.Sin(player.angle) * 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		changed = true
		player.x -= math.Cos(player.angle) * 0.02
		player.y -= math.Sin(player.angle) * 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		changed = true
		player.angle -= 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		changed = true
		player.angle += 0.02
	}

	player.x = math.Max(0, math.Min(float64(env.cellsx), player.x))
	player.y = math.Max(0, math.Min(float64(env.cellsy), player.y))

	if changed {
		player.frontRay = player.ray(env, 0)[0]
		player.rays = [][2]float64{}
		width := screen.Bounds().Max.X
		for ray := 0; ray < width; ray++ {
			angledif := (float64(ray)/float64(width))*player.fov - player.fov/2
			r := player.ray(env, angledif)
			var real float64
			if int(r[0]) < player.maxDist {
				real = math.Abs(math.Cos(angledif) * r[0])
			} else {
				real = r[0]
			}

			player.rays = append(player.rays, [2]float64{real, r[1]})
		}
	}
}

func (player *Player) ray(env *Enviroment, angledif float64) [2]float64 {
	angle := player.angle + angledif
	dist := float64(0)

	cos := math.Cos(angle)
	sin := math.Sin(angle)

	curcell := [2]int{int(player.x), int(player.y)}
	curpos := [2]float64{player.x, player.y}

	for int(dist) < player.maxDist {

		relx := curpos[0] - float64(curcell[0])
		rely := curpos[1] - float64(curcell[1])

		hor := 1 - relx
		if cos <= 0 {
			hor = -relx
		}
		ver := 1 - rely
		if sin <= 0 {
			ver = -rely
		}

		var hormult float64
		if cos != 0 {
			hormult = hor / cos
		} else {
			hormult = 1000
		}
		var vermult float64
		if sin != 0 {
			vermult = ver / sin
		} else {
			vermult = 1000
		}

		horlen := math.Pow(hor, 2) + math.Pow(sin*hormult, 2)
		verlen := math.Pow(ver, 2) + math.Pow(cos*vermult, 2)

		var textureIndex float64
		if horlen < verlen {
			if hor < 0 {
				curcell = [2]int{curcell[0] - 1, curcell[1]}
			} else {
				curcell = [2]int{curcell[0] + 1, curcell[1]}
			}
			curpos = [2]float64{curpos[0] + hor, curpos[1] + sin*hormult}
			textureIndex = curpos[1] - float64(curcell[1])

		} else {
			if ver < 0 {
				curcell = [2]int{curcell[0], curcell[1] - 1}
			} else {
				curcell = [2]int{curcell[0], curcell[1] + 1}
			}
			curpos = [2]float64{curpos[0] + cos*vermult, curpos[1] + ver}
			textureIndex = curpos[0] - float64(curcell[0])
		}
		player.intersection = curpos
		dist = math.Min(dist+math.Sqrt(math.Min(horlen, verlen)), float64(player.maxDist))
		if curcell[0] < 0 || curcell[0] >= env.cellsx ||
			curcell[1] < 0 || curcell[1] >= env.cellsy {
			return [2]float64{dist, textureIndex}
		}
		if env.cells[curcell[1]][curcell[0]] != 0 {
			return [2]float64{dist, textureIndex}
		}
	}
	return [2]float64{float64(player.maxDist), 0}
}

func (player *Player) draw3D(screen *ebiten.Image, env *Enviroment) {
	//Draw the floor
	twidth := env.textures[0].Bounds().Max.X
	theight := env.textures[0].Bounds().Max.Y
	swidth := screen.Bounds().Max.X
	sheight := screen.Bounds().Max.Y
	rayDirY0 := math.Sin(player.angle - player.fov/2)
	rayDirX0 := math.Cos(player.angle-player.fov/2) * (1 / rayDirY0)
	rayDirY1 := math.Sin(player.angle + player.fov/2)
	rayDirX1 := math.Cos(player.angle+player.fov/2) * (1 / rayDirY1)
	rayDirY0, rayDirY1 = 1, 1
	fmt.Println([]float64{math.Sin(player.angle - player.fov/2), math.Sin(player.angle + player.fov/2)})
	for y := 0; y < sheight/2; y++ {
		rowDist := ((float64(sheight) / 2) / (float64(sheight)/2 - float64(y))) * 1

		floorStepX := rowDist * (rayDirX1 - rayDirX0) / float64(swidth)
		floorStepY := rowDist * (rayDirY1 - rayDirY0) / float64(swidth)

		floorX := player.x + rowDist*rayDirX0
		floorY := player.y + rowDist*rayDirY0

		for x := 0; x < swidth; x++ {
			cellX := int(floorX)
			cellY := int(floorY)

			tx := int(float64(twidth) * (floorX - float64(cellX)))
			ty := int(float64(theight) * (floorY - float64(cellY)))

			floorX += floorStepX
			floorY += floorStepY

			color := env.textures[0].At(tx, ty)
			screen.Set(x, y, color)
			screen.Set(x, sheight-y-1, color)
		}
	}

	//Draw the walls
	for x, ray := range player.rays {
		height := 0
		if ray[0] != 0 {
			height = int(float64(sheight) / ray[0])
		}
		if height != 0 {
			imgx := int(ray[1] * float64(env.textures[1].Bounds().Max.X))
			for y := int(math.Max(float64(sheight-height)/2, 0)); y < int(math.Min(float64((sheight-height)/2+height), float64(sheight))); y++ {
				imgy := int((float64(y-(sheight-height)/2) / float64(height)) * float64(env.textures[1].Bounds().Max.Y))
				screen.Set(x, y, env.textures[1].At(imgx, imgy))
			}
		}
	}
}

func (player *Player) draw2D(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, player.x-5, player.y-5, 10, 10, color.RGBA{255, 255, 255, 255})
	ebitenutil.DrawRect(screen, player.intersection[0]-5, player.intersection[1]-5, 10, 10, color.RGBA{0, 0, 255, 255})
	ebitenutil.DrawLine(screen, player.x, player.y, player.x+math.Cos(player.angle)*player.frontRay, player.y+math.Sin(player.angle)*player.frontRay, color.RGBA{0, 255, 0, 255})
}
