package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spriteTeleporter, _, _ = ebitenutil.NewImageFromFile("resources/blue.png")

type Teleporter struct {
	x           float64
	y           float64
	destination *Teleporter
	spriteIndex int
}

func (t *Teleporter) X() float64 {
	return t.x
}

func (t *Teleporter) Y() float64 {
	return t.y
}

func (t *Teleporter) Image() *ebiten.Image {
	return spriteTeleporter.SubImage(image.Rect(16*t.spriteIndex, 0, 16*t.spriteIndex+16, 16)).(*ebiten.Image)
}

func (t *Teleporter) Active() bool {
	return true
}

func (t *Teleporter) Update(world *World, inputs []ebiten.Key) {
	if t.destination != nil {
		for _, player := range world.players {
			if math.Abs(t.y-player.newY) < 16 && math.Abs(t.x-player.newX) < 16 {
				player.newX = t.destination.x
				player.newY = t.destination.y
			}
		}
	}
}
