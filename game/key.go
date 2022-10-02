package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spriteKey, _, _ = ebitenutil.NewImageFromFile("resources/sprites/key.png")

type Key struct {
	x            float64
	originX      float64
	y            float64
	originY      float64
	active       bool
	spriteIndex  int
	frame        int
	frameCounter int
}

func (k *Key) X() float64 {
	return k.x
}

func (k *Key) Y() float64 {
	return k.y
}

func (k *Key) Image() *ebiten.Image {
	return spriteKey.SubImage(image.Rect(16*k.spriteIndex, 16*k.frame, 16+16*k.spriteIndex, 16+16*k.frame)).(*ebiten.Image)
}

func (k *Key) Active() bool {
	return k.active
}

func (k *Key) Update(world *World, inputs []ebiten.Key) {
	k.frameCounter++
	if k.frameCounter == 10 {
		k.frame++
		k.frameCounter = 0
	}
	if k.frame == 4 {
		k.frame = 0
	}
}

func (w *Key) UpdateAnyway() {

}
