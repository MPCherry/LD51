package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Object interface {
	X() float64
	Y() float64
	Image() *ebiten.Image
	Update(*World, []ebiten.Key)
}

type World struct {
	keys    []ebiten.Key
	objects []Object
}

func NewWorld() *World {
	world := &World{}
	world.objects = append(world.objects, &Player{x: 10, y: 200})

	return world
}

func UpdateWorld(world *World) {
	world.keys = inpututil.AppendPressedKeys(world.keys[:0])

	for _, object := range world.objects {
		object.Update(world, world.keys)
	}
}

func DrawObjects(screen *ebiten.Image, world *World) {
	screen.Fill(color.White)
	for _, object := range world.objects {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(object.X(), object.Y())
		screen.DrawImage(object.Image(), op)
	}
}
