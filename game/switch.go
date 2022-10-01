package game

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spriteSwitch, _, _ = ebitenutil.NewImageFromFile("resources/green.png")

type Switch struct {
	x           float64
	y           float64
	activated   bool
	key         *Key
	resets      bool
	spriteIndex int
	final       bool
	cover       int
}

func (s *Switch) X() float64 {
	return s.x
}

func (s *Switch) Y() float64 {
	return s.y
}

func (s *Switch) Image() *ebiten.Image {
	return spriteSwitch.SubImage(image.Rect(16*s.spriteIndex, 0, 16*s.spriteIndex+16, 16)).(*ebiten.Image)
}

func (s *Switch) Active() bool {
	return !s.activated
}

func (s *Switch) Update(world *World, inputs []ebiten.Key) {
	for _, player := range world.players {
		if math.Abs(s.y-player.newY) < 16 && math.Abs(s.x-player.newX) < 16 {
			if s.key == nil {
				s.activated = true
				if s.cover != 0 {
					switch s.cover {
					case 2:
						bottomCoverDraw = false
					case 3:
						middleCoverDraw = false
					case 4:
						rightCoverDraw = false
					}
				}
			} else {
				if player.carrying != nil && player.carrying == s.key {
					s.activated = true
					if s.cover != 0 {
						switch s.cover {
						case 2:
							bottomCoverDraw = false
						case 3:
							middleCoverDraw = false
						case 4:
							rightCoverDraw = false
						}
					}
					player.carrying = nil
					s.key.active = false
				}
			}

			if s.final {
				gameover = true
				goCause = "win"
				fmt.Println("Won")
			}
		}
	}
}
