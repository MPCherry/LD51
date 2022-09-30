package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spritePlayer, _, _ = ebitenutil.NewImageFromFile("resources/small_black_square.png")

const (
	runSpeed   = 2
	jumpHeight = 5
)

type Player struct {
	x             float64
	newX          float64
	y             float64
	newY          float64
	verticalSpeed float64
	jumped        bool
}

func (p *Player) X() float64 {
	return p.x
}

func (p *Player) Y() float64 {
	return p.y
}

func (p *Player) Image() *ebiten.Image {
	return spritePlayer
}

func (p *Player) Update(world *World, inputs []ebiten.Key) {
	for _, k := range inputs {
		switch k {
		case ebiten.KeyLeft:
			p.newX = p.x - runSpeed
		case ebiten.KeyRight:
			p.newX = p.x + runSpeed
		case ebiten.KeyUp:
			if !p.jumped {
				p.verticalSpeed = -jumpHeight
				p.jumped = true
			}
		}
	}

	p.newY = p.y + p.verticalSpeed
	if p.verticalSpeed < 3 {
		p.verticalSpeed += 0.2
	}
}
