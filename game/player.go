package game

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

//Player is the player in the game
type Player struct {
	PosX   float64
	PosY   float64
	DirX   float64
	DirY   float64
	PlaneX float64
	PlaneY float64

	Pitch float64
	PosZ  float64
	velZ  float64

	rays []Ray

	maxDist int

	walkSpeed float64
	rotXSpeed float64
	rotYSpeed float64

	frontRay float64

	cursorX float64
	cursorY float64

	ScreenWidth  int
	ScreenHeight int

	size float64
}

//Init initializes the player with the specified x, y and angle. Width and height of screen is also required for some calculations
func (player *Player) Init(x float64, y float64, angle float64, screenWidth int, screenHeight int) {
	player.PosX = x
	player.PosY = y
	player.ScreenWidth = screenWidth
	player.ScreenHeight = screenHeight

	player.DirX, player.DirY = rotate(1, 0, angle)
	player.PlaneX, player.PlaneY = rotate(0, 0.5*(float64(player.ScreenWidth)/float64(player.ScreenHeight)), angle)

	player.Pitch = 0
	player.PosZ = 0

	player.InitConstants()
}

//InitConstants initializes all unexported variables
func (player *Player) InitConstants() {
	player.maxDist = 20

	player.walkSpeed = 0.05
	sensitivity := 1.0
	player.rotXSpeed = 1.5 * sensitivity
	player.rotYSpeed = 2.5 * sensitivity
	cursorX, cursorY := ebiten.CursorPosition()
	player.cursorX = float64(cursorX)
	player.cursorY = float64(cursorY)

	player.size = 0.1
}

//Update updates the player's position based on user input
func (player *Player) Update(screen *ebiten.Image, env *Enviroment) {
	player.handleInput(screen, env)

	player.velZ -= 0.003
	player.PosZ = math.Min(math.Max(player.PosZ+player.velZ, 0), float64(player.ScreenHeight)*0.001)

	player.frontRay = player.ray(env, 0).length

}

func (player *Player) ray(env *Enviroment, cameraX float64) Ray {
	dist := float64(0)

	rayDirX := player.DirX + player.PlaneX*cameraX
	rayDirY := player.DirY + player.PlaneY*cameraX

	curcell := [2]int{int(player.PosX), int(player.PosY)}
	curpos := [2]float64{player.PosX, player.PosY}

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
			realDist = (curpos[0] - player.PosX) / rayDirX
		} else {
			realDist = (curpos[1] - player.PosY) / rayDirY
		}

		if curcell[0] < 0 || curcell[0] >= env.Cellsx ||
			curcell[1] < 0 || curcell[1] >= env.Cellsy {
			return Ray{realDist, textureIndex, env.Cells[curcell[1]][curcell[0]]}
		}
		if env.Cells[curcell[1]][curcell[0]] != 0 {
			return Ray{realDist, textureIndex, env.Cells[curcell[1]][curcell[0]]}
		}
	}
	return Ray{float64(player.maxDist), 0, 0}
}

//Draw3D draws walls, floor, ceiling and sprites from the first person view of the player
func (player *Player) Draw3D(screen *ebiten.Image, env *Enviroment, sprites []Sprite) {
	player.rayCast(screen, env)
	player.drawFloorCeiling(screen, env)
	player.drawWalls(screen, env)
	player.drawSprites(screen, env, append(env.Sprites, sprites...))

}

func (player *Player) rayCast(screen *ebiten.Image, env *Enviroment) {
	player.rays = []Ray{}
	for x := 0; x < player.ScreenWidth; x++ {
		cameraX := (float64(x)/float64(player.ScreenWidth))*2 - 1
		ray := player.ray(env, cameraX)
		player.rays = append(player.rays, ray)
	}
}

func (player *Player) drawFloorCeiling(screen *ebiten.Image, env *Enviroment) {
	twidth := env.textures[0].Bounds().Max.X
	theight := env.textures[0].Bounds().Max.Y
	rayDirX0 := player.DirX - player.PlaneX
	rayDirY0 := player.DirY - player.PlaneY
	rayDirX1 := player.DirX + player.PlaneX
	rayDirY1 := player.DirY + player.PlaneY

	for y := 0; y < player.ScreenHeight; y++ {
		isFloor := y > player.ScreenHeight/2+int(player.Pitch)

		p := float64(player.ScreenHeight)/2 - float64(y) + player.Pitch
		if isFloor {
			p = float64(y) - float64(player.ScreenHeight)/2 - player.Pitch
		}

		camZ := (float64(player.ScreenHeight) / 2) - player.PosZ*float64(player.ScreenHeight)
		if isFloor {
			camZ = (float64(player.ScreenHeight) / 2) + player.PosZ*float64(player.ScreenHeight)
		}
		rowDist := camZ / p

		floorStepX := rowDist * (rayDirX1 - rayDirX0) / float64(player.ScreenWidth)
		floorStepY := rowDist * (rayDirY1 - rayDirY0) / float64(player.ScreenWidth)

		floorX := player.PosX + rowDist*rayDirX0
		floorY := player.PosY + rowDist*rayDirY0

		for x := 0; x < player.ScreenWidth; x++ {
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
				height = int(float64(player.ScreenHeight) / ray.length)
			}
			if height != 0 {
				imgx := int(ray.texIndex * float64(env.textures[ray.texture-1].Bounds().Max.X))
				start := int(math.Max(float64(player.ScreenHeight-height)/2+player.Pitch+player.PosZ*float64(player.ScreenHeight)/ray.length, 0))
				end := int(math.Min(float64((player.ScreenHeight-height)/2+height)+player.Pitch+player.PosZ*float64(player.ScreenHeight)/ray.length, float64(player.ScreenHeight)))
				for y := start; y < end; y++ {
					imgy := int(((float64(y-(player.ScreenHeight-height)/2) - player.Pitch - player.PosZ*float64(player.ScreenHeight)/ray.length) / float64(height)) * float64(env.textures[1].Bounds().Max.Y))
					screen.Set(x, y, env.textures[ray.texture-1].At(imgx, imgy))
				}
			}
		}

	}
}

func (player *Player) drawSprites(screen *ebiten.Image, env *Enviroment, sprites []Sprite) {
	for i, sprite := range sprites {
		sprites[i].distance = math.Pow(player.PosX-sprite.PosX, 2) + math.Pow(player.PosY-sprite.PosY, 2)
	}
	sort.SliceStable(sprites, func(i, j int) bool {
		return sprites[i].distance > sprites[j].distance
	})
	for _, sprite := range sprites {
		relX := sprite.PosX - player.PosX
		relY := sprite.PosY - player.PosY

		invDet := 1 / (player.PlaneX*player.DirY - player.DirX*player.PlaneY)

		transformX := invDet * (player.DirY*relX - player.DirX*relY)
		transformY := invDet * (-player.PlaneY*relX + player.PlaneX*relY)

		screenX := int((float64(player.ScreenWidth) / 2) * (1 + transformX/transformY))
		spriteSize := int(math.Abs(float64(player.ScreenHeight) / transformY))

		for x := int(math.Max(float64(screenX-spriteSize/2), 0)); x < int(math.Min(float64(screenX+spriteSize/2), float64(player.ScreenWidth-1))); x++ {
			imgx := int((float64(x-(screenX-spriteSize/2)) / float64(spriteSize)) * float64(env.textures[sprite.Texture].Bounds().Max.X))
			if player.rays[x].length > transformY && transformY > 0 {
				start := int(math.Max(float64((player.ScreenHeight-spriteSize)/2)+player.Pitch+player.PosZ*float64(player.ScreenHeight)/transformY, 0))
				end := int(math.Min(float64((player.ScreenHeight-spriteSize)/2+spriteSize)+player.Pitch+player.PosZ*float64(player.ScreenHeight)/transformY, float64(player.ScreenHeight)))
				for y := start; y < end; y++ {
					imgy := int(((float64(y-(player.ScreenHeight-spriteSize)/2) - player.Pitch - player.PosZ*float64(player.ScreenHeight)/transformY) / float64(spriteSize)) * float64(env.textures[sprite.Texture].Bounds().Max.Y))
					r, g, b, a := env.textures[sprite.Texture].At(imgx, imgy).RGBA()
					if a != 0 {
						screen.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
					}
				}
			}

		}
	}
}

//Draw2D draws the player's position and direction from top-down view
func (player *Player) Draw2D(screen *ebiten.Image, env *Enviroment) {
	cellsizex := float64(screen.Bounds().Max.X / env.Cellsx)
	cellsizey := float64(screen.Bounds().Max.Y / env.Cellsy)
	ebitenutil.DrawRect(screen, player.PosX*cellsizex-5, player.PosY*cellsizey-5, 10, 10, color.RGBA{255, 255, 255, 255})
	//ebitenutil.DrawRect(screen, player.intersection[0]*cellsizex-5, player.intersection[1]*cellsizey-5, 10, 10, color.RGBA{0, 0, 255, 255})
	ebitenutil.DrawLine(screen, player.PosX*cellsizex, player.PosY*cellsizey, (player.PosX+player.DirX*player.frontRay)*cellsizex, (player.PosY+player.DirY*player.frontRay)*cellsizey, color.RGBA{0, 255, 0, 255})

	rayDirX0 := player.DirX - player.PlaneX
	rayDirY0 := player.DirY - player.PlaneY
	rayDirX1 := player.DirX + player.PlaneX
	rayDirY1 := player.DirY + player.PlaneY
	ebitenutil.DrawLine(screen, player.PosX*cellsizex, player.PosY*cellsizey, (player.PosX+rayDirX0)*cellsizex, (player.PosY+rayDirY0)*cellsizey, color.RGBA{0, 255, 255, 255})
	ebitenutil.DrawLine(screen, player.PosX*cellsizex, player.PosY*cellsizey, (player.PosX+rayDirX1)*cellsizex, (player.PosY+rayDirY1)*cellsizey, color.RGBA{0, 255, 255, 255})
	ebitenutil.DrawLine(screen, player.PosX*cellsizex, player.PosY*cellsizey, (player.PosX+player.DirX)*cellsizex, (player.PosY+player.DirY)*cellsizey, color.RGBA{255, 0, 255, 255})
}

func (player *Player) handleInput(screen *ebiten.Image, env *Enviroment) {
	nextX, nextY := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		nextX += player.DirX
		nextY += player.DirY

	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		nextX -= player.DirX
		nextY -= player.DirY

	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dirX, dirY := rotate(player.DirX, player.DirY, -math.Pi/2)
		nextX += dirX
		nextY += dirY
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dirX, dirY := rotate(player.DirX, player.DirY, math.Pi/2)
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
	goX, goY := player.collision(env, player.PosX+nextX*player.walkSpeed, player.PosY+nextY*player.walkSpeed)
	if goX {
		player.PosX += nextX * player.walkSpeed
	}
	if goY {
		player.PosY += nextY * player.walkSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && player.PosZ == 0 {
		player.velZ = 0.045
	}
	cursorX, cursorY := ebiten.CursorPosition()
	newCursorX := float64(cursorX)
	newCursorY := float64(cursorY)

	changeX := (newCursorX - player.cursorX) / float64(player.ScreenWidth)
	changeY := (newCursorY - player.cursorY)

	if changeX > -1000 && changeX < 1000 {
		player.DirX, player.DirY = rotate(player.DirX, player.DirY, changeX*player.rotYSpeed)
		player.PlaneX, player.PlaneY = rotate(player.PlaneX, player.PlaneY, changeX*player.rotYSpeed)
		player.Pitch = math.Max(math.Min(player.Pitch-changeY*player.rotXSpeed, float64(player.ScreenHeight)), float64(-player.ScreenHeight))
	}
	player.cursorX, player.cursorY = newCursorX, newCursorY
}

func (player *Player) handleInput2D(screen *ebiten.Image, env *Enviroment) {
	nextX, nextY := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {

		nextX += player.DirX
		nextY += player.DirY
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		nextX -= player.DirX
		nextY -= player.DirY
	}
	fmt.Println(nextX, nextY)

	angle := math.Atan2(nextY, nextX)
	if nextX != 0 {
		nextX = math.Cos(angle)
	}
	if nextY != 0 {
		nextY = math.Sin(angle)
	}
	goX, goY := player.collision(env, player.PosX+nextX*player.walkSpeed, player.PosY+nextY*player.walkSpeed)
	fmt.Println(goX, goY, nextX*player.walkSpeed, nextY*player.walkSpeed, player.PosY, player.PosY+nextY*player.walkSpeed)
	if goX {
		player.PosX += nextX * player.walkSpeed
	}
	if goY {
		player.PosY += nextY * player.walkSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		player.DirX, player.DirY = rotate(player.DirX, player.DirY, -0.02)
		player.PlaneX, player.PlaneY = rotate(player.PlaneX, player.PlaneY, -0.02)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		player.DirX, player.DirY = rotate(player.DirX, player.DirY, 0.02)
		player.PlaneX, player.PlaneY = rotate(player.PlaneX, player.PlaneY, 0.02)
	}
}

func (player *Player) collision(env *Enviroment, nextX, nextY float64) (bool, bool) {
	return doCollide(player.PosX, nextX+player.size, player.PosY+player.size, true, env) &&
			doCollide(player.PosX, nextX+player.size, player.PosY-player.size, true, env) &&
			doCollide(player.PosX, nextX-player.size, player.PosY+player.size, true, env) &&
			doCollide(player.PosX, nextX-player.size, player.PosY-player.size, true, env),
		doCollide(player.PosY, nextY+player.size, player.PosX+player.size, false, env) &&
			doCollide(player.PosY, nextY+player.size, player.PosX-player.size, false, env) &&
			doCollide(player.PosY, nextY-player.size, player.PosX+player.size, false, env) &&
			doCollide(player.PosY, nextY-player.size, player.PosX-player.size, false, env)
}

func doCollide(start, end, other float64, isX bool, env *Enviroment) bool {
	if round(end) == round(start) {
		return true
	} else if end > 0 && end < float64(env.Cellsx) {
		if isX {
			if env.Cells[int(other)][int(end)] == 0 {
				return true
			}
		} else {
			if env.Cells[int(end)][int(other)] == 0 {
				return true
			}
		}
	} else if end > 0 && end < float64(env.Cellsx) {
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
