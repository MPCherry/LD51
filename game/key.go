package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spriteKey, _, _ = ebitenutil.NewImageFromFile("resources/small_pink_square.png")

type Key struct {
	x       float64
	originX float64
	y       float64
	originY float64
	active  bool
}

func (k *Key) X() float64 {
	return k.x
}

func (k *Key) Y() float64 {
	return k.y
}

func (k *Key) Image() *ebiten.Image {
	return spriteKey
}

func (k *Key) Active() bool {
	return k.active
}

func (k *Key) Update(world *World, inputs []ebiten.Key) {
}