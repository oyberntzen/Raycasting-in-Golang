package graphics

import (
	"image"
	"image/color"

	//Used to decode the PNG textures
	_ "image/png"
	"math"
	"os"
	"sort"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/nfnt/resize"

	"github.com/oyberntzen/Raycasting-in-Golang/game"
	"github.com/oyberntzen/Raycasting-in-Golang/game/levels"
)

var (
	textures []image.Image
	images   []ebiten.Image
)

//Init inits all textures
func Init(dir string, width, height int) {
	textures = []image.Image{
		load(dir + "brick_circle.png"),
		load(dir + "brick.png"),
		load(dir + "purplestone.png"),
		load(dir + "greystone.png"),
		load(dir + "bluestone.png"),
		load(dir + "stone.png"),
		load(dir + "planks.png"),
		load(dir + "colorstone.png"),
		load(dir + "barrel.png"),
		load(dir + "pillar.png"),
		load(dir + "greenlight.png"),
		load(dir + "player.png"),
		load(dir + "pistol.png"),
		load(dir + "cursor.png"),
	}
	images = append(images, *resizeImage(textures[13], int(float64(width)/50)+1, 0)) //Cursor
	images = append(images, *resizeImage(textures[12], int(float64(width)/5)+1, 0))  //Pistol
}

//Draw3D draws walls, floor, ceiling and sprites from the first person view of the player
func Draw3D(screen *ebiten.Image, player game.Player, cells [][]uint8, sprites []game.Sprite, players []game.Player, width, height int, playerSize float64) {
	dirX, dirY := game.Rotate(1, 0, player.Angle)
	planeX, planeY := game.Rotate(0, 0.5*(float64(width)/float64(height)), player.Angle)

	for _, p := range players {
		sprites = append(sprites, levels.CreateSprite(levels.PlayerInfo, p.X, p.Y, p.Z, playerSize, 0, levels.SpriteZFree))
	}

	dists, indicies, texs := rayCast(player, cells, width, dirX, dirY, planeX, planeY)
	drawFloorCeiling(screen, player, width, height, dirX, dirY, planeX, planeY)
	drawWalls(screen, player, dists, indicies, texs, height)
	drawSprites(screen, player, sprites, dists, dirX, dirY, planeX, planeY, width, height)
	drawUI(screen, width, height)
}

//Ray shoots ray from player and calculates distance to wall
func Ray(player game.Player, cells [][]uint8, rayDirX, rayDirY float64) (float64, float64, uint8) {
	dist := float64(0)

	curcell := [2]int{int(player.X), int(player.Y)}
	curpos := [2]float64{player.X, player.Y}

	side := 0

	for {

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
		dist = dist + math.Sqrt(math.Min(horlen, verlen))
		var realDist float64
		if side == 0 {
			realDist = (curpos[0] - player.X) / rayDirX
		} else {
			realDist = (curpos[1] - player.Y) / rayDirY
		}

		if curcell[0] < 0 || curcell[0] >= len(cells[0]) ||
			curcell[1] < 0 || curcell[1] >= len(cells) {
			return realDist, textureIndex, 1
		}
		if cells[curcell[1]][curcell[0]] != 0 {
			return realDist, textureIndex, cells[curcell[1]][curcell[0]]
		}
	}
}

func rayCast(player game.Player, cells [][]uint8, width int, dirX, dirY, planeX, planeY float64) ([]float64, []float64, []uint8) {
	dists := []float64{}
	indicies := []float64{}
	texs := []uint8{}

	for x := 0; x < width; x++ {
		cameraX := (float64(x)/float64(width))*2 - 1
		rayDirX := dirX + planeX*cameraX
		rayDirY := dirY + planeY*cameraX

		dist, index, texture := Ray(player, cells, rayDirX, rayDirY)
		dists = append(dists, dist)
		indicies = append(indicies, index)
		texs = append(texs, texture)
	}

	return dists, indicies, texs
}

func drawFloorCeiling(screen *ebiten.Image, player game.Player, width, height int, dirX, dirY, planeX, planeY float64) {
	twidthf := textures[3].Bounds().Max.X
	theightf := textures[3].Bounds().Max.Y
	twidthc := textures[6].Bounds().Max.X
	theightc := textures[6].Bounds().Max.Y

	rayDirX0 := dirX - planeX
	rayDirY0 := dirY - planeY
	rayDirX1 := dirX + planeX
	rayDirY1 := dirY + planeY

	pitch := player.Pitch * float64(height)
	for y := 0; y < height; y++ {
		isFloor := y > height/2+int(pitch)

		p := float64(height)/2 - float64(y) + pitch
		if isFloor {
			p = float64(y) - float64(height)/2 - pitch
		}

		camZ := (float64(height) / 2) - player.Z*float64(height)
		if isFloor {
			camZ = (float64(height) / 2) + player.Z*float64(height)
		}
		rowDist := camZ / p

		floorStepX := rowDist * (rayDirX1 - rayDirX0) / float64(width)
		floorStepY := rowDist * (rayDirY1 - rayDirY0) / float64(width)

		floorX := player.X + rowDist*rayDirX0
		floorY := player.Y + rowDist*rayDirY0

		for x := 0; x < width; x++ {
			cellX := int(floorX)
			cellY := int(floorY)

			txf := int(float64(twidthf) * (floorX - float64(cellX)))
			tyf := int(float64(theightf) * (floorY - float64(cellY)))

			txc := int(float64(twidthc) * (floorX - float64(cellX)))
			tyc := int(float64(theightc) * (floorY - float64(cellY)))

			floorX += floorStepX
			floorY += floorStepY

			color := textures[6].At(txc, tyc)
			if isFloor {
				color = textures[3].At(txf, tyf)
			}
			screen.Set(x, y, color)
		}
	}
}

func drawWalls(screen *ebiten.Image, player game.Player, dists, indicies []float64, texs []uint8, height int) {
	pitch := player.Pitch * float64(height)
	for x := 0; x < len(dists); x++ {
		if texs[x] != 0 {
			wheight := 0
			if dists[x] != 0 {
				wheight = int(float64(height) / dists[x])
			}
			twidth := textures[texs[x]-1].Bounds().Max.X
			theight := textures[texs[x]-1].Bounds().Max.Y
			if wheight != 0 {
				imgx := int(indicies[x] * float64(twidth))
				start := int(math.Max(float64(height-wheight)/2+pitch+player.Z*float64(height)/dists[x], 0))
				end := int(math.Min(float64((height-wheight)/2+wheight)+pitch+player.Z*float64(height)/dists[x], float64(height)))
				for y := start; y < end; y++ {
					imgy := int(((float64(y-(height-wheight)/2) - pitch - player.Z*float64(height)/dists[x]) / float64(wheight)) * float64(theight))
					screen.Set(x, y, textures[texs[x]-1].At(imgx, imgy))
				}
			}
		}
	}
}

func drawSprites(screen *ebiten.Image, player game.Player, sprites []game.Sprite, dists []float64, dirX, dirY, planeX, planeY float64, width, height int) {
	distances := [][2]float64{}
	for i, sprite := range sprites {
		distances = append(distances, [2]float64{math.Pow(player.X-sprite.X, 2) + math.Pow(player.Y-sprite.Y, 2), float64(i)})
	}
	sort.SliceStable(distances, func(i, j int) bool {
		return distances[i][0] > distances[j][0]
	})
	pitch := player.Pitch * float64(height)

	for _, distance := range distances {
		sprite := sprites[int(distance[1])]

		relX := sprite.X - player.X
		relY := sprite.Y - player.Y

		invDet := 1 / (planeX*dirY - dirX*planeY)

		transformX := invDet * (dirY*relX - dirX*relY)
		transformY := invDet * (-planeY*relX + planeX*relY)

		screenX := int((float64(width) / 2) * (1 + transformX/transformY))

		textureWidth := float64(textures[sprite.Texture].Bounds().Max.X)
		textureHeight := float64(textures[sprite.Texture].Bounds().Max.Y)
		spriteWidth := int(math.Abs(float64(height)/transformY) * sprite.W)
		spriteHeight := int(math.Abs(float64(height)/transformY) * sprite.H)

		yMoveScreen := sprite.Z / transformY * float64(height)
		change := pitch - yMoveScreen

		start := int(math.Max(float64(screenX-spriteWidth/2), 0))
		end := int(math.Min(float64(screenX+spriteWidth/2), float64(width-1)))
		for x := start; x < end; x++ {

			imgx := int((float64(x-(screenX-spriteWidth/2)) / float64(spriteWidth)) * textureWidth)
			if dists[x] > transformY && transformY > 0 {

				start := int(math.Max(float64((height-spriteHeight)/2)+change+player.Z*float64(height)/transformY, 0))
				end := int(math.Min(float64((height-spriteHeight)/2+spriteHeight)+change+player.Z*float64(height)/transformY, float64(height)))
				for y := start; y < end; y++ {

					imgy := int(((float64(y-(height-spriteHeight)/2) - change - player.Z*float64(height)/transformY) / float64(spriteHeight)) * textureHeight)
					r, g, b, a := textures[sprite.Texture].At(imgx, imgy).RGBA()
					if a != 0 {
						screen.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
					}
				}
			}

		}
	}
}

func drawUI(screen *ebiten.Image, width, height int) {
	//Draw the crosshair
	geoM := ebiten.GeoM{}
	twidth, theight := images[0].Size()
	geoM.Translate(float64((width-twidth)/2), float64((height-theight)/2))
	screen.DrawImage(&images[0], &ebiten.DrawImageOptions{GeoM: geoM})

	//Draw the pistol
	geoM = ebiten.GeoM{}
	twidth, theight = images[1].Size()
	geoM.Translate(float64((width-twidth)/2), float64(height-theight))
	screen.DrawImage(&images[1], &ebiten.DrawImageOptions{GeoM: geoM})
}

//Draw2D draws a top-down view
func Draw2D(screen *ebiten.Image, cells [][]uint8, players []game.Player, playerSize float64, width, height int) {
	cellsizex := float64(width) / float64(len(cells[0]))
	cellsizey := float64(height) / float64(len(cells))
	for y := 0; y < len(cells); y++ {
		for x := 0; x < len(cells[0]); x++ {
			if cells[y][x] != 0 {
				ebitenutil.DrawRect(screen, float64(x)*cellsizex, float64(y)*cellsizey,
					cellsizex, cellsizey, color.RGBA{255, 0, 0, 255})
			}
		}
	}

	playerSizeX := cellsizex * playerSize
	playerSizeY := cellsizey * playerSize
	for _, player := range players {
		ebitenutil.DrawRect(screen, player.X*cellsizex-playerSizeX/2, player.Y*cellsizey-playerSizeY/2, playerSizeX, playerSizeY, color.RGBA{0, 255, 0, 255})
		ebitenutil.DrawLine(screen, player.X*cellsizex, player.Y*cellsizey, (player.X+math.Cos(player.Angle))*cellsizex, (player.Y+math.Sin(player.Angle))*cellsizey, color.RGBA{255, 255, 255, 255})
	}

	/*if len(players) == 2 {
		relX := players[1].X - players[0].X
		relY := players[1].Y - players[0].Y

		playerAngle := math.Atan2(relY, relX)
		dist := math.Sqrt(math.Pow(relX, 2) + math.Pow(relY, 2))
		angleWidth := math.Atan((playerSize / 2) / dist)
		//ebitenutil.DrawLine(screen, players[0].X*cellsizex, players[0].Y*cellsizey, (players[0].X+math.Cos(playerAngle)*dist)*cellsizex, (players[0].Y+math.Sin(playerAngle)*dist)*cellsizey, color.RGBA{0, 0, 255, 255})
		ebitenutil.DrawLine(screen, players[0].X*cellsizex, players[0].Y*cellsizey, (players[0].X+math.Cos(playerAngle+angleWidth)*dist)*cellsizex, (players[0].Y+math.Sin(playerAngle+angleWidth)*dist)*cellsizey, color.RGBA{0, 0, 255, 255})
		ebitenutil.DrawLine(screen, players[0].X*cellsizex, players[0].Y*cellsizey, (players[0].X+math.Cos(playerAngle-angleWidth)*dist)*cellsizex, (players[0].Y+math.Sin(playerAngle-angleWidth)*dist)*cellsizey, color.RGBA{0, 0, 255, 255})
	}*/
}

func load(path string) image.Image {
	file, err := os.Open(path)
	game.HandleError(err)
	img, _, err := image.Decode(file)
	game.HandleError(err)
	return img
}

func resizeImage(image image.Image, width, height int) *ebiten.Image {
	resized := resize.Resize(uint(width), uint(height), image, resize.NearestNeighbor)

	if cimg, ok := resized.(changeableImg); ok {
		w, h := cimg.Bounds().Max.X, cimg.Bounds().Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r, g, b, a := cimg.At(x, y).RGBA()
				if a > 0 {
					cimg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), 65535})
				}
			}
		}
	}
	img, _ := ebiten.NewImageFromImage(resized, ebiten.FilterDefault)
	return img
}

type changeableImg interface {
	Set(x, y int, c color.Color)
	Bounds() image.Rectangle
	At(x, y int) color.Color
}
