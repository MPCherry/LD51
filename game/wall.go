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

func (w *Wall) Update(world *World, inputs []ebiten.Key) {
	if math.Abs(w.y-world.player.newY) < 16 && world.player.newY < w.y && math.Abs(w.x-world.player.x) < 16 {
		world.player.newY = w.y - 16
		world.player.verticalSpeed = 0
		world.player.jumped = false
	}

	if math.Abs(w.y-world.player.newY) < 16 && world.player.newY > w.y && math.Abs(w.x-world.player.x) < 16 {
		world.player.newY = w.y + 16
		world.player.verticalSpeed = 0
	}
	if math.Abs(w.x-world.player.newX) < 16 && world.player.newX > w.x && math.Abs(w.y-world.player.y) < 16 {
		world.player.newX = w.x + 16
	}

	if math.Abs(w.x-world.player.newX) < 16 && world.player.newX < w.x && math.Abs(w.y-world.player.y) < 16 {
		world.player.newX = w.x - 16
	}
}
