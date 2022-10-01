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
	Active() bool
}

type World struct {
	keys          []ebiten.Key
	playerObjects []Object
	wallObjects   []Object
	players       []*Player
}

func NewWorld() *World {
	world := &World{}
	player := &Player{x: 16, y: 640 - 32, newX: 16, newY: 640 - 32, first: true, active: true}
	world.players = append(world.players, player)
	world.playerObjects = append(world.playerObjects, player)
	for i := 0; i < 60; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 640 - 16})
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 0})
	}
	for i := 1; i < 39; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 0, y: 16 * float64(i)})
		world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16, y: 16 * float64(i)})
	}
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 3, y: 640 - 16*3})

	return world
}

var recordCounter = 0
var respawnCounter = 0
var respawnShadowCounter = 0
var respawning = false

var shadowCounter = 0

func UpdateWorld(world *World) {
	world.keys = inpututil.AppendPressedKeys(world.keys[:0])

	for _, object := range world.playerObjects {
		if object.Active() {
			object.Update(world, world.keys)
		}
	}
	for _, object := range world.wallObjects {
		object.Update(world, world.keys)
	}
	for _, player := range world.players {
		player.x = player.newX
		player.y = player.newY
	}

	recordCounter++
	if recordCounter%300 == 0 {
		keyRecordCopy := make([][]ebiten.Key, len(keyRecording))
		copy(keyRecordCopy, keyRecording)
		keyRecording = [][]ebiten.Key{}

		player := &Player{x: 16, y: 640 - 32, newX: 16, newY: 640 - 32, first: false, keyRecord: keyRecordCopy, active: true}
		world.players = append(world.players, player)
		world.playerObjects = append(world.playerObjects, player)
		shadowCounter++
		respawning = true
		respawnShadowCounter = 0
		respawnCounter = 0
		for _, player := range world.players {
			player.x = -32
			player.newX = -32
			player.y = -32
			player.newY = -32
			player.active = false
		}
	}

	if respawning {
		if respawning && respawnCounter%60 == 0 {
			if respawnShadowCounter == len(world.players)-1 {
				world.players[0].x = 16
				world.players[0].newX = 16
				world.players[0].y = 640 - 32
				world.players[0].newY = 640 - 32
				world.players[0].jumped = false
				world.players[0].active = true
				respawning = false
				recordCounter = 0
			} else {
				world.players[1+respawnShadowCounter].x = 16
				world.players[1+respawnShadowCounter].newX = 16
				world.players[1+respawnShadowCounter].y = 640 - 32
				world.players[1+respawnShadowCounter].newY = 640 - 32
				world.players[1+respawnShadowCounter].jumped = false
				world.players[1+respawnShadowCounter].keyIndex = 0
				world.players[1+respawnShadowCounter].active = true
				respawnShadowCounter++
			}
		}
		respawnCounter++
	}
}

func DrawObjects(screen *ebiten.Image, world *World) {
	screen.Fill(color.White)
	for _, object := range append(world.playerObjects, world.wallObjects...) {
		if object.Active() {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(object.X(), object.Y())
			screen.DrawImage(object.Image(), op)
		}
	}
}
