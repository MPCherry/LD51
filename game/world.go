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
	player  *Player
}

func NewWorld() *World {
	world := &World{}
	player := &Player{x: 16, y: 640 - 32, newX: 16, newY: 640 - 32}
	world.player = player
	world.objects = append(world.objects, player)
	for i := 0; i < 60; i++ {
		world.objects = append(world.objects, &Wall{x: 16 * float64(i), y: 640 - 16})
		world.objects = append(world.objects, &Wall{x: 16 * float64(i), y: 0})
	}
	for i := 1; i < 39; i++ {
		world.objects = append(world.objects, &Wall{x: 0, y: 16 * float64(i)})
		world.objects = append(world.objects, &Wall{x: 960 - 16, y: 16 * float64(i)})
	}
	world.objects = append(world.objects, &Wall{x: 16 * 3, y: 640 - 16*3})

	return world
}

func UpdateWorld(world *World) {
	world.keys = inpututil.AppendPressedKeys(world.keys[:0])

	for _, object := range world.objects {
		object.Update(world, world.keys)
	}
	world.player.x = world.player.newX
	world.player.y = world.player.newY
}

func DrawObjects(screen *ebiten.Image, world *World) {
	screen.Fill(color.White)
	for _, object := range world.objects {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(object.X(), object.Y())
		screen.DrawImage(object.Image(), op)
	}
}
