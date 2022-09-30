package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var sprite, _, _ = ebitenutil.NewImageFromFile("resources/sprites/dog.png")

type Dog struct {
	x float64
	y float64
}

func (d *Dog) X() float64 {
	return d.x
}

func (d *Dog) Y() float64 {
	return d.y
}

func (d *Dog) Image() *ebiten.Image {
	return sprite
}

func (d *Dog) Update(world *World, inputs []ebiten.Key) {
	for _, k := range inputs {
		switch k {
		case ebiten.KeyLeft:
			d.x--
		case ebiten.KeyRight:
			d.x++
		}
	}
}
