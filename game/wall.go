package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spriteWall, _, _ = ebitenutil.NewImageFromFile("resources/small_red_square.png")

type Wall struct {
	x float64
	y float64
}

func (w *Wall) X() float64 {
	return w.x
}

func (w *Wall) Y() float64 {
	return w.y
}

func (w *Wall) Image() *ebiten.Image {
	return spriteWall
}

func (w *Wall) Active() bool {
	return true
}

func (w *Wall) Update(world *World, inputs []ebiten.Key) {
	for _, player := range world.players {
		if math.Abs(w.y-player.newY) < 16 && player.newY < w.y && math.Abs(w.x-player.x) < 16 {
			player.newY = w.y - 16
			player.verticalSpeed = 0
			player.jumped = false
		}

		if math.Abs(w.y-player.newY) < 16 && player.newY > w.y && math.Abs(w.x-player.x) < 16 {
			player.newY = w.y + 16
			player.verticalSpeed = 0
		}
		if math.Abs(w.x-player.newX) < 16 && player.newX > w.x && math.Abs(w.y-player.y) < 16 {
			player.newX = w.x + 16
		}

		if math.Abs(w.x-player.newX) < 16 && player.newX < w.x && math.Abs(w.y-player.y) < 16 {
			player.newX = w.x - 16
		}
	}
}
