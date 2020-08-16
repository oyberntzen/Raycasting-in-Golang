package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

//Player is the player in the game
type Player struct {
	posX   float64
	posY   float64
	dirX   float64
	dirY   float64
	planeX float64
	planeY float64

	rays [][2]float64

	maxDist   int
	walkSpeed float64

	frontRay     float64
	intersection [2]float64
}

func (player *Player) init(ratio float64, x float64, y float64, angle float64) {
	player.posX = x
	player.posY = y

	player.dirX = math.Cos(angle)
	player.dirY = math.Sin(angle)

	player.planeX = 0
	player.planeY = 0.5

	player.maxDist = 20
	player.walkSpeed = 0.03
}

func (player *Player) update(screen *ebiten.Image, env *Enviroment, start bool) {
	changed := false
	forwardX := player.posX + player.dirX*player.walkSpeed
	forwardY := player.posY + player.dirY*player.walkSpeed
	if ebiten.IsKeyPressed(ebiten.KeyW) && env.cells[int(math.Max(0, math.Min(float64(env.cellsx-1), forwardX)))][int(math.Max(0, math.Min(float64(env.cellsx-1), forwardY)))] == 0 {
		changed = true
		player.posX += player.dirX * player.walkSpeed
		player.posY += player.dirY * player.walkSpeed
	}
	backwardX := player.posX - player.dirX*player.walkSpeed
	backwardY := player.posY - player.dirY*player.walkSpeed
	if ebiten.IsKeyPressed(ebiten.KeyS) && env.cells[int(math.Max(0, math.Min(float64(env.cellsx-1), backwardX)))][int(math.Max(0, math.Min(float64(env.cellsx-1), backwardY)))] == 0 {
		changed = true
		player.posX -= player.dirX * player.walkSpeed
		player.posY -= player.dirY * player.walkSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		changed = true

		oldDirX := player.dirX
		player.dirX = player.dirX*math.Cos(-0.02) - player.dirY*math.Sin(-0.02)
		player.dirY = player.dirY*math.Cos(-0.02) + oldDirX*math.Sin(-0.02)

		oldPlaneX := player.planeX
		player.planeX = player.planeX*math.Cos(-0.02) - player.planeY*math.Sin(-0.02)
		player.planeY = player.planeY*math.Cos(-0.02) + oldPlaneX*math.Sin(-0.02)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		changed = true

		oldDirX := player.dirX
		player.dirX = player.dirX*math.Cos(0.02) - player.dirY*math.Sin(0.02)
		player.dirY = player.dirY*math.Cos(0.02) + oldDirX*math.Sin(0.02)

		oldPlaneX := player.planeX
		player.planeX = player.planeX*math.Cos(0.02) - player.planeY*math.Sin(0.02)
		player.planeY = player.planeY*math.Cos(0.02) + oldPlaneX*math.Sin(0.02)
	}

	player.posX = math.Max(0.001, math.Min(float64(env.cellsx)-0.001, player.posX))
	player.posY = math.Max(0.001, math.Min(float64(env.cellsy)-0.001, player.posY))

	if changed || start {
		player.frontRay = player.ray(env, 0)[0]
		player.rays = [][2]float64{}
		width := screen.Bounds().Max.X
		for x := 0; x < width; x++ {
			cameraX := (float64(x)/float64(width))*2 - 1
			ray := player.ray(env, cameraX)
			//real := ray[0] / math.Sqrt(cameraX*cameraX+1)
			//fmt.Println([2]float64{real, ray[0]})
			player.rays = append(player.rays, ray) //[2]float64{real, ray[1]})
		}
	}
}

func (player *Player) ray(env *Enviroment, cameraX float64) [2]float64 {
	dist := float64(0)

	rayDirX := player.dirX + player.planeX*cameraX
	rayDirY := player.dirY + player.planeY*cameraX

	curcell := [2]int{int(player.posX), int(player.posY)}
	curpos := [2]float64{player.posX, player.posY}

	side := 0

	for int(dist) < player.maxDist {

		relx := curpos[0] - float64(curcell[0])
		rely := curpos[1] - float64(curcell[1])

		hor := 1 - relx
		if rayDirX <= 0 {
			hor = -relx
		}
		ver := 1 - rely
		if rayDirY <= 0 {
			ver = -rely
		}

		var hormult float64
		if rayDirX != 0 {
			hormult = hor / rayDirX
		} else {
			hormult = 1000
		}
		var vermult float64
		if rayDirY != 0 {
			vermult = ver / rayDirY
		} else {
			vermult = 1000
		}

		horlen := math.Pow(hor, 2) + math.Pow(rayDirY*hormult, 2)
		verlen := math.Pow(ver, 2) + math.Pow(rayDirX*vermult, 2)

		var textureIndex float64
		if horlen < verlen {
			side = 0
			if hor < 0 {
				curcell = [2]int{curcell[0] - 1, curcell[1]}
			} else {
				curcell = [2]int{curcell[0] + 1, curcell[1]}
			}
			curpos = [2]float64{curpos[0] + hor, curpos[1] + rayDirY*hormult}
			textureIndex = curpos[1] - float64(curcell[1])

		} else {
			side = 1
			if ver < 0 {
				curcell = [2]int{curcell[0], curcell[1] - 1}
			} else {
				curcell = [2]int{curcell[0], curcell[1] + 1}
			}
			curpos = [2]float64{curpos[0] + rayDirX*vermult, curpos[1] + ver}
			textureIndex = curpos[0] - float64(curcell[0])
		}
		//player.intersection = curpos
		dist = math.Min(dist+math.Sqrt(math.Min(horlen, verlen)), float64(player.maxDist))
		var realDist float64
		if side == 0 {
			realDist = (curpos[0] - player.posX) / rayDirX
		} else {
			realDist = (curpos[1] - player.posY) / rayDirY
		}

		if curcell[0] < 0 || curcell[0] >= env.cellsx ||
			curcell[1] < 0 || curcell[1] >= env.cellsy {
			return [2]float64{realDist, textureIndex}
		}
		if env.cells[curcell[1]][curcell[0]] != 0 {
			return [2]float64{realDist, textureIndex}
		}
	}
	return [2]float64{float64(player.maxDist), 0}
}

func (player *Player) draw3D(screen *ebiten.Image, env *Enviroment) {
	//Draw the floor and ceiling
	twidth := env.textures[0].Bounds().Max.X
	theight := env.textures[0].Bounds().Max.Y
	swidth := screen.Bounds().Max.X
	sheight := screen.Bounds().Max.Y
	rayDirX0 := player.dirX - player.planeX
	rayDirY0 := player.dirY - player.planeY
	rayDirX1 := player.dirX + player.planeX
	rayDirY1 := player.dirY + player.planeY

	for y := 0; y < sheight/2; y++ {
		rowDist := ((float64(sheight) / 2) / (float64(sheight)/2 - float64(y)))

		floorStepX := rowDist * (rayDirX1 - rayDirX0) / float64(swidth)
		floorStepY := rowDist * (rayDirY1 - rayDirY0) / float64(swidth)

		floorX := player.posX + rowDist*rayDirX0
		floorY := player.posY + rowDist*rayDirY0

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
	cellsizex := float64(screen.Bounds().Max.X / env.cellsx)
	cellsizey := float64(screen.Bounds().Max.Y / env.cellsy)
	ebitenutil.DrawRect(screen, player.posX*cellsizex-5, player.posY*cellsizey-5, 10, 10, color.RGBA{255, 255, 255, 255})
	//ebitenutil.DrawRect(screen, player.intersection[0]*cellsizex-5, player.intersection[1]*cellsizey-5, 10, 10, color.RGBA{0, 0, 255, 255})
	ebitenutil.DrawLine(screen, player.posX*cellsizex, player.posY*cellsizey, (player.posX+player.dirX*player.frontRay)*cellsizex, (player.posY+player.dirY*player.frontRay)*cellsizey, color.RGBA{0, 255, 0, 255})

	rayDirX0 := player.dirX - player.planeX
	rayDirY0 := player.dirY - player.planeY
	rayDirX1 := player.dirX + player.planeX
	rayDirY1 := player.dirY + player.planeY
	ebitenutil.DrawLine(screen, player.posX*cellsizex, player.posY*cellsizey, (player.posX+rayDirX0)*cellsizex, (player.posY+rayDirY0)*cellsizey, color.RGBA{0, 255, 255, 255})
	ebitenutil.DrawLine(screen, player.posX*cellsizex, player.posY*cellsizey, (player.posX+rayDirX1)*cellsizex, (player.posY+rayDirY1)*cellsizey, color.RGBA{0, 255, 255, 255})
	ebitenutil.DrawLine(screen, player.posX*cellsizex, player.posY*cellsizey, (player.posX+player.dirX)*cellsizex, (player.posY+player.dirY)*cellsizey, color.RGBA{255, 0, 255, 255})
}
