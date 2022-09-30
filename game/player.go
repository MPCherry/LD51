package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var playerSprite, _, _ = ebitenutil.NewImageFromFile("resources/sprites/dog.png")

type Player struct {
	x float64
	y float64
}

func (d *Player) X() float64 {
	return d.x
}

func (d *Player) Y() float64 {
	return d.y
}

func (d *Player) Image() *ebiten.Image {
	return playerSprite
}

func (d *Player) Update(world *World) {}

func (d *Player) UpdatePlayer(world *World, inputs []ebiten.Key) {
	for _, k := range inputs {
		switch k {
		case ebiten.KeyLeft:
			d.x--
		case ebiten.KeyRight:
			d.x++
		}
	}
}
