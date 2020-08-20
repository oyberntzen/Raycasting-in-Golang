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

	pitch float64
	posZ  float64
	velZ  float64

	rays []Ray

	maxDist int

	walkSpeed float64
	rotXSpeed float64
	rotYSpeed float64

	frontRay float64

	cursorX float64
	cursorY float64

	screenWidth  int
	screenHeight int

	size float64
}

func (player *Player) init(x float64, y float64, angle float64, screenWidth int, screenHeight int) {
	player.posX = x
	player.posY = y
	player.screenWidth = screenWidth
	player.screenHeight = screenHeight

	player.dirX, player.dirY = rotate(1, 0, angle)
	player.planeX, player.planeY = rotate(0, 0.5*(float64(player.screenWidth)/float64(player.screenHeight)), angle)

	player.maxDist = 20

	player.walkSpeed = 0.05
	sensitivity := 1.0
	player.rotXSpeed = 1.5 * sensitivity
	player.rotYSpeed = 2.5 * sensitivity

	player.pitch = 0
	player.posZ = 0

	cursorX, cursorY := ebiten.CursorPosition()
	player.cursorX = float64(cursorX)
	player.cursorY = float64(cursorY)

	player.size = 0.1

}

func (player *Player) update(screen *ebiten.Image, env *Enviroment) {
	player.handleInput(screen, env)

	player.velZ -= 0.003
	player.posZ = math.Min(math.Max(player.posZ+player.velZ, 0), float64(player.screenHeight)*0.001)

	player.frontRay = player.ray(env, 0).length

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
	player.rayCast(screen, env)
	player.drawFloorCeiling(screen, env)
	player.drawWalls(screen, env)
	player.drawSprites(screen, env)

}

func (player *Player) rayCast(screen *ebiten.Image, env *Enviroment) {
	player.rays = []Ray{}
	for x := 0; x < player.screenWidth; x++ {
		cameraX := (float64(x)/float64(player.screenWidth))*2 - 1
		ray := player.ray(env, cameraX)
		player.rays = append(player.rays, ray)
	}
}

func (player *Player) drawFloorCeiling(screen *ebiten.Image, env *Enviroment) {
	twidth := env.textures[0].Bounds().Max.X
	theight := env.textures[0].Bounds().Max.Y
	rayDirX0 := player.dirX - player.planeX
	rayDirY0 := player.dirY - player.planeY
	rayDirX1 := player.dirX + player.planeX
	rayDirY1 := player.dirY + player.planeY

	for y := 0; y < player.screenHeight; y++ {
		isFloor := y > player.screenHeight/2+int(player.pitch)

		p := float64(player.screenHeight)/2 - float64(y) + player.pitch
		if isFloor {
			p = float64(y) - float64(player.screenHeight)/2 - player.pitch
		}

		camZ := (float64(player.screenHeight) / 2) - player.posZ*float64(player.screenHeight)
		if isFloor {
			camZ = (float64(player.screenHeight) / 2) + player.posZ*float64(player.screenHeight)
		}
		rowDist := camZ / p

		floorStepX := rowDist * (rayDirX1 - rayDirX0) / float64(player.screenWidth)
		floorStepY := rowDist * (rayDirY1 - rayDirY0) / float64(player.screenWidth)

		floorX := player.posX + rowDist*rayDirX0
		floorY := player.posY + rowDist*rayDirY0

		for x := 0; x < player.screenWidth; x++ {
			cellX := int(floorX)
			cellY := int(floorY)

			tx := int(float64(twidth) * (floorX - float64(cellX)))
			ty := int(float64(theight) * (floorY - float64(cellY)))

			floorX += floorStepX
			floorY += floorStepY

			color := env.textures[6].At(tx, ty)
			if isFloor {
				color = env.textures[3].At(tx, ty)
			}
			screen.Set(x, y, color)
		}
	}
}

func (player *Player) drawWalls(screen *ebiten.Image, env *Enviroment) {
	for x, ray := range player.rays {
		if ray.texIndex != 0 {
			height := 0
			if ray.length != 0 {
				height = int(float64(player.screenHeight) / ray.length)
			}
			if height != 0 {
				imgx := int(ray.texIndex * float64(env.textures[ray.texture-1].Bounds().Max.X))
				start := int(math.Max(float64(player.screenHeight-height)/2+player.pitch+player.posZ*float64(player.screenHeight)/ray.length, 0))
				end := int(math.Min(float64((player.screenHeight-height)/2+height)+player.pitch+player.posZ*float64(player.screenHeight)/ray.length, float64(player.screenHeight)))
				for y := start; y < end; y++ {
					imgy := int(((float64(y-(player.screenHeight-height)/2) - player.pitch - player.posZ*float64(player.screenHeight)/ray.length) / float64(height)) * float64(env.textures[1].Bounds().Max.Y))
					screen.Set(x, y, env.textures[ray.texture-1].At(imgx, imgy))
				}
			}
		}

	}
}

func (player *Player) drawSprites(screen *ebiten.Image, env *Enviroment) {
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

		screenX := int((float64(player.screenWidth) / 2) * (1 + transformX/transformY))
		spriteSize := int(math.Abs(float64(player.screenHeight) / transformY))

		for x := int(math.Max(float64(screenX-spriteSize/2), 0)); x < int(math.Min(float64(screenX+spriteSize/2), float64(player.screenWidth-1))); x++ {
			imgx := int((float64(x-(screenX-spriteSize/2)) / float64(spriteSize)) * float64(env.textures[sprite.texture].Bounds().Max.X))
			if player.rays[x].length > transformY && transformY > 0 {
				start := int(math.Max(float64((player.screenHeight-spriteSize)/2)+player.pitch+player.posZ*float64(player.screenHeight)/transformY, 0))
				end := int(math.Min(float64((player.screenHeight-spriteSize)/2+spriteSize)+player.pitch+player.posZ*float64(player.screenHeight)/transformY, float64(player.screenHeight)))
				for y := start; y < end; y++ {
					imgy := int(((float64(y-(player.screenHeight-spriteSize)/2) - player.pitch - player.posZ*float64(player.screenHeight)/transformY) / float64(spriteSize)) * float64(env.textures[sprite.texture].Bounds().Max.Y))
					r, g, b, a := env.textures[sprite.texture].At(imgx, imgy).RGBA()
					if a != 0 {
						screen.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
					}
				}
			}

		}
	}
}

func (player *Player) draw2D(screen *ebiten.Image, env *Enviroment) {
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

func (player *Player) handleInput(screen *ebiten.Image, env *Enviroment) {
	nextX, nextY := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		nextX += player.dirX
		nextY += player.dirY

	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		nextX -= player.dirX
		nextY -= player.dirY

	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dirX, dirY := rotate(player.dirX, player.dirY, -math.Pi/2)
		nextX += dirX
		nextY += dirY
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dirX, dirY := rotate(player.dirX, player.dirY, math.Pi/2)
		nextX += dirX
		nextY += dirY

	}
	angle := math.Atan2(nextY, nextX)
	if nextX != 0 {
		nextX = math.Cos(angle)
	}
	if nextY != 0 {
		nextY = math.Sin(angle)
	}
	//nextX, nextY = math.Cos(angle), math.Sin(angle)
	goX, goY := player.collision(env, player.posX+nextX*player.walkSpeed, player.posY+nextY*player.walkSpeed)
	if goX {
		player.posX += nextX * player.walkSpeed
	}
	if goY {
		player.posY += nextY * player.walkSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && player.posZ == 0 {
		player.velZ = 0.045
	}
	cursorX, cursorY := ebiten.CursorPosition()
	newCursorX := float64(cursorX)
	newCursorY := float64(cursorY)

	changeX := (newCursorX - player.cursorX) / float64(player.screenWidth)
	changeY := (newCursorY - player.cursorY)

	if changeX > -1000 && changeX < 1000 {
		player.dirX, player.dirY = rotate(player.dirX, player.dirY, changeX*player.rotYSpeed)
		player.planeX, player.planeY = rotate(player.planeX, player.planeY, changeX*player.rotYSpeed)
		player.pitch = math.Max(math.Min(player.pitch-changeY*player.rotXSpeed, float64(player.screenHeight)), float64(-player.screenHeight))
	}
	player.cursorX, player.cursorY = newCursorX, newCursorY
}

func (player *Player) collision(env *Enviroment, nextX, nextY float64) (bool, bool) {
	return doCollide(player.posX, nextX+player.size, player.posY+player.size, true, env) &&
			doCollide(player.posX, nextX+player.size, player.posY-player.size, true, env) &&
			doCollide(player.posX, nextX-player.size, player.posY+player.size, true, env) &&
			doCollide(player.posX, nextX-player.size, player.posY-player.size, true, env),
		doCollide(player.posY, nextY+player.size, player.posX+player.size, false, env) &&
			doCollide(player.posY, nextY+player.size, player.posX-player.size, false, env) &&
			doCollide(player.posY, nextY-player.size, player.posX+player.size, false, env) &&
			doCollide(player.posY, nextY-player.size, player.posX-player.size, false, env)
}

func doCollide(start, end, other float64, isX bool, env *Enviroment) bool {
	if round(end) == round(start) {
		return true
	} else if end > 0 && end < float64(env.cellsx) {
		if isX {
			if env.cells[int(other)][int(end)] == 0 {
				return true
			}
		} else {
			if env.cells[int(end)][int(other)] == 0 {
				return true
			}
		}
	} else if end > 0 && end < float64(env.cellsx) {
		return true
	}
	return false
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
