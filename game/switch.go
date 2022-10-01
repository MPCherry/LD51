package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spriteSwitch, _, _ = ebitenutil.NewImageFromFile("resources/small_green_square.png")

type Switch struct {
	x         float64
	y         float64
	activated bool
	key       *Key
}

func (s *Switch) X() float64 {
	return s.x
}

func (s *Switch) Y() float64 {
	return s.y
}

func (s *Switch) Image() *ebiten.Image {
	return spriteSwitch
}

func (s *Switch) Active() bool {
	return true
}

func (s *Switch) Update(world *World, inputs []ebiten.Key) {
	for _, player := range world.players {
		if math.Abs(s.y-player.newY) < 16 && math.Abs(s.x-player.newX) < 16 {
			if s.key == nil {
				s.activated = true
			} else {
				if player.carrying != nil && player.carrying == s.key {
					s.activated = true
				}
			}
		}
	}
}
