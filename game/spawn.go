package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spriteSpawn, _, _ = ebitenutil.NewImageFromFile("resources/sprites/spawn.png")

type Spawn struct {
	x            float64
	y            float64
	active       bool
	spriteIndex  int
	frame        int
	frameCounter int
	life         int
}

func (s *Spawn) X() float64 {
	return s.x
}

func (s *Spawn) Y() float64 {
	return s.y
}

func (s *Spawn) Image() *ebiten.Image {
	return spriteSpawn.SubImage(image.Rect(16*s.spriteIndex, 16*s.frame, 16+16*s.spriteIndex, 16+16*s.frame)).(*ebiten.Image)
}

func (s *Spawn) Active() bool {
	return s.active
}

func (s *Spawn) Update(world *World, inputs []ebiten.Key) {

}

func (s *Spawn) UpdateAnyway() {
	s.frameCounter++
	if s.frameCounter == 10 {
		s.frame++
		s.frameCounter = 0
	}
	if s.frame == 4 {
		s.frame = 0
	}
	s.life--
	if s.life == 0 {
		s.active = false
	}
}
