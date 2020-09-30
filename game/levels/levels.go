package levels

import (
	"github.com/oyberntzen/Raycasting-in-Golang/networking"
)

//SpriteInfo is info needed to create a new sprite
type SpriteInfo struct {
	Texture       uint8
	Width, Height uint
}

var (
	//LampInfo is used to create a lamp sprite in CreateSprite
	LampInfo SpriteInfo = SpriteInfo{10, 19, 10}
	//PillarInfo is used to create a pillar sprite in CreateSprite
	PillarInfo SpriteInfo = SpriteInfo{9, 16, 64}
	//BarrelInfo is used to create a barrel sprite in CreateSprite
	BarrelInfo SpriteInfo = SpriteInfo{8, 58, 64}
	//PlayerInfo is used to create a player sprite in CreateSprite
	PlayerInfo SpriteInfo = SpriteInfo{11, 32, 56}
)

//SpriteZOption is used in CreateSprite. You can choose between the sprite hanging in the ceiling, sitting on the floor or a specified Z value.
type SpriteZOption uint8

const (
	//SpriteZFloor is SpriteZOption for setting the sprite on the floor
	SpriteZFloor SpriteZOption = 0
	//SpriteZCeiling is SpriteZOption for hanging the sprite in the ceiling
	SpriteZCeiling SpriteZOption = 1
	//SpriteZFree is SpriteZOption for specifing your own z
	SpriteZFree SpriteZOption = 3
)

//Level is a struct for levels
type Level struct {
	Cells   [][]uint8
	Sprites []networking.Sprite
}

//Level01 is the first level
var Level01 Level = Level{[][]uint8{
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
	/*[]game.Sprite{
		CreateLamp(20.5, 11.5, 0.3),
		CreateLamp(18.5, 4.5, 0.3),
		CreateLamp(10.0, 4.5, 0.3),
		CreateLamp(10.0, 12.5, 0.3),
		CreateLamp(3.5, 6.5, 0.3),
		CreateLamp(3.5, 20.5, 0.3),
		CreateLamp(3.5, 14.5, 0.3),
		CreateLamp(14.5, 20.5, 0.3),
		CreatePillar(18.5, 10.5),
		CreatePillar(18.5, 11.5),
		CreatePillar(18.5, 12.5),
		CreateBarrel(21.5, 1.5, 0.4),
		CreateBarrel(15.5, 1.5, 0.4),
		CreateBarrel(16.0, 1.8, 0.4),
		CreateBarrel(16.2, 1.2, 0.4),
		CreateBarrel(3.5, 2.5, 0.4),
		CreateBarrel(9.5, 15.5, 0.4),
		CreateBarrel(10.0, 15.1, 0.4),
		CreateBarrel(10.5, 15.8, 0.4),
	},*/
	[]networking.Sprite{
		CreateSprite(LampInfo, 20.5, 11.5, 0, 0.3, 0, SpriteZCeiling),
		CreateSprite(LampInfo, 18.5, 4.5, 0, 0.3, 0, SpriteZCeiling),
		CreateSprite(LampInfo, 10.0, 4.5, 0, 0.3, 0, SpriteZCeiling),
		CreateSprite(LampInfo, 10.0, 12.5, 0, 0.3, 0, SpriteZCeiling),
		CreateSprite(LampInfo, 3.5, 6.5, 0, 0.3, 0, SpriteZCeiling),
		CreateSprite(LampInfo, 3.5, 20.5, 0, 0.3, 0, SpriteZCeiling),
		CreateSprite(LampInfo, 3.5, 14.5, 0, 0.3, 0, SpriteZCeiling),
		CreateSprite(LampInfo, 14.5, 20.5, 0, 0.3, 0, SpriteZCeiling),

		CreateSprite(PillarInfo, 18.5, 10.5, 0, 0, 1, SpriteZFloor),
		CreateSprite(PillarInfo, 18.5, 11.5, 0, 0, 1, SpriteZFloor),
		CreateSprite(PillarInfo, 18.5, 12.5, 0, 0, 1, SpriteZFloor),

		CreateSprite(BarrelInfo, 21.5, 1.5, 0, 0, 0.4, SpriteZFloor),
		CreateSprite(BarrelInfo, 15.5, 1.5, 0, 0, 0.4, SpriteZFloor),
		CreateSprite(BarrelInfo, 16.0, 1.8, 0, 0, 0.4, SpriteZFloor),
		CreateSprite(BarrelInfo, 16.2, 1.2, 0, 0, 0.4, SpriteZFloor),
		CreateSprite(BarrelInfo, 3.5, 2.5, 0, 0, 0.4, SpriteZFloor),
		CreateSprite(BarrelInfo, 9.5, 15.5, 0, 0, 0.4, SpriteZFloor),
		CreateSprite(BarrelInfo, 10.0, 15.1, 0, 0, 0.4, SpriteZFloor),
		CreateSprite(BarrelInfo, 10.5, 15.8, 0, 0, 0.4, SpriteZFloor),
	},
}

//CreateSprite creates a new sprite
func CreateSprite(spriteInfo SpriteInfo, x, y, z, width, height float64, zOption SpriteZOption) networking.Sprite {
	if width == 0 && height == 0 {
		if spriteInfo.Width > spriteInfo.Height {
			width = 1
		} else {
			height = 1
		}
	}
	if width == 0 {
		width = (float64(spriteInfo.Width) / float64(spriteInfo.Height)) * height
	} else if height == 0 {
		height = (float64(spriteInfo.Height) / float64(spriteInfo.Width)) * width
	}

	if zOption == SpriteZFloor {
		z = -(1 - float64(height)) / 2
	} else if zOption == SpriteZCeiling {
		z = (1 - float64(height)) / 2
	}

	return networking.Sprite{X: x, Y: y, Z: z, W: width, H: height, Texture: spriteInfo.Texture}
}

/*//CalculateZ calculates Z value of sprite when it is on the ground from height. Inverse value gives Z value of sprite when it hangs in the ceiling.
func CalculateZ(H float64) float64 {
	return -(1 - H) / 2
}

//CreateBarrel creates a new barrel sprite
func CreateBarrel(X, Y, Size float64) game.Sprite {
	textureWidth, textureHeight := 58.0, 64.0
	spriteWidth, spriteHeight := Size, Size
	if textureWidth > textureHeight {
		spriteHeight *= textureHeight / textureWidth
	} else {
		spriteWidth *= textureWidth / textureHeight
	}
	Z := CalculateZ(spriteHeight)
	return game.Sprite{X, Y, Z, spriteWidth, spriteHeight, 8}
}

//CreateLamp creates a new lamp sprite
func CreateLamp(X, Y, Size float64) game.Sprite {
	textureWidth, textureHeight := 19.0, 10.0
	spriteWidth, spriteHeight := Size, Size
	if textureWidth > textureHeight {
		spriteHeight *= textureHeight / textureWidth
	} else {
		spriteWidth *= textureWidth / textureHeight
	}
	Z := -CalculateZ(spriteHeight)
	return game.Sprite{X, Y, Z, spriteWidth, spriteHeight, 10}
}

//CreatePillar creates a new pillar sprite
func CreatePillar(X, Y float64) game.Sprite {
	textureWidth, textureHeight := 16.0, 64.0
	spriteWidth, spriteHeight := 1.0, 1.0
	if textureWidth > textureHeight {
		spriteHeight *= textureHeight / textureWidth
	} else {
		spriteWidth *= textureWidth / textureHeight
	}
	Z := 0.0
	return game.Sprite{X, Y, Z, spriteWidth, spriteHeight, 9}
}

//CreatePlayer creates a new player sprite
func CreatePlayer(X, Y float64) game.Sprite {
	textureWidth, textureHeight := 32.0, 56.0
	spriteWidth, spriteHeight := 0.6, 0.6
	if textureWidth > textureHeight {
		spriteHeight *= textureHeight / textureWidth
	} else {
		spriteWidth *= textureWidth / textureHeight
	}
	Z := CalculateZ(spriteHeight)
	return game.Sprite{X, Y, Z, spriteWidth, spriteHeight, 11}
}
*/
