package main

import (
	"image/color"
	"math"
	"sort"

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

	rays []Ray

	maxDist   int
	walkSpeed float64

	frontRay     float64
	intersection [2]float64
}

func (player *Player) init(x float64, y float64, angle float64) {
	player.posX = x
	player.posY = y

	player.dirX, player.dirY = rotate(1, 0, angle)
	player.planeX, player.planeY = rotate(0, 0.5, angle)

	player.maxDist = 20
	player.walkSpeed = 0.03
}

func (player *Player) update(screen *ebiten.Image, env *Enviroment, start bool) {
	changed := false
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		nextX := player.posX + player.dirX*player.walkSpeed
		nextY := player.posY + player.dirY*player.walkSpeed
		var goX bool
		var goY bool
		if round(nextX) == round(player.posX) {
			goX = true
		} else if nextX > 0 && nextX < float64(env.cellsx) {
			if env.cells[int(player.posY)][int(nextX)] == 0 {
				goX = true
			}
		} else if nextX > 0 && nextX < float64(env.cellsx) {
			goX = true
		}
		if round(nextY) == round(player.posY) {
			goY = true
		} else if nextY > 0 && nextY < float64(env.cellsy) {
			if env.cells[int(nextY)][int(player.posX)] == 0 {
				goY = true
			}
		} else if nextY > 0 && nextY < float64(env.cellsy) {
			goY = true
		}
		if goX {
			player.posX += player.dirX * player.walkSpeed
			changed = true
		}
		if goY {
			player.posY += player.dirY * player.walkSpeed
			changed = true
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		nextX := player.posX - player.dirX*player.walkSpeed
		nextY := player.posY - player.dirY*player.walkSpeed
		var goX bool
		var goY bool
		if round(nextX) == round(player.posX) {
			goX = true
		} else if nextX > 0 && nextX < float64(env.cellsx) {
			if env.cells[int(player.posY)][int(nextX)] == 0 {
				goX = true
			}
		} else if nextX > 0 && nextX < float64(env.cellsx) {
			goX = true
		}
		if round(nextY) == round(player.posY) {
			goY = true
		} else if nextY > 0 && nextY < float64(env.cellsy) {
			if env.cells[int(nextY)][int(player.posX)] == 0 {
				goY = true
			}
		} else if nextY > 0 && nextY < float64(env.cellsy) {
			goY = true
		}
		if goX {
			player.posX -= player.dirX * player.walkSpeed
			changed = true
		}
		if goY {
			player.posY -= player.dirY * player.walkSpeed
			changed = true
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		changed = true

		player.dirX, player.dirY = rotate(player.dirX, player.dirY, -0.02)
		player.planeX, player.planeY = rotate(player.planeX, player.planeY, -0.02)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		changed = true

		player.dirX, player.dirY = rotate(player.dirX, player.dirY, 0.02)
		player.planeX, player.planeY = rotate(player.planeX, player.planeY, 0.02)
	}

	if changed || start {
		player.frontRay = player.ray(env, 0).length
		player.rays = []Ray{}
		width := screen.Bounds().Max.X
		for x := 0; x < width; x++ {
			cameraX := (float64(x)/float64(width))*2 - 1
			ray := player.ray(env, cameraX)
			player.rays = append(player.rays, ray)
		}
	}
}

func (player *Player) ray(env *Enviroment, cameraX float64) Ray {
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
		dist = math.Min(dist+math.Sqrt(math.Min(horlen, verlen)), float64(player.maxDist))
		var realDist float64
		if side == 0 {
			realDist = (curpos[0] - player.posX) / rayDirX
		} else {
			realDist = (curpos[1] - player.posY) / rayDirY
		}

		if curcell[0] < 0 || curcell[0] >= env.cellsx ||
			curcell[1] < 0 || curcell[1] >= env.cellsy {
			return Ray{realDist, textureIndex, env.cells[curcell[1]][curcell[0]]}
		}
		if env.cells[curcell[1]][curcell[0]] != 0 {
			return Ray{realDist, textureIndex, env.cells[curcell[1]][curcell[0]]}
		}
	}
	return Ray{float64(player.maxDist), 0, 0}
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

			color := env.textures[6].At(tx, ty)
			screen.Set(x, y, color)
			color = env.textures[3].At(tx, ty)
			screen.Set(x, sheight-y-1, color)
		}
	}

	//Draw the walls
	for x, ray := range player.rays {
		if ray.texIndex != 0 {
			height := 0
			if ray.length != 0 {
				height = int(float64(sheight) / ray.length)
			}
			if height != 0 {
				imgx := int(ray.texIndex * float64(env.textures[ray.texture-1].Bounds().Max.X))
				for y := int(math.Max(float64(sheight-height)/2, 0)); y < int(math.Min(float64((sheight-height)/2+height), float64(sheight))); y++ {
					imgy := int((float64(y-(sheight-height)/2) / float64(height)) * float64(env.textures[1].Bounds().Max.Y))
					screen.Set(x, y, env.textures[ray.texture-1].At(imgx, imgy))
				}
			}
		}

	}

	for i, sprite := range env.sprites {
		env.sprites[i].distance = math.Pow(player.posX-sprite.posX, 2) + math.Pow(player.posY-sprite.posY, 2)
	}
	sort.SliceStable(env.sprites, func(i, j int) bool {
		return env.sprites[i].distance > env.sprites[j].distance
	})
	for _, sprite := range env.sprites {
		relX := sprite.posX - player.posX
		relY := sprite.posY - player.posY

		invDet := 1 / (player.planeX*player.dirY - player.dirX*player.planeY)

		transformX := invDet * (player.dirY*relX - player.dirX*relY)
		transformY := invDet * (-player.planeY*relX + player.planeX*relY)

		screenX := int((float64(swidth) / 2) * (1 + transformX/transformY))
		spriteSize := int(math.Abs(float64(sheight) / transformY))

		for x := int(math.Max(float64(screenX-spriteSize/2), 0)); x < int(math.Min(float64(screenX+spriteSize/2), float64(swidth-1))); x++ {
			imgx := int((float64(x-(screenX-spriteSize/2)) / float64(spriteSize)) * float64(env.textures[sprite.texture].Bounds().Max.X))
			if player.rays[x].length > transformY && transformY > 0 {
				for y := int(math.Max(float64((sheight-spriteSize)/2), 0)); y < int(math.Min(float64((sheight-spriteSize)/2+spriteSize), float64(sheight))); y++ {
					imgy := int((float64(y-(sheight-spriteSize)/2) / float64(spriteSize)) * float64(env.textures[sprite.texture].Bounds().Max.Y))
					r, g, b, a := env.textures[sprite.texture].At(imgx, imgy).RGBA()
					if a != 0 {
						screen.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
					}
					//fmt.Println(col)
				}
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

func round(number float64) int {
	if number < 0 {
		number--
	}
	return int(number)
}

func rotate(x, y, a float64) (float64, float64) {
	return x*math.Cos(a) - y*math.Sin(a), y*math.Cos(a) + x*math.Sin(a)
}

//Ray is a struct for rays
type Ray struct {
	length   float64
	texIndex float64
	texture  int
}
