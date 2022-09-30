package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spriteWall, _, _ = ebitenutil.NewImageFromFile("resources/small_black_square.png")

type Wall struct {
	x float64
	y float64
}

func (d *Wall) X() float64 {
	return d.x
}

func (d *Wall) Y() float64 {
	return d.y
}

func (d *Wall) Image() *ebiten.Image {
	return spriteWall
}

func (d *Wall) Update(world *World, inputs []ebiten.Key) {

}
